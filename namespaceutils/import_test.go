package namespaceutils

import (
	"testing"

	squallmodels "github.com/aporeto-inc/gaia/squallmodels/current/golang"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_importComputeNamespace(t *testing.T) {
	Convey("Given test data is prepared", t, func() {
		apiAuthorizationPolicy1 := map[string]interface{}{"authorizedNamespace": "/3/4"}
		apiAuthorizationPolicy2 := map[string]interface{}{"authorizedNamespace": "/3"}

		namespaceMappingPolicy1 := map[string]interface{}{"mappedNamespace": "/3/4", "object": []interface{}{[]interface{}{"$namespace=/3"}, []interface{}{"$namespace=/3/4"}, []interface{}{"$namespace=/3/5"}}}
		namespaceMappingPolicy2 := map[string]interface{}{"mappedNamespace": "/3", "subject": []interface{}{[]interface{}{"$namespace=/3"}, []interface{}{"$namespace=/3/4"}, []interface{}{"$namespace=/3/5"}}}

		Convey("Then I try to compute the data with namespace /1/2", func() {
			namespace := "/1/2"
			importComputeNamespace(namespace, squallmodels.APIAuthorizationPolicyIdentity.Category, apiAuthorizationPolicy1)
			importComputeNamespace(namespace, squallmodels.APIAuthorizationPolicyIdentity.Category, apiAuthorizationPolicy2)
			importComputeNamespace(namespace, squallmodels.NamespaceMappingPolicyIdentity.Category, namespaceMappingPolicy1)
			importComputeNamespace(namespace, squallmodels.NamespaceMappingPolicyIdentity.Category, namespaceMappingPolicy2)

			So(apiAuthorizationPolicy1["authorizedNamespace"], ShouldEqual, "/1/2/3/4")
			So(apiAuthorizationPolicy2["authorizedNamespace"], ShouldEqual, "/1/2/3")
			So(namespaceMappingPolicy1["mappedNamespace"], ShouldEqual, "/1/2/3/4")
			So(namespaceMappingPolicy2["mappedNamespace"], ShouldEqual, "/1/2/3")
			So(namespaceMappingPolicy1["object"], ShouldResemble, []interface{}{[]interface{}{"$namespace=/1/2/3"}, []interface{}{"$namespace=/1/2/3/4"}, []interface{}{"$namespace=/1/2/3/5"}})
			So(namespaceMappingPolicy2["subject"], ShouldResemble, []interface{}{[]interface{}{"$namespace=/1/2/3"}, []interface{}{"$namespace=/1/2/3/4"}, []interface{}{"$namespace=/1/2/3/5"}})
		})
	})

	Convey("Given a second test data are prepared", t, func() {
		apiAuthorizationPolicy1 := map[string]interface{}{"authorizedNamespace": "/apomux/4"}
		apiAuthorizationPolicy2 := map[string]interface{}{"authorizedNamespace": "/apomux"}

		namespaceMappingPolicy1 := map[string]interface{}{"mappedNamespace": "/apomux/4", "object": []interface{}{[]interface{}{"$namespace=/apomux"}, []interface{}{"$namespace=/apomux/4"}, []interface{}{"$namespace=/apomux/5"}}}
		namespaceMappingPolicy2 := map[string]interface{}{"mappedNamespace": "/apomux", "subject": []interface{}{[]interface{}{"$namespace=/apomux"}, []interface{}{"$namespace=/apomux/4"}, []interface{}{"$namespace=/apomux/5"}}}

		Convey("Then I try to compute the data with namespace /1/2", func() {
			namespace := "/1/2"
			importComputeNamespace(namespace, squallmodels.APIAuthorizationPolicyIdentity.Category, apiAuthorizationPolicy1)
			importComputeNamespace(namespace, squallmodels.APIAuthorizationPolicyIdentity.Category, apiAuthorizationPolicy2)
			importComputeNamespace(namespace, squallmodels.NamespaceMappingPolicyIdentity.Category, namespaceMappingPolicy1)
			importComputeNamespace(namespace, squallmodels.NamespaceMappingPolicyIdentity.Category, namespaceMappingPolicy2)

			So(apiAuthorizationPolicy1["authorizedNamespace"], ShouldEqual, "/1/2/apomux/4")
			So(apiAuthorizationPolicy2["authorizedNamespace"], ShouldEqual, "/1/2/apomux")
			So(namespaceMappingPolicy1["mappedNamespace"], ShouldEqual, "/1/2/apomux/4")
			So(namespaceMappingPolicy2["mappedNamespace"], ShouldEqual, "/1/2/apomux")
			So(namespaceMappingPolicy1["object"], ShouldResemble, []interface{}{[]interface{}{"$namespace=/1/2/apomux"}, []interface{}{"$namespace=/1/2/apomux/4"}, []interface{}{"$namespace=/1/2/apomux/5"}})
			So(namespaceMappingPolicy2["subject"], ShouldResemble, []interface{}{[]interface{}{"$namespace=/1/2/apomux"}, []interface{}{"$namespace=/1/2/apomux/4"}, []interface{}{"$namespace=/1/2/apomux/5"}})
		})
	})

	Convey("Given a third test data are prepared", t, func() {
		apiAuthorizationPolicy1 := map[string]interface{}{"authorizedNamespace": "/apomux/4"}
		apiAuthorizationPolicy2 := map[string]interface{}{"authorizedNamespace": "/apomux"}

		namespaceMappingPolicy1 := map[string]interface{}{"mappedNamespace": "/apomux/4", "object": []interface{}{[]interface{}{"$namespace=/apomux"}, []interface{}{"$namespace=/apomux/4"}, []interface{}{"$namespace=/apomux/5"}}}
		namespaceMappingPolicy2 := map[string]interface{}{"mappedNamespace": "/apomux", "subject": []interface{}{[]interface{}{"$namespace=/apomux"}, []interface{}{"$namespace=/apomux/4"}, []interface{}{"$namespace=/apomux/5"}}}

		Convey("Then I try to compute the data with namespace /", func() {
			namespace := "/"
			importComputeNamespace(namespace, squallmodels.APIAuthorizationPolicyIdentity.Category, apiAuthorizationPolicy1)
			importComputeNamespace(namespace, squallmodels.APIAuthorizationPolicyIdentity.Category, apiAuthorizationPolicy2)
			importComputeNamespace(namespace, squallmodels.NamespaceMappingPolicyIdentity.Category, namespaceMappingPolicy1)
			importComputeNamespace(namespace, squallmodels.NamespaceMappingPolicyIdentity.Category, namespaceMappingPolicy2)

			So(apiAuthorizationPolicy1["authorizedNamespace"], ShouldEqual, "/apomux/4")
			So(apiAuthorizationPolicy2["authorizedNamespace"], ShouldEqual, "/apomux")
			So(namespaceMappingPolicy1["mappedNamespace"], ShouldEqual, "/apomux/4")
			So(namespaceMappingPolicy2["mappedNamespace"], ShouldEqual, "/apomux")
			So(namespaceMappingPolicy1["object"], ShouldResemble, []interface{}{[]interface{}{"$namespace=/apomux"}, []interface{}{"$namespace=/apomux/4"}, []interface{}{"$namespace=/apomux/5"}})
			So(namespaceMappingPolicy2["subject"], ShouldResemble, []interface{}{[]interface{}{"$namespace=/apomux"}, []interface{}{"$namespace=/apomux/4"}, []interface{}{"$namespace=/apomux/5"}})
		})
	})
}
