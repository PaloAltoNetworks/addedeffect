package namespaceutils

import (
	"encoding/json"
	"fmt"
	"strings"

	squallmodels "github.com/aporeto-inc/gaia/squallmodels/current/golang"

	"github.com/aporeto-inc/elemental"
	"github.com/aporeto-inc/manipulate"
)

// Import the given content to the given namespace
// If shouldClean is set, Import will clean previous namespaces if it overlapps with a new one
// If shouldClean is not set, Import will clean the content of the namespace if it overlapps with a new one, just the namespace will remain
func Import(manipulator manipulate.Manipulator, namespace string, content map[string]interface{}, shouldClean bool) error {
	return importNamespaceContent(manipulator, namespace, namespace, content, shouldClean)
}

// importNamespaceContent is a recursive function
// The function will create namespaces first and then content of them
// It will first check if the namespace exists, if yes it will delete it if shouldClean is set, otherwise it will retrieve the content of it
// Then we create the namesapce if needed, create the content and finally delete the previous content
func importNamespaceContent(manipulator manipulate.Manipulator, topNamespace string, currentNamespace string, content map[string]interface{}, shouldClean bool) error {

	previousContent := elemental.IdentifiablesList{}
	originalNamespaceName := ""

	for key, value := range content {

		key = strings.ToLower(key)

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

			if err = json.Unmarshal(jsonRaw, &namespace); err != nil {
				return err
			}

			isNamespaceExists := true

			if namespace.Name != "" {
				originalNamespaceName = namespace.Name
				if currentNamespace == "/" {
					namespace.Name = "/" + namespace.Name
				} else {
					namespace.Name = currentNamespace + "/" + namespace.Name
				}

				isNamespaceExists, err = isNamespaceExist(manipulator, currentNamespace, namespace)

				if err != nil {
					return err
				}

				if shouldClean && isNamespaceExists {
					if err = deleteNamespace(manipulator, currentNamespace, namespace); err != nil {
						return err
					}
				}
			}

			if (isNamespaceExists && !shouldClean) || (originalNamespaceName == "") {
				previousContent, err = ContentOfNamespace(manipulator, currentNamespace, false)

				if err != nil {
					return err
				}
			}

			if (shouldClean || !isNamespaceExists) && originalNamespaceName != "" {
				// When we create a namespace, we are not allowed to put some /
				newNamespace := &squallmodels.Namespace{}
				newNamespace.Name = originalNamespaceName
				if err := createNamespace(manipulator, currentNamespace, newNamespace); err != nil {
					return err
				}
			}

			if err := importNamespaceContent(manipulator, topNamespace, namespace.Name, namespaceContent, shouldClean); err != nil {
				return err
			}
		}
	}

	if err := createContent(manipulator, topNamespace, currentNamespace, content); err != nil {
		return err
	}

	if !shouldClean || (originalNamespaceName == "") {
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

func deleteNamespace(manipulator manipulate.Manipulator, namespaceSession string, namespace *squallmodels.Namespace) error {
	mctx := manipulate.NewContext()
	mctx.Namespace = namespaceSession
	mctx.OverrideProtection = true

	return manipulator.Delete(mctx, namespace)
}

func createContent(manipulator manipulate.Manipulator, topNamespace string, namespace string, content map[string]interface{}) error {
	for key, value := range content {

		key = strings.ToLower(key)

		if key == squallmodels.NamespaceIdentity.Category {
			continue
		}

		for _, object := range value.([]interface{}) {
			dest := squallmodels.IdentifiableForCategory(key).(elemental.Identifiable)

			if dest == nil {
				return fmt.Errorf("The given key %s is not valid", key)
			}

			// For instance, values as /apomux needs to be /level/apomux when adding the content in /level/apomux
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

		if value.Identity().Name == squallmodels.NamespaceIdentity.Name {
			continue
		}

		if err := manipulator.Delete(mctx, value); err != nil {
			return err
		}
	}

	return nil
}

func isNamespaceExist(manipulator manipulate.Manipulator, namespaceSession string, namespace *squallmodels.Namespace) (bool, error) {
	mctx := manipulate.NewContext()
	mctx.Namespace = namespaceSession
	mctx.Filter = manipulate.NewFilterComposer().WithKey("name").Equals(namespace.Name).Done()

	dest := squallmodels.NamespacesList{}

	if err := manipulator.RetrieveMany(mctx, &dest); err != nil {
		return false, err
	}

	if len(dest) > 0 {
		namespace.SetIdentifier(dest[0].Identifier())
		return true, nil
	}

	return false, nil
}
