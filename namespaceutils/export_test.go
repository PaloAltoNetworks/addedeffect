package namespaceutils

import (
	"testing"

	"github.com/aporeto-inc/elemental"
	squallmodels "github.com/aporeto-inc/gaia/squallmodels/current/golang"
	. "github.com/smartystreets/goconvey/convey"
)

var coucou *testing.T

func Test_TreeContentOfNamespace(t *testing.T) {
	Convey("Given test data is prepared", t, func() {
		namespace := "/3"

		namespace1 := &squallmodels.Namespace{Name: "4", Namespace: "/3"}
		namespace2 := &squallmodels.Namespace{Name: "4.1", Namespace: "/3"}
		namespace3 := &squallmodels.Namespace{Name: "5", Namespace: "/3/4"}

		externalService1 := &squallmodels.ExternalService{Name: "externalService1", Namespace: "/3"}
		externalService2 := &squallmodels.ExternalService{Name: "externalService2", Namespace: "/3"}
		externalService3 := &squallmodels.ExternalService{Name: "externalService3", Namespace: "/3/4/5"}
		externalService4 := &squallmodels.ExternalService{Name: "externalService4", Namespace: "/3/4.1"}

		filepath1 := &squallmodels.FilePath{Name: "filePath1", Namespace: "/3"}
		filepath2 := &squallmodels.FilePath{Name: "filePath2", Namespace: "/3/4/5"}
		filepath3 := &squallmodels.FilePath{Name: "filePath3", Namespace: "/3/4/5"}
		filepath4 := &squallmodels.FilePath{Name: "filePath4", Namespace: "/3/4.1"}

		apiAuthorizationPolicy1 := &squallmodels.APIAuthorizationPolicy{Name: "api1", Namespace: "/3", AuthorizedNamespace: "/3/4/5", Subject: [][]string{[]string{"$namespace=/3/4/5"}}}
		apiAuthorizationPolicy2 := &squallmodels.APIAuthorizationPolicy{Name: "api2", Namespace: "/3", AuthorizedNamespace: "/3/4/5", Subject: [][]string{}}

		identifiables := elemental.IdentifiablesList{namespace1, namespace2, namespace3, externalService1, externalService2, externalService3, externalService4, filepath1, filepath2, filepath3, filepath4, apiAuthorizationPolicy1, apiAuthorizationPolicy2}

		Convey("Then I create my tree", func() {
			coucou = t
			tree, err := TreeContentOfNamespace(namespace, identifiables)
			namespaceMap1 := map[string]interface{}{"name": "4"}
			namespaceMap2 := map[string]interface{}{"name": "4.1"}
			namespaceMap3 := map[string]interface{}{"name": "5"}

			externalServiceMap1 := map[string]interface{}{"name": "externalService1"}
			externalServiceMap2 := map[string]interface{}{"name": "externalService2"}
			externalServiceMap3 := map[string]interface{}{"name": "externalService3"}
			externalServiceMap4 := map[string]interface{}{"name": "externalService4"}

			filepathMap1 := map[string]interface{}{"name": "filePath1"}
			filepathMap2 := map[string]interface{}{"name": "filePath2"}
			filepathMap3 := map[string]interface{}{"name": "filePath3"}
			filepathMap4 := map[string]interface{}{"name": "filePath4"}

			apiAuthorizationPolicyMap1 := map[string]interface{}{"name": "api1", "authorizedNamespace": "/3/4/5", "subject": []interface{}{[]interface{}{"$namespace=/3/4/5"}}}
			apiAuthorizationPolicyMap2 := map[string]interface{}{"name": "api2", "authorizedNamespace": "/3/4/5"}
			topNamespace := tree["namespaces"].(map[string]interface{})

			So(err, ShouldBeNil)

			So(topNamespace["name"], ShouldEqual, "3")
			So(len(topNamespace[namespaceContentKey].(map[string][]map[string]interface{})["namespaces"]), ShouldEqual, 2)
			So(len(topNamespace[namespaceContentKey].(map[string][]map[string]interface{})["externalservices"]), ShouldEqual, 2)
			So(len(topNamespace[namespaceContentKey].(map[string][]map[string]interface{})["filepaths"]), ShouldEqual, 1)
			So(len(topNamespace[namespaceContentKey].(map[string][]map[string]interface{})["apiauthorizationpolicies"]), ShouldEqual, 2)
			So(topNamespace[namespaceContentKey].(map[string][]map[string]interface{})["filepaths"][0], ShouldResemble, filepathMap1)
			So(topNamespace[namespaceContentKey].(map[string][]map[string]interface{})["namespaces"][0]["name"], ShouldResemble, namespaceMap1["name"])
			So(topNamespace[namespaceContentKey].(map[string][]map[string]interface{})["namespaces"][1]["name"], ShouldResemble, namespaceMap2["name"])
			So(topNamespace[namespaceContentKey].(map[string][]map[string]interface{})["externalservices"][0], ShouldResemble, externalServiceMap1)
			So(topNamespace[namespaceContentKey].(map[string][]map[string]interface{})["externalservices"][1], ShouldResemble, externalServiceMap2)
			So(topNamespace[namespaceContentKey].(map[string][]map[string]interface{})["apiauthorizationpolicies"][0], ShouldResemble, apiAuthorizationPolicyMap1)
			So(topNamespace[namespaceContentKey].(map[string][]map[string]interface{})["apiauthorizationpolicies"][1], ShouldResemble, apiAuthorizationPolicyMap2)

			// namespace1
			ns := topNamespace[namespaceContentKey].(map[string][]map[string]interface{})["namespaces"][0]
			So(len(ns[namespaceContentKey].(map[string][]map[string]interface{})["namespaces"]), ShouldEqual, 1)
			So(len(ns[namespaceContentKey].(map[string][]map[string]interface{})["externalservices"]), ShouldEqual, 0)
			So(len(ns[namespaceContentKey].(map[string][]map[string]interface{})["filepaths"]), ShouldEqual, 0)
			So(ns[namespaceContentKey].(map[string][]map[string]interface{})["namespaces"][0]["name"], ShouldResemble, namespaceMap3["name"])

			// namespace3
			ns = ns[namespaceContentKey].(map[string][]map[string]interface{})["namespaces"][0]
			So(len(ns[namespaceContentKey].(map[string][]map[string]interface{})["namespaces"]), ShouldEqual, 0)
			So(len(ns[namespaceContentKey].(map[string][]map[string]interface{})["externalservices"]), ShouldEqual, 1)
			So(len(ns[namespaceContentKey].(map[string][]map[string]interface{})["filepaths"]), ShouldEqual, 2)
			So(ns[namespaceContentKey].(map[string][]map[string]interface{})["filepaths"][0], ShouldResemble, filepathMap2)
			So(ns[namespaceContentKey].(map[string][]map[string]interface{})["filepaths"][1], ShouldResemble, filepathMap3)
			So(ns[namespaceContentKey].(map[string][]map[string]interface{})["externalservices"][0], ShouldResemble, externalServiceMap3)

			// namespace2
			ns = topNamespace[namespaceContentKey].(map[string][]map[string]interface{})["namespaces"][1]
			So(len(ns[namespaceContentKey].(map[string][]map[string]interface{})["namespaces"]), ShouldEqual, 0)
			So(len(ns[namespaceContentKey].(map[string][]map[string]interface{})["externalservices"]), ShouldEqual, 1)
			So(len(ns[namespaceContentKey].(map[string][]map[string]interface{})["filepaths"]), ShouldEqual, 1)
			So(ns[namespaceContentKey].(map[string][]map[string]interface{})["filepaths"][0], ShouldResemble, filepathMap4)
			So(ns[namespaceContentKey].(map[string][]map[string]interface{})["externalservices"][0], ShouldResemble, externalServiceMap4)
		})
	})
}

