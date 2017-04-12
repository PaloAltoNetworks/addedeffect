package namespaceutils

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/aporeto-inc/elemental"
	"github.com/aporeto-inc/manipulate"

	squallmodels "github.com/aporeto-inc/gaia/squallmodels/current/golang"
)

const namespaceContentKey = "content"

// ContentOfNamespace returns the content of the given namespace, if recursive is set it will return the content of its child namespace
func ContentOfNamespace(manipulator manipulate.Manipulator, namespace string, recursive bool, trackingID string) (elemental.IdentifiablesList, error) {

	identifiablesChannel := make(chan elemental.IdentifiablesList)
	errorsChannel := make(chan error)
	identifiables := elemental.IdentifiablesList{}

	mctx := manipulate.NewContext()
	mctx.Recursive = recursive
	mctx.Namespace = namespace
	mctx.ExternalTrackingID = trackingID
	mctx.ExternalTrackingType = "addedeffect.namespaceutils.export.contentofnamespace"

	for _, identity := range exportNamespacesObjects {
		go func() {
			dest := squallmodels.ContentIdentifiableForCategory(identity.Category)

			if err := manipulator.RetrieveMany(mctx, dest); err != nil {
				errorsChannel <- err
			}

			identifiablesChannel <- dest.List()
		}()

		select {
		case err := <-errorsChannel:
			return nil, err
		case ids := <-identifiablesChannel:
			identifiables = append(identifiables, ids...)
		}
	}

	return identifiables, nil
}

// TreeContentOfNamespace returns a tree of the given identifiables
// The main object of the tree is the namespace, it will have the public keys of the namespace + the key content
// content will contain the resources of the namespace
func TreeContentOfNamespace(namespace string, identifiables elemental.IdentifiablesList, trackingID string) (map[string]interface{}, error) {

	ns := &squallmodels.Namespace{}
	ns.Name = namespace

	namespaceContentRegistry := map[string][]map[string]map[string]interface{}{}
	identifiables = append(identifiables, ns)
	root := map[string]interface{}{}

	for _, identifiable := range identifiables {
		b, err := json.Marshal(identifiable)

		if err != nil {
			return nil, err
		}

		object := make(map[string]interface{})

		if err := json.Unmarshal(b, &object); err != nil {
			return nil, err
		}

		objectNamespace := reflect.ValueOf(identifiable).Elem().FieldByName("Namespace").String()

		if objectNamespace == namespace {
			objectNamespace = namespace[strings.LastIndex(namespace, "/")+1:]
		} else {
			objectNamespace = namespace[strings.LastIndex(namespace, "/")+1:] + strings.Replace(objectNamespace, namespace, "", 1)
		}

		attributeSpecifications := identifiable.(elemental.AttributeSpecifiable).AttributeSpecifications()
		FilterResourceField(attributeSpecifications, object)
		exportComputeNamespaceAttributes(namespace, identifiable.Identity().Category, object)

		if ns == identifiable {
			root = object
		} else {
			namespaceContentRegistry[objectNamespace] = append(namespaceContentRegistry[objectNamespace], map[string]map[string]interface{}{identifiable.Identity().Category: object})
		}
	}

	fillTreeForNamespace("", root, namespaceContentRegistry)
	return map[string]interface{}{squallmodels.NamespaceIdentity.Category: []interface{}{root}}, nil
}

func fillTreeForNamespace(namespace string, currentNamespace map[string]interface{}, namespaceContentRegistry map[string][]map[string]map[string]interface{}) {

	currentNamespace[namespaceContentKey] = map[string][]map[string]interface{}{}
	fullNamespaceName := namespace + currentNamespace["name"].(string)

	for _, objects := range namespaceContentRegistry[fullNamespaceName] {

		for identity, object := range objects {

			if identity == squallmodels.NamespaceIdentity.Category {
				newNamespace := fullNamespaceName + "/"

				if fullNamespaceName == "" {
					newNamespace = ""
				}
				fillTreeForNamespace(newNamespace, object, namespaceContentRegistry)
			}

			currentNamespace[namespaceContentKey].(map[string][]map[string]interface{})[identity] = append(currentNamespace[namespaceContentKey].(map[string][]map[string]interface{})[identity], object)
		}
	}
}

func exportComputeNamespace(namespace string, objectNamespace string) string {

	if objectNamespace == namespace {
		return namespace[strings.LastIndex(namespace, "/"):]
	}

	return namespace[strings.LastIndex(namespace, "/"):] + strings.Replace(objectNamespace, namespace, "", 1)
}

func exportComputeNamespaceAttributes(namespace string, identityName string, object map[string]interface{}) {

	if identityName == squallmodels.NamespaceIdentity.Category {
		object["name"] = object["name"].(string)[strings.LastIndex(object["name"].(string), "/")+1:]
	}

	if identityName == squallmodels.APIAuthorizationPolicyIdentity.Category {
		object["authorizedNamespace"] = exportComputeNamespace(namespace, object["authorizedNamespace"].(string))
	}

	if identityName == squallmodels.NamespaceMappingPolicyIdentity.Category {
		object["mappedNamespace"] = exportComputeNamespace(namespace, object["mappedNamespace"].(string))
	}

	keys := []string{"subject", "object"}

	for _, key := range keys {
		if values, ok := object[key]; ok {
			for _, vs := range values.([]interface{}) {
				for i, v := range vs.([]interface{}) {
					s := strings.SplitN(v.(string), "=", 2)

					if s[0] == "$namespace" {
						newNamespace := exportComputeNamespace(namespace, s[1])
						vs.([]interface{})[i] = s[0] + "=" + newNamespace
					}
				}
			}
		}
	}

}
