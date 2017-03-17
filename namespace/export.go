package namespace

import (
	"encoding/json"
	"reflect"
	"strings"

	squallmodels "github.com/aporeto-inc/gaia/squallmodels/current/golang"

	"github.com/aporeto-inc/elemental"
	"github.com/aporeto-inc/manipulate"
)

func ContentOfNamespace(manipulator manipulate.Manipulator, namespace string) (elemental.IdentifiablesList, error) {

	identifiablesChannel := make(chan elemental.IdentifiablesList)
	errorsChannel := make(chan error)
	identifiables := elemental.IdentifiablesList{}

	mctx := manipulate.NewContext()
	mctx.Recursive = true
	mctx.Namespace = namespace

	for _, identity := range exportNamespacesObjects {
		go func() {
			dest := squallmodels.ContentIdentifiableForIdentity(identity.Name)

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

func TreeContentOfNamespace(manipulator manipulate.Manipulator, namespace string) (map[string]interface{}, error) {

	ns := &squallmodels.Namespace{}
	identifiables, err := ContentOfNamespace(manipulator, namespace)

	if err != nil {
		return nil, err
	}

	namespaceContentRegistry := map[string][]map[string]map[string]interface{}{}
	identifiables = append(identifiables, ns)
	yamlNamespace := map[string]interface{}{}

	for _, identifiable := range identifiables {
		b, err := json.Marshal(identifiable)

		if err != nil {
			return nil, err
		}

		object := make(map[string]interface{})

		if err := json.Unmarshal(b, &object); err != nil {
			return nil, err
		}

		objectNamespace := reflect.ValueOf(identifiable).FieldByName("Namespace").String()

		if objectNamespace == namespace {
			objectNamespace = namespace[strings.LastIndex(namespace, "/"):]
		} else {
			objectNamespace = namespace[strings.LastIndex(namespace, "/"):] + strings.Replace(objectNamespace, namespace, "", 1)
		}

		attributeSpecifications := identifiable.(elemental.AttributeSpecifiable).AttributeSpecifications()
		FilterResourceField(attributeSpecifications, object)
		computeNamespaceAttribute(namespace, identifiable.Identity().Name, object)

		if ns == identifiable {
			yamlNamespace = object
		} else {
			namespaceContentRegistry[objectNamespace] = append(namespaceContentRegistry[strings.Replace(objectNamespace, namespace, "", 1)], map[string]map[string]interface{}{identifiable.Identity().Name: object})
		}
	}

	fillTreeNamespaceContent(yamlNamespace, namespaceContentRegistry)

	return map[string]interface{}{"namespace": yamlNamespace}, nil
}

func fillTreeNamespaceContent(currentNamespace map[string]interface{}, namespaceContentRegistry map[string][]map[string]map[string]interface{}) {

	currentNamespace["resources"] = map[string][]map[string]interface{}{}
	for _, objects := range namespaceContentRegistry[currentNamespace["name"].(string)] {

		for identity, object := range objects {

			if identity == squallmodels.NamespaceIdentity.Name {
				fillTreeNamespaceContent(object, namespaceContentRegistry)
			}

			currentNamespace["resources"].(map[string][]map[string]interface{})[identity] = append(currentNamespace["resources"].(map[string][]map[string]interface{})[identity], object)
		}
	}
}

func computeNamespaceAttribute(namespace string, identityName string, object map[string]interface{}) {

	if identityName == squallmodels.NamespaceIdentity.Name {
		object["name"] = namespace[strings.LastIndex(namespace, "/")+1:]
	}

	if identityName == squallmodels.APIAuthorizationPolicyIdentity.Name {
		if object["authorizedNamespace"] == namespace {
			object["authorizedNamespace"] = namespace[strings.LastIndex(namespace, "/"):]
		} else {
			object["authorizedNamespace"] = namespace[strings.LastIndex(namespace, "/"):] + strings.Replace(object["authorizatedNamespace"].(string), namespace, "", 1)
		}
	}

	if identityName == squallmodels.NamespaceMappingPolicyIdentity.Name {
		if object["mappedNamespace"] == namespace {
			object["mappedNamespace"] = namespace[strings.LastIndex(namespace, "/"):]
		} else {
			object["mappedNamespace"] = namespace[strings.LastIndex(namespace, "/"):] + strings.Replace(object["mappedNamespace"].(string), namespace, "", 1)
		}
	}
}

// filterResourceField receives a set of attribute specifications, and a set of
// key-value pairs belonging to that category, and removes a field if any of
// the following criteria applies:
//  1) is readOnly, or is not exposed;
//	2) is an empty value, i.e. an empty string, integer, or boolean.
func FilterResourceField(attribMap map[string]elemental.AttributeSpecification, object map[string]interface{}) {

	if attribMap == nil {
		return
	}

	for _, spec := range attribMap {
		key := spec.Name
		val, ok := object[key]
		if !ok {
			continue
		}
		// 1) If (readOnly || !exposed), remove.
		doRemove := spec.ReadOnly || !spec.Exposed
		// 2) If value is empty, remove.
		if !doRemove {
			switch spec.Type {
			case "string":
				doRemove = reflect.ValueOf(val).Len() == 0
			case "integer":
				doRemove = reflect.DeepEqual(val, 0)
			case "boolean":
				doRemove = reflect.DeepEqual(val, false)
			}
		}
		// Do remove
		if doRemove {
			delete(object, key)
		}
	}
}