func Test_fillTreeNamespaceContent(t *testing.T) {
	Convey("Given test data is prepared with namespace /3", t, func() {
		topNamespace := map[string]interface{}{"name": "3"}
		namespaceContentRegistry := map[string][]map[string]map[string]interface{}{}

		namespace1 := map[string]interface{}{"name": "4"}
		namespace2 := map[string]interface{}{"name": "4.1"}
		namespace3 := map[string]interface{}{"name": "5"}

		externalService1 := map[string]interface{}{"name": "externalService1"}
		externalService2 := map[string]interface{}{"name": "externalService2"}
		externalService3 := map[string]interface{}{"name": "externalService3"}
		externalService4 := map[string]interface{}{"name": "externalService4"}

		filepath1 := map[string]interface{}{"name": "filePath1"}
		filepath2 := map[string]interface{}{"name": "filePath2"}
		filepath3 := map[string]interface{}{"name": "filePath3"}
		filepath4 := map[string]interface{}{"name": "filePath4"}

		namespaceContentRegistry["3"] = []map[string]map[string]interface{}{}
		namespaceContentRegistry["3"] = append(namespaceContentRegistry["3"], map[string]map[string]interface{}{"namespaces": namespace1})
		namespaceContentRegistry["3"] = append(namespaceContentRegistry["3"], map[string]map[string]interface{}{"namespaces": namespace2})
		namespaceContentRegistry["3"] = append(namespaceContentRegistry["3"], map[string]map[string]interface{}{"externalservices": externalService1})
		namespaceContentRegistry["3"] = append(namespaceContentRegistry["3"], map[string]map[string]interface{}{"externalservices": externalService2})
		namespaceContentRegistry["3"] = append(namespaceContentRegistry["3"], map[string]map[string]interface{}{"filepaths": filepath1})

		namespaceContentRegistry["3/4"] = []map[string]map[string]interface{}{}
		namespaceContentRegistry["3/4"] = append(namespaceContentRegistry["3/4"], map[string]map[string]interface{}{"namespaces": namespace3})

		namespaceContentRegistry["3/4.1"] = []map[string]map[string]interface{}{}
		namespaceContentRegistry["3/4.1"] = append(namespaceContentRegistry["3/4.1"], map[string]map[string]interface{}{"externalservices": externalService4})
		namespaceContentRegistry["3/4.1"] = append(namespaceContentRegistry["3/4.1"], map[string]map[string]interface{}{"filepaths": filepath4})

		namespaceContentRegistry["3/4/5"] = []map[string]map[string]interface{}{}
		namespaceContentRegistry["3/4/5"] = append(namespaceContentRegistry["3/4/5"], map[string]map[string]interface{}{"externalservices": externalService3})
		namespaceContentRegistry["3/4/5"] = append(namespaceContentRegistry["3/4/5"], map[string]map[string]interface{}{"filepaths": filepath2})
		namespaceContentRegistry["3/4/5"] = append(namespaceContentRegistry["3/4/5"], map[string]map[string]interface{}{"filepaths": filepath3})

		Convey("Then I fill my top namespace with the data", func() {
			fillTreeForNamespace("", topNamespace, namespaceContentRegistry)

			So(topNamespace["name"], ShouldEqual, "3")
			So(len(topNamespace[namespaceContentKey].(map[string][]map[string]interface{})["namespaces"]), ShouldEqual, 2)
			So(len(topNamespace[namespaceContentKey].(map[string][]map[string]interface{})["externalservices"]), ShouldEqual, 2)
			So(len(topNamespace[namespaceContentKey].(map[string][]map[string]interface{})["filepaths"]), ShouldEqual, 1)
			So(topNamespace[namespaceContentKey].(map[string][]map[string]interface{})["filepaths"][0], ShouldResemble, filepath1)
			So(topNamespace[namespaceContentKey].(map[string][]map[string]interface{})["namespaces"][0], ShouldResemble, namespace1)
			So(topNamespace[namespaceContentKey].(map[string][]map[string]interface{})["namespaces"][1], ShouldResemble, namespace2)
			So(topNamespace[namespaceContentKey].(map[string][]map[string]interface{})["externalservices"][0], ShouldResemble, externalService1)
			So(topNamespace[namespaceContentKey].(map[string][]map[string]interface{})["externalservices"][1], ShouldResemble, externalService2)

			// namespace1
			ns := topNamespace[namespaceContentKey].(map[string][]map[string]interface{})["namespaces"][0]
			So(len(ns[namespaceContentKey].(map[string][]map[string]interface{})["namespaces"]), ShouldEqual, 1)
			So(len(ns[namespaceContentKey].(map[string][]map[string]interface{})["externalservices"]), ShouldEqual, 0)
			So(len(ns[namespaceContentKey].(map[string][]map[string]interface{})["filepaths"]), ShouldEqual, 0)
			So(ns[namespaceContentKey].(map[string][]map[string]interface{})["namespaces"][0], ShouldResemble, namespace3)

			// namespace3
			ns = ns[namespaceContentKey].(map[string][]map[string]interface{})["namespaces"][0]
			So(len(ns[namespaceContentKey].(map[string][]map[string]interface{})["namespaces"]), ShouldEqual, 0)
			So(len(ns[namespaceContentKey].(map[string][]map[string]interface{})["externalservices"]), ShouldEqual, 1)
			So(len(ns[namespaceContentKey].(map[string][]map[string]interface{})["filepaths"]), ShouldEqual, 2)
			So(ns[namespaceContentKey].(map[string][]map[string]interface{})["filepaths"][0], ShouldResemble, filepath2)
			So(ns[namespaceContentKey].(map[string][]map[string]interface{})["filepaths"][1], ShouldResemble, filepath3)
			So(ns[namespaceContentKey].(map[string][]map[string]interface{})["externalservices"][0], ShouldResemble, externalService3)

			// namespace2
			ns = topNamespace[namespaceContentKey].(map[string][]map[string]interface{})["namespaces"][1]
			So(len(ns[namespaceContentKey].(map[string][]map[string]interface{})["namespaces"]), ShouldEqual, 0)
			So(len(ns[namespaceContentKey].(map[string][]map[string]interface{})["externalservices"]), ShouldEqual, 1)
			So(len(ns[namespaceContentKey].(map[string][]map[string]interface{})["filepaths"]), ShouldEqual, 1)
			So(ns[namespaceContentKey].(map[string][]map[string]interface{})["filepaths"][0], ShouldResemble, filepath4)
			So(ns[namespaceContentKey].(map[string][]map[string]interface{})["externalservices"][0], ShouldResemble, externalService4)
		})
	})

	Convey("Given test data is prepared with namespace /", t, func() {
		topNamespace := map[string]interface{}{"name": ""}
		namespaceContentRegistry := map[string][]map[string]map[string]interface{}{}

		namespace1 := map[string]interface{}{"name": "4"}
		namespace2 := map[string]interface{}{"name": "4.1"}
		namespace3 := map[string]interface{}{"name": "5"}

		externalService1 := map[string]interface{}{"name": "externalService1"}
		externalService2 := map[string]interface{}{"name": "externalService2"}
		externalService3 := map[string]interface{}{"name": "externalService3"}
		externalService4 := map[string]interface{}{"name": "externalService4"}

		filepath1 := map[string]interface{}{"name": "filePath1"}
		filepath2 := map[string]interface{}{"name": "filePath2"}
		filepath3 := map[string]interface{}{"name": "filePath3"}
		filepath4 := map[string]interface{}{"name": "filePath4"}

		namespaceContentRegistry[""] = []map[string]map[string]interface{}{}
		namespaceContentRegistry[""] = append(namespaceContentRegistry[""], map[string]map[string]interface{}{"namespaces": namespace1})
		namespaceContentRegistry[""] = append(namespaceContentRegistry[""], map[string]map[string]interface{}{"namespaces": namespace2})
		namespaceContentRegistry[""] = append(namespaceContentRegistry[""], map[string]map[string]interface{}{"externalservices": externalService1})
		namespaceContentRegistry[""] = append(namespaceContentRegistry[""], map[string]map[string]interface{}{"externalservices": externalService2})
		namespaceContentRegistry[""] = append(namespaceContentRegistry[""], map[string]map[string]interface{}{"filepaths": filepath1})

		namespaceContentRegistry["4"] = []map[string]map[string]interface{}{}
		namespaceContentRegistry["4"] = append(namespaceContentRegistry["4"], map[string]map[string]interface{}{"namespaces": namespace3})

		namespaceContentRegistry["4.1"] = []map[string]map[string]interface{}{}
		namespaceContentRegistry["4.1"] = append(namespaceContentRegistry["4.1"], map[string]map[string]interface{}{"externalservices": externalService4})
		namespaceContentRegistry["4.1"] = append(namespaceContentRegistry["4.1"], map[string]map[string]interface{}{"filepaths": filepath4})

		namespaceContentRegistry["4/5"] = []map[string]map[string]interface{}{}
		namespaceContentRegistry["4/5"] = append(namespaceContentRegistry["4/5"], map[string]map[string]interface{}{"externalservices": externalService3})
		namespaceContentRegistry["4/5"] = append(namespaceContentRegistry["4/5"], map[string]map[string]interface{}{"filepaths": filepath2})
		namespaceContentRegistry["4/5"] = append(namespaceContentRegistry["4/5"], map[string]map[string]interface{}{"filepaths": filepath3})

		Convey("Then I fill my top namespace with the data", func() {
			fillTreeForNamespace("", topNamespace, namespaceContentRegistry)

			So(topNamespace["name"], ShouldEqual, "")
			So(len(topNamespace[namespaceContentKey].(map[string][]map[string]interface{})["namespaces"]), ShouldEqual, 2)
			So(len(topNamespace[namespaceContentKey].(map[string][]map[string]interface{})["externalservices"]), ShouldEqual, 2)
			So(len(topNamespace[namespaceContentKey].(map[string][]map[string]interface{})["filepaths"]), ShouldEqual, 1)
			So(topNamespace[namespaceContentKey].(map[string][]map[string]interface{})["filepaths"][0], ShouldResemble, filepath1)
			So(topNamespace[namespaceContentKey].(map[string][]map[string]interface{})["namespaces"][0], ShouldResemble, namespace1)
			So(topNamespace[namespaceContentKey].(map[string][]map[string]interface{})["namespaces"][1], ShouldResemble, namespace2)
			So(topNamespace[namespaceContentKey].(map[string][]map[string]interface{})["externalservices"][0], ShouldResemble, externalService1)
			So(topNamespace[namespaceContentKey].(map[string][]map[string]interface{})["externalservices"][1], ShouldResemble, externalService2)

			// namespace1
			ns := topNamespace[namespaceContentKey].(map[string][]map[string]interface{})["namespaces"][0]
			So(len(ns[namespaceContentKey].(map[string][]map[string]interface{})["namespaces"]), ShouldEqual, 1)
			So(len(ns[namespaceContentKey].(map[string][]map[string]interface{})["externalservices"]), ShouldEqual, 0)
			So(len(ns[namespaceContentKey].(map[string][]map[string]interface{})["filepaths"]), ShouldEqual, 0)
			So(ns[namespaceContentKey].(map[string][]map[string]interface{})["namespaces"][0], ShouldResemble, namespace3)

			// namespace3
			ns = ns[namespaceContentKey].(map[string][]map[string]interface{})["namespaces"][0]
			So(len(ns[namespaceContentKey].(map[string][]map[string]interface{})["namespaces"]), ShouldEqual, 0)
			So(len(ns[namespaceContentKey].(map[string][]map[string]interface{})["externalservices"]), ShouldEqual, 1)
			So(len(ns[namespaceContentKey].(map[string][]map[string]interface{})["filepaths"]), ShouldEqual, 2)
			So(ns[namespaceContentKey].(map[string][]map[string]interface{})["filepaths"][0], ShouldResemble, filepath2)
			So(ns[namespaceContentKey].(map[string][]map[string]interface{})["filepaths"][1], ShouldResemble, filepath3)
			So(ns[namespaceContentKey].(map[string][]map[string]interface{})["externalservices"][0], ShouldResemble, externalService3)

			// namespace2
			ns = topNamespace[namespaceContentKey].(map[string][]map[string]interface{})["namespaces"][1]
			So(len(ns[namespaceContentKey].(map[string][]map[string]interface{})["namespaces"]), ShouldEqual, 0)
			So(len(ns[namespaceContentKey].(map[string][]map[string]interface{})["externalservices"]), ShouldEqual, 1)
			So(len(ns[namespaceContentKey].(map[string][]map[string]interface{})["filepaths"]), ShouldEqual, 1)
			So(ns[namespaceContentKey].(map[string][]map[string]interface{})["filepaths"][0], ShouldResemble, filepath4)
			So(ns[namespaceContentKey].(map[string][]map[string]interface{})["externalservices"][0], ShouldResemble, externalService4)
		})
	})
}

