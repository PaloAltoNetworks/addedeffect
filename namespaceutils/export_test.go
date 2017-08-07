package namespaceutils

import (
	"testing"

	"github.com/aporeto-inc/elemental"
	"github.com/aporeto-inc/gaia/squallmodels/v1/golang"
	"github.com/aporeto-inc/manipulate"
	"github.com/aporeto-inc/manipulate/maniptest"
	. "github.com/smartystreets/goconvey/convey"
)

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
			tree, err := TreeContentOfNamespace(namespace, identifiables, "")
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

			apiAuthorizationPolicyMap1 := map[string]interface{}{"activeSchedule": "", "name": "api1", "authorizedNamespace": "/3/4/5", "subject": []interface{}{[]interface{}{"$namespace=/3/4/5"}}}
			apiAuthorizationPolicyMap2 := map[string]interface{}{"activeSchedule": "", "name": "api2", "authorizedNamespace": "/3/4/5"}
			topNamespace := tree["namespaces"].([]interface{})[0].(map[string]interface{})

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

func Test_ContentOfNamespace(t *testing.T) {
	Convey("Given test data is prepared", t, func() {

		manipulator := maniptest.NewTestManipulator()

		namespaceMappingPolicy1 := squallmodels.NewNamespaceMappingPolicy()
		namespaceMappingPolicy1.Name = "namespaceMappingPolicy1"

		namespaceMappingPolicy2 := squallmodels.NewNamespaceMappingPolicy()
		namespaceMappingPolicy2.Name = "namespaceMappingPolicy2"

		networksAccessPolicy1 := squallmodels.NewNetworkAccessPolicy()
		networksAccessPolicy1.Name = "networksAccessPolicy1"

		fileAccessPolicy1 := squallmodels.NewFileAccessPolicy()
		fileAccessPolicy1.Name = "fileAccessPolicy1"

		fileAccessPolicy2 := squallmodels.NewFileAccessPolicy()
		fileAccessPolicy2.Name = "fileAccessPolicy2"

		enforcerMappingPolicy1 := squallmodels.NewEnforcerProfileMappingPolicy()
		enforcerMappingPolicy1.Name = "enforcerMappingPolicy1"

		namespace1 := squallmodels.NewNamespace()
		namespace1.Name = "namespace1"

		externalService1 := squallmodels.NewExternalService()
		externalService1.Name = "externalService1"

		filePath1 := squallmodels.NewFilePath()
		filePath1.Name = "filePath1"

		enforcerProfile1 := squallmodels.NewEnforcerProfile()
		enforcerProfile1.Name = "enforcerProfile1"

		Convey("Then we get an error", func() {

			manipulator.MockRetrieveMany(t, func(context *manipulate.Context, dest elemental.ContentIdentifiable) error {

				if dest.ContentIdentity().Name == squallmodels.NamespaceMappingPolicyIdentity.Name {
					policies := dest.(*squallmodels.NamespaceMappingPoliciesList)
					*policies = append(*policies, namespaceMappingPolicy2, namespaceMappingPolicy1)
					dest = policies
					_ = dest
				}

				if dest.ContentIdentity().Name == squallmodels.NetworkAccessPolicyIdentity.Name {
					policies := dest.(*squallmodels.NetworkAccessPoliciesList)
					*policies = append(*policies, networksAccessPolicy1)
					dest = policies
					_ = dest
				}

				if dest.ContentIdentity().Name == squallmodels.FileAccessPolicyIdentity.Name {
					policies := dest.(*squallmodels.FileAccessPoliciesList)
					*policies = append(*policies, fileAccessPolicy1, fileAccessPolicy2)
					dest = policies
					_ = dest
				}

				if dest.ContentIdentity().Name == squallmodels.EnforcerProfileMappingPolicyIdentity.Name {
					policies := dest.(*squallmodels.EnforcerProfileMappingPoliciesList)
					*policies = append(*policies, enforcerMappingPolicy1)
					dest = policies
					_ = dest
				}

				if dest.ContentIdentity().Name == squallmodels.NamespaceIdentity.Name {
					namespaces := dest.(*squallmodels.NamespacesList)
					*namespaces = append(*namespaces, namespace1)
					dest = namespaces
					_ = dest
				}

				if dest.ContentIdentity().Name == squallmodels.ExternalServiceIdentity.Name {
					externalServices := dest.(*squallmodels.ExternalServicesList)
					*externalServices = append(*externalServices, externalService1)
					dest = externalServices
					_ = dest
				}

				if dest.ContentIdentity().Name == squallmodels.FilePathIdentity.Name {
					filePaths := dest.(*squallmodels.FilePathsList)
					*filePaths = append(*filePaths, filePath1)
					dest = filePaths
					_ = dest
				}

				if dest.ContentIdentity().Name == squallmodels.EnforcerProfileIdentity.Name {
					enforcers := dest.(*squallmodels.EnforcerProfilesList)
					*enforcers = append(*enforcers, enforcerProfile1)
					dest = enforcers
					_ = dest
				}

				return nil
			})

			// FIXME !!: content, err := ContentOfNamespace(manipulator, "/coucou", true, "")
			// FIXME !!: So(content, ShouldBeNil)
			// FIXME !!: So(err, ShouldNotBeNil)
		})

		Convey("Then we get data from the namespaces", func() {

			var expectedNamespace string
			var expectedRecursive bool

			manipulator.MockRetrieveMany(t, func(context *manipulate.Context, dest elemental.ContentIdentifiable) error {

				expectedNamespace = context.Namespace
				expectedRecursive = context.Recursive

				if dest.ContentIdentity().Name == squallmodels.NamespaceMappingPolicyIdentity.Name {
					policies := dest.(*squallmodels.NamespaceMappingPoliciesList)
					*policies = append(*policies, namespaceMappingPolicy2, namespaceMappingPolicy1)
					dest = policies
					_ = dest
				}

				if dest.ContentIdentity().Name == squallmodels.NetworkAccessPolicyIdentity.Name {
					policies := dest.(*squallmodels.NetworkAccessPoliciesList)
					*policies = append(*policies, networksAccessPolicy1)
					dest = policies
					_ = dest
				}

				if dest.ContentIdentity().Name == squallmodels.FileAccessPolicyIdentity.Name {
					policies := dest.(*squallmodels.FileAccessPoliciesList)
					*policies = append(*policies, fileAccessPolicy1, fileAccessPolicy2)
					dest = policies
					_ = dest
				}

				if dest.ContentIdentity().Name == squallmodels.EnforcerProfileMappingPolicyIdentity.Name {
					policies := dest.(*squallmodels.EnforcerProfileMappingPoliciesList)
					*policies = append(*policies, enforcerMappingPolicy1)
					dest = policies
					_ = dest
				}

				if dest.ContentIdentity().Name == squallmodels.NamespaceIdentity.Name {
					namespaces := dest.(*squallmodels.NamespacesList)
					*namespaces = append(*namespaces, namespace1)
					dest = namespaces
					_ = dest
				}

				if dest.ContentIdentity().Name == squallmodels.ExternalServiceIdentity.Name {
					externalServices := dest.(*squallmodels.ExternalServicesList)
					*externalServices = append(*externalServices, externalService1)
					dest = externalServices
					_ = dest
				}

				if dest.ContentIdentity().Name == squallmodels.FilePathIdentity.Name {
					filePaths := dest.(*squallmodels.FilePathsList)
					*filePaths = append(*filePaths, filePath1)
					dest = filePaths
					_ = dest
				}

				if dest.ContentIdentity().Name == squallmodels.EnforcerProfileIdentity.Name {
					enforcers := dest.(*squallmodels.EnforcerProfilesList)
					*enforcers = append(*enforcers, enforcerProfile1)
					dest = enforcers
					_ = dest
				}

				return nil
			})

			content, err := ContentOfNamespace(manipulator, "/coucou", true, "")
			So(err, ShouldBeNil)
			So(len(content), ShouldEqual, 10)
			So(content, ShouldContain, namespaceMappingPolicy1)
			So(content, ShouldContain, namespaceMappingPolicy2)
			So(content, ShouldContain, networksAccessPolicy1)
			So(content, ShouldContain, fileAccessPolicy1)
			So(content, ShouldContain, fileAccessPolicy2)
			So(content, ShouldContain, enforcerMappingPolicy1)
			So(content, ShouldContain, namespace1)
			So(content, ShouldContain, externalService1)
			So(content, ShouldContain, filePath1)
			So(content, ShouldContain, enforcerProfile1)

			So(expectedNamespace, ShouldEqual, "/coucou")
			So(expectedRecursive, ShouldBeTrue)
		})
	})
}
