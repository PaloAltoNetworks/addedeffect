package namespaceutils

import (
	"encoding/json"
	"fmt"
	"strings"

	squallmodels "github.com/aporeto-inc/gaia/squallmodels/current/golang"

	"github.com/aporeto-inc/elemental"
	"github.com/aporeto-inc/manipulate"
)

func Import(manipulator manipulate.Manipulator, namespace string, content map[string]interface{}, shouldClean bool) error {

	if _, ok := content["namespace"]; !ok {
		return fmt.Errorf("The given content should have a key namespace")
	}

	topNamespace := content["namespace"].(map[string]interface{})["name"].(string)

	mctx := manipulate.NewContext()
	mctx.Namespace = namespace
	mctx.OverrideProtection = true
	mctx.Filter = manipulate.NewFilterComposer().WithKey("namespace").Equals(topNamespace).Done()

	previousContent := elemental.IdentifiablesList{}
	namespaces := squallmodels.NamespacesList{}

	manipulator.RetrieveMany(mctx, namespaces)

	if len(namespaces) == 1 && shouldClean {
		mctx = manipulate.NewContext()
		mctx.Namespace = namespace
		mctx.OverrideProtection = true

		if err := manipulator.Delete(mctx, namespaces[0]); err != nil {
			return err
		}
	} else if len(namespaces) == 1 {
		for _, value := range exportNamespacesObjects {
			mctx = manipulate.NewContext()
			mctx.Namespace = topNamespace
			mctx.Recursive = true

			dest := squallmodels.ContentIdentifiableForCategory(value.Category)

			if err := manipulator.RetrieveMany(mctx, dest); err != nil {
				return err
			}

			previousContent = append(previousContent, dest.List()...)
		}
	}

	if err := importNamespaceContent(manipulator, namespace, content["namespace"].(map[string]interface{})); err != nil {
		return err
	}

	if !shouldClean && len(namespaces) == 1 {
		for _, value := range previousContent {
			mctx := manipulate.NewContext()
			mctx.Namespace = namespace
			mctx.OverrideProtection = true

			if err := manipulator.Delete(mctx, value); err != nil {
				return err
			}
		}
	}

	return nil
}

func importNamespaceContent(manipulator manipulate.Manipulator, currentNamespace string, content map[string]interface{}) error {

	namespaceContent := content[namespaceContentKey].(map[string]interface{})
	delete(content, namespaceContentKey)

	namespace := &squallmodels.Namespace{}
	jsonRaw, err := json.Marshal(content)

	if err != nil {
		return err
	}

	if err := json.Unmarshal(jsonRaw, &namespace); err != nil {
		return err
	}

	mctx := manipulate.NewContext()
	mctx.Namespace = currentNamespace

	if err := manipulator.Create(mctx, namespace); err != nil {
		return err
	}

	for key, value := range namespaceContent {
		if key == squallmodels.NamespaceIdentity.Category {
			for _, n := range value.([]interface{}) {
				if err := importNamespaceContent(manipulator, namespace.Name, n.(map[string]interface{})); err != nil {
					return err
				}
			}
		}
	}

	for key, value := range namespaceContent {

		if key == squallmodels.NamespaceIdentity.Category {
			continue
		}

		for _, object := range value.([]interface{}) {

			dest := squallmodels.ContentIdentifiableForCategory(key).(elemental.Identifiable)
			importComputeNamespace(namespace.Name, key, object.(map[string]interface{}))
			jsonRaw, err := json.Marshal(object)

			if err != nil {
				return err
			}

			if err := json.Unmarshal(jsonRaw, &dest); err != nil {
				return err
			}

			if err := manipulator.Create(mctx, dest); err != nil {
				return err
			}

		}
	}

	return nil
}

func importComputeNamespace(namespace string, identityName string, object map[string]interface{}) {

	if identityName == squallmodels.APIAuthorizationPolicyIdentity.Category {
		object["authorizedNamespace"] = namespace[:len(namespace)-2] + object["authorizedNamespace"].(string)
	}

	if identityName == squallmodels.NamespaceMappingPolicyIdentity.Category {
		object["mappedNamespace"] = namespace[:len(namespace)-2] + object["mappedNamespace"].(string)
	}

	keys := []string{"subject", "object"}

	for _, key := range keys {
		if values, ok := object[key]; ok {
			for _, vs := range values.([]interface{}) {
				for i, v := range vs.([]interface{}) {
					s := strings.SplitN(v.(string), "=", 2)

					if s[0] == "$namespace" {
						newNamespace := namespace[:len(namespace)-2] + s[1]
						vs.([]interface{})[i] = s[0] + "=" + newNamespace
					}
				}
			}
		}
	}
}