func Test_computeNamespaceAttributes(t *testing.T) {
	Convey("Given test data is prepared", t, func() {
		apiAuthorizationPolicy1 := map[string]interface{}{"authorizedNamespace": "/1/2/3/4"}
		apiAuthorizationPolicy2 := map[string]interface{}{"authorizedNamespace": "/1/2/3"}

		namespaceMappingPolicy1 := map[string]interface{}{"mappedNamespace": "/1/2/3/4", "object": []interface{}{[]interface{}{"$namespace=/1/2/3"}, []interface{}{"$namespace=/1/2/3/4"}, []interface{}{"$namespace=/1/2/3/5"}}}
		namespaceMappingPolicy2 := map[string]interface{}{"mappedNamespace": "/1/2/3", "subject": []interface{}{[]interface{}{"$namespace=/1/2/3"}, []interface{}{"$namespace=/1/2/3/4"}, []interface{}{"$namespace=/1/2/3/5"}}}

		namespace1 := map[string]interface{}{"name": "/1/2/3/4"}
		namespace2 := map[string]interface{}{"name": "/1/2/3"}

		Convey("Then I try to compute the data with namespace /1/2", func() {
			namespace := "/1/2"
			exportComputeNamespaceAttributes(namespace, squallmodels.APIAuthorizationPolicyIdentity.Category, apiAuthorizationPolicy1)
			exportComputeNamespaceAttributes(namespace, squallmodels.APIAuthorizationPolicyIdentity.Category, apiAuthorizationPolicy2)
			exportComputeNamespaceAttributes(namespace, squallmodels.NamespaceMappingPolicyIdentity.Category, namespaceMappingPolicy1)
			exportComputeNamespaceAttributes(namespace, squallmodels.NamespaceMappingPolicyIdentity.Category, namespaceMappingPolicy2)
			exportComputeNamespaceAttributes(namespace, squallmodels.NamespaceIdentity.Category, namespace1)
			exportComputeNamespaceAttributes(namespace, squallmodels.NamespaceIdentity.Category, namespace2)

			So(apiAuthorizationPolicy1["authorizedNamespace"], ShouldEqual, "/2/3/4")
			So(apiAuthorizationPolicy2["authorizedNamespace"], ShouldEqual, "/2/3")
			So(namespaceMappingPolicy1["mappedNamespace"], ShouldEqual, "/2/3/4")
			So(namespaceMappingPolicy2["mappedNamespace"], ShouldEqual, "/2/3")
			So(namespaceMappingPolicy1["object"], ShouldResemble, []interface{}{[]interface{}{"$namespace=/2/3"}, []interface{}{"$namespace=/2/3/4"}, []interface{}{"$namespace=/2/3/5"}})
			So(namespaceMappingPolicy2["subject"], ShouldResemble, []interface{}{[]interface{}{"$namespace=/2/3"}, []interface{}{"$namespace=/2/3/4"}, []interface{}{"$namespace=/2/3/5"}})
			So(namespace1["name"], ShouldEqual, "4")
			So(namespace2["name"], ShouldEqual, "3")
		})

		Convey("Then I try to compute the data with namespace /1/2/3", func() {
			namespace := "/1/2/3"
			exportComputeNamespaceAttributes(namespace, squallmodels.APIAuthorizationPolicyIdentity.Category, apiAuthorizationPolicy1)
			exportComputeNamespaceAttributes(namespace, squallmodels.APIAuthorizationPolicyIdentity.Category, apiAuthorizationPolicy2)
			exportComputeNamespaceAttributes(namespace, squallmodels.NamespaceMappingPolicyIdentity.Category, namespaceMappingPolicy1)
			exportComputeNamespaceAttributes(namespace, squallmodels.NamespaceMappingPolicyIdentity.Category, namespaceMappingPolicy2)
			exportComputeNamespaceAttributes(namespace, squallmodels.NamespaceIdentity.Category, namespace1)
			exportComputeNamespaceAttributes(namespace, squallmodels.NamespaceIdentity.Category, namespace2)

			So(apiAuthorizationPolicy1["authorizedNamespace"], ShouldEqual, "/3/4")
			So(apiAuthorizationPolicy2["authorizedNamespace"], ShouldEqual, "/3")
			So(namespaceMappingPolicy1["mappedNamespace"], ShouldEqual, "/3/4")
			So(namespaceMappingPolicy2["mappedNamespace"], ShouldEqual, "/3")
			So(namespaceMappingPolicy1["object"], ShouldResemble, []interface{}{[]interface{}{"$namespace=/3"}, []interface{}{"$namespace=/3/4"}, []interface{}{"$namespace=/3/5"}})
			So(namespaceMappingPolicy2["subject"], ShouldResemble, []interface{}{[]interface{}{"$namespace=/3"}, []interface{}{"$namespace=/3/4"}, []interface{}{"$namespace=/3/5"}})
			So(namespace1["name"], ShouldEqual, "4")
			So(namespace2["name"], ShouldEqual, "3")
		})
	})
}
