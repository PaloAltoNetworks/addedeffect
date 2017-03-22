package namespaceutils

import (
	"reflect"

	"github.com/aporeto-inc/elemental"
)

// FilterResourceField receives a set of attribute specifications, and a set of
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
		val := object[key]

		// 1) If (readOnly || !exposed), remove.
		doRemove := spec.ReadOnly || !spec.Exposed || val == nil
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

		if !doRemove {
			if (reflect.ValueOf(val).Kind() == reflect.Map || reflect.ValueOf(val).Kind() == reflect.Slice || reflect.ValueOf(val).Kind() == reflect.Array) && reflect.ValueOf(val).Len() == 0 {
				doRemove = true
			}
		}

		// Do remove
		if doRemove {
			delete(object, key)
		}
	}
}
