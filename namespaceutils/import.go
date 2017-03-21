package namespaceutils

import (
	"encoding/json"
	"strings"

	squallmodels "github.com/aporeto-inc/gaia/squallmodels/current/golang"

	"github.com/aporeto-inc/elemental"
	"github.com/aporeto-inc/manipulate"
)

func Import(manipulator manipulate.Manipulator, namespace string, content map[string]interface{}, shouldClean bool) error {

	if err := importNamespaceContent(manipulator, namespace, namespace, content, shouldClean); err != nil {
		return err
	}

	return nil
}

func importNamespaceContent(manipulator manipulate.Manipulator, topNamespace string, currentNamespace string, content map[string]interface{}, shouldClean bool) error {

	previousContent := elemental.IdentifiablesList{}

	for key, value := range content {
		if key != squallmodels.NamespaceIdentity.Category {
			continue
		}

		for _, n := range value.([]interface{}) {

			namespace := &squallmodels.Namespace{}
			namespaceContent := n.(map[string]interface{})[namespaceContentKey].(map[string]interface{})
			jsonRaw, err := json.Marshal(n)

			if err != nil {
				return err
			}

			if err := json.Unmarshal(jsonRaw, &namespace); err != nil {
				return err
			}

			isNsExist, err := isNamespaceExist(manipulator, currentNamespace, currentNamespace+"/"+namespace.Name)

			if err != nil {
				return err
			}

			if shouldClean && isNsExist {
				if err := deleteNamespace(manipulator, currentNamespace, currentNamespace+"/"+namespace.Name); err != nil {
					return err
				}
			}

			if isNsExist && !shouldClean {
				previousContent, err = ContentOfNamespace(manipulator, currentNamespace, false)

				if err != nil {
					return err
				}
			}

			if shouldClean || !isNsExist {
				if err := createNamespace(manipulator, currentNamespace, namespace); err != nil {
					return err
				}
			} else if isNsExist {
				namespace.Name = currentNamespace + "/" + namespace.Name
			}

			if err := importNamespaceContent(manipulator, topNamespace, namespace.Name, namespaceContent, shouldClean); err != nil {
				return err
			}
		}
	}

	if err := createContent(manipulator, topNamespace, currentNamespace, content); err != nil {
		return err
	}

	if !shouldClean {
		if err := deleteContent(manipulator, currentNamespace, previousContent); err != nil {
			return err
		}
	}

	return nil
}

func importComputeNamespace(namespace string, identityName string, object map[string]interface{}) {

	if namespace == "/" {
		namespace = ""
	}

	if identityName == squallmodels.APIAuthorizationPolicyIdentity.Category {
		object["authorizedNamespace"] = namespace + object["authorizedNamespace"].(string)
	}

	if identityName == squallmodels.NamespaceMappingPolicyIdentity.Category {
		object["mappedNamespace"] = namespace + object["mappedNamespace"].(string)
	}

	keys := []string{"subject", "object"}

	for _, key := range keys {
		if values, ok := object[key]; ok {
			for _, vs := range values.([]interface{}) {
				for i, v := range vs.([]interface{}) {
					s := strings.SplitN(v.(string), "=", 2)

					if s[0] == "$namespace" {
						newNamespace := namespace + s[1]
						vs.([]interface{})[i] = s[0] + "=" + newNamespace
					}
				}
			}
		}
	}
}

func createNamespace(manipulator manipulate.Manipulator, namespaceSession string, namespace *squallmodels.Namespace) error {
	mctx := manipulate.NewContext()
	mctx.Namespace = namespaceSession

	return manipulator.Create(mctx, namespace)
}

func deleteNamespace(manipulator manipulate.Manipulator, namespaceSession string, namespaceName string) error {
	mctx := manipulate.NewContext()
	mctx.Namespace = namespaceSession
	mctx.Filter = manipulate.NewFilterComposer().WithKey("namespace").Equals(namespaceName).Done()
	mctx.OverrideProtection = true

	return manipulator.DeleteMany(mctx, squallmodels.NamespaceIdentity)
}

func createContent(manipulator manipulate.Manipulator, topNamespace string, namespace string, content map[string]interface{}) error {
	for key, value := range content {

		if key == squallmodels.NamespaceIdentity.Category {
			continue
		}

		for _, object := range value.([]interface{}) {
			dest := squallmodels.IdentifiableForCategory(key).(elemental.Identifiable)
			importComputeNamespace(topNamespace, key, object.(map[string]interface{}))
			jsonRaw, err := json.Marshal(object)

			if err != nil {
				return err
			}

			if err := json.Unmarshal(jsonRaw, &dest); err != nil {
				return err
			}

			mctx := manipulate.NewContext()
			mctx.Namespace = namespace

			if err := manipulator.Create(mctx, dest); err != nil {
				return err
			}
		}
	}
	return nil
}

func deleteContent(manipulator manipulate.Manipulator, namespace string, content elemental.IdentifiablesList) error {
	for _, value := range content {
		mctx := manipulate.NewContext()
		mctx.Namespace = namespace
		mctx.OverrideProtection = true

		if value.Identity().Name == squallmodels.NamespaceIdentity.Category {
			continue
		}

		if err := manipulator.Delete(mctx, value); err != nil {
			return err
		}
	}

	return nil
}

func isNamespaceExist(manipulator manipulate.Manipulator, namespaceSession string, namespaceName string) (bool, error) {
	mctx := manipulate.NewContext()
	mctx.Namespace = namespaceSession
	mctx.Filter = manipulate.NewFilterComposer().WithKey("namespace").Equals(namespaceName).Done()

	dest := squallmodels.NamespacesList{}

	if err := manipulator.RetrieveMany(mctx, dest); err != nil {
		return false, err
	}

	if len(dest) > 0 {
		return true, nil
	}

	return false, nil
}
