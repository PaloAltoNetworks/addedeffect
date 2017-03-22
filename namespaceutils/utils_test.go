package namespaceutils

import (
	"testing"

	"github.com/aporeto-inc/elemental"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_FilterResourceField(t *testing.T) {
	Convey("Given test data is prepared", t, func() {
		attribMap := map[string]elemental.AttributeSpecification{
			"StringField": elemental.AttributeSpecification{
				Name:     "stringField",
				Type:     "string",
				ReadOnly: false,
				Exposed:  true,
			},
			"IntegerField": elemental.AttributeSpecification{
				Name:     "integerField",
				Type:     "integer",
				ReadOnly: false,
				Exposed:  true,
			},
			"BooleanField": elemental.AttributeSpecification{
				Name:     "booleanField",
				Type:     "boolean",
				ReadOnly: false,
				Exposed:  true,
			},
			"ReadOnlyField": elemental.AttributeSpecification{
				Name:     "readOnlyField",
				Type:     "string",
				ReadOnly: true,
				Exposed:  true,
			},
			"NonExposedField": elemental.AttributeSpecification{
				Name:     "nonExposedField",
				Type:     "string",
				ReadOnly: false,
				Exposed:  false,
			},
		}

		testObj := map[string]interface{}{
			"unrelatedField": "someValue",
		}
		benchmarkObj := map[string]interface{}{
			"unrelatedField": "someValue",
		}

		// [<key>, <removedValue>, <keptValue>]
		singleFieldTestCases := [][3]interface{}{
			{"stringField", "", "testStringValue"},
			{"integerField", 0, 11223344},
			{"booleanField", false, true},
			{"readOnlyField", "anyValue", nil},
			{"nonExposedField", "anyValue", nil},
		}

		for _, testCase := range singleFieldTestCases {
			key := testCase[0].(string)
			removedValue := testCase[1]
			keptValue := testCase[2]

			Convey("It should remove "+key, func() {
				testObj[key] = removedValue
				FilterResourceField(attribMap, testObj)
				So(testObj, ShouldResemble, benchmarkObj)
			})

			if keptValue != nil {
				Convey("It should not remove non-empty "+key, func() {
					testObj[key] = keptValue
					benchmarkObj[key] = keptValue
					FilterResourceField(attribMap, testObj)
					So(testObj, ShouldResemble, benchmarkObj)
				})
			}
		}
	})
}
