package namespaceutils

import (
	"testing"

	"github.com/aporeto-inc/elemental"
	squallmodels "github.com/aporeto-inc/gaia/squallmodels/current/golang"
	"github.com/aporeto-inc/manipulate"
	"github.com/aporeto-inc/manipulate/maniptest"
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

func Test_createNamespace(t *testing.T) {
	Convey("Given test data is prepared", t, func() {
		manipulator := maniptest.NewTestManipulator()
		namespace := squallmodels.NewNamespace()
		namespace.Name = "coucou"

		Convey("Given the creation of the namespace is a success", func() {

			var expectedNamespace *squallmodels.Namespace
			var expectedNamespaceSession string

			manipulator.MockCreate(t, func(ctx *manipulate.Context, objects ...elemental.Identifiable) error {
				expectedNamespaceSession = ctx.Namespace
				expectedNamespace = objects[0].(*squallmodels.Namespace)
				return nil
			})

			err := createNamespace(manipulator, "/2", namespace)
			So(err, ShouldBeNil)
			So(expectedNamespace, ShouldEqual, namespace)
			So(expectedNamespaceSession, ShouldEqual, "/2")
		})

		Convey("Given the creation of the namespace is a failure", func() {
			manipulator.MockCreate(t, func(ctx *manipulate.Context, objects ...elemental.Identifiable) error {
				return elemental.NewError("Invalid Entity", "", "", 500)
			})

			err := createNamespace(manipulator, "/2", namespace)
			So(err, ShouldNotBeNil)
		})
	})
}

func Test_deleteNamespace(t *testing.T) {
	Convey("Given test data is prepared", t, func() {
		manipulator := maniptest.NewTestManipulator()
		namespace := squallmodels.NewNamespace()
		namespace.Name = "coucou"

		Convey("Given the delete of the namespace is a success", func() {

			var expectedNamespace *squallmodels.Namespace
			var expectedNamespaceSession string

			manipulator.MockDelete(t, func(ctx *manipulate.Context, objects ...elemental.Identifiable) error {
				expectedNamespaceSession = ctx.Namespace
				expectedNamespace = objects[0].(*squallmodels.Namespace)
				return nil
			})

			err := deleteNamespace(manipulator, "/2", namespace)
			So(err, ShouldBeNil)
			So(expectedNamespace, ShouldEqual, namespace)
			So(expectedNamespaceSession, ShouldEqual, "/2")
		})

		Convey("Given the delete of the namespace is a failure", func() {
			manipulator.MockDelete(t, func(ctx *manipulate.Context, objects ...elemental.Identifiable) error {
				return elemental.NewError("Invalid Entity", "", "", 500)
			})

			err := deleteNamespace(manipulator, "/2", namespace)
			So(err, ShouldNotBeNil)
		})
	})
}

func Test_deleteContent(t *testing.T) {
	Convey("Given test data is prepared", t, func() {
		manipulator := maniptest.NewTestManipulator()

		namespace1 := &squallmodels.Namespace{Name: "5", Namespace: "/3/4"}

		externalService1 := &squallmodels.ExternalService{Name: "externalService1", Namespace: "/3"}
		externalService2 := &squallmodels.ExternalService{Name: "externalService2", Namespace: "/3"}

		filepath1 := &squallmodels.FilePath{Name: "filePath1", Namespace: "/3"}

		content := elemental.IdentifiablesList{namespace1, externalService1, externalService2, filepath1}

		Convey("Given the delete of the content is a success", func() {

			var expectedNamespaceSession string
			expectedContent := elemental.IdentifiablesList{}

			manipulator.MockDelete(t, func(ctx *manipulate.Context, objects ...elemental.Identifiable) error {
				expectedNamespaceSession = ctx.Namespace
				expectedContent = append(expectedContent, objects[0])
				return nil
			})

			err := deleteContent(manipulator, "/2", content)
			So(err, ShouldBeNil)
			So(len(expectedContent), ShouldEqual, 3)
			So(expectedContent, ShouldContain, externalService1)
			So(expectedContent, ShouldContain, externalService2)
			So(expectedContent, ShouldContain, filepath1)
			So(expectedNamespaceSession, ShouldEqual, "/2")
		})

		Convey("Given the delete of the content is a failure", func() {
			manipulator.MockDelete(t, func(ctx *manipulate.Context, objects ...elemental.Identifiable) error {
				return elemental.NewError("Invalid Entity", "", "", 500)
			})

			err := deleteContent(manipulator, "/2", elemental.IdentifiablesList{externalService1})
			So(err, ShouldNotBeNil)
		})
	})
}

func Test_isNamespaceExist(t *testing.T) {
	Convey("Given test data is prepared", t, func() {

		manipulator := maniptest.NewTestManipulator()

		Convey("Given we have an existing namespace", func() {

			namespace := squallmodels.NewNamespace()
			namespace.Name = "/a/c"

			manipulator.MockRetrieveMany(t, func(context *manipulate.Context, dest elemental.ContentIdentifiable) error {
				namespaces := dest.(*squallmodels.NamespacesList)
				*namespaces = append(*namespaces, namespace)
				dest = namespaces
				_ = dest
				return nil
			})

			result, err := isNamespaceExist(manipulator, "/", namespace)
			So(err, ShouldBeNil)
			So(result, ShouldBeTrue)
		})

		Convey("Given we don't have an existing namespace", func() {

			namespace := squallmodels.NewNamespace()
			namespace.Name = "/a/b"

			manipulator.MockRetrieveMany(t, func(context *manipulate.Context, dest elemental.ContentIdentifiable) error {
				namespaces := dest.(*squallmodels.NamespacesList)
				dest = namespaces
				_ = dest
				return nil
			})

			result, err := isNamespaceExist(manipulator, "/", namespace)
			So(err, ShouldBeNil)
			So(result, ShouldBeFalse)
		})

		Convey("Given we get an error from manipulate", func() {

			namespace := squallmodels.NewNamespace()
			namespace.Name = "/a/b"

			manipulator.MockRetrieveMany(t, func(context *manipulate.Context, dest elemental.ContentIdentifiable) error {
				return elemental.NewError("Invalid Entity", "", "", 500)
			})

			result, err := isNamespaceExist(manipulator, "/", namespace)
			So(err, ShouldNotBeNil)
			So(result, ShouldBeFalse)
		})
	})
}

func Test_createContent(t *testing.T) {
	Convey("Given test data is prepared", t, func() {
		topNamespace := "/coucou"
		namespace := "/coucou/java"

		manipulator := maniptest.NewTestManipulator()

		namespace1 := map[string]interface{}{"name": "5"}
		externalService1 := map[string]interface{}{"name": "externalService1"}
		externalService2 := map[string]interface{}{"name": "externalService2"}
		filepath1 := map[string]interface{}{"name": "filePath1"}
		apiAuthorizationPolicy1 := map[string]interface{}{"name": "apiAuthorizationPolicy1", "authorizedNamespace": "/java/3", "subject": []interface{}{[]interface{}{"$namespace=/java/4/5"}}}

		content := make(map[string]interface{})
		content[squallmodels.APIAuthorizationPolicyIdentity.Category] = []interface{}{apiAuthorizationPolicy1}
		content[squallmodels.ExternalServiceIdentity.Category] = []interface{}{externalService1, externalService2}
		content[squallmodels.FilePathIdentity.Category] = []interface{}{filepath1}
		content[squallmodels.NamespaceIdentity.Category] = []interface{}{namespace1}

		Convey("Given we the creation of the content is a success", func() {

			var expectedNamespaceSession string
			expectedContent := elemental.IdentifiablesList{}
			var expectedExternalService1 *squallmodels.ExternalService
			var expectedExternalService2 *squallmodels.ExternalService
			var expectedFilePath *squallmodels.FilePath
			var expectedAPIAuthorization *squallmodels.APIAuthorizationPolicy

			manipulator.MockCreate(t, func(ctx *manipulate.Context, objects ...elemental.Identifiable) error {
				expectedNamespaceSession = ctx.Namespace
				expectedContent = append(expectedContent, objects[0])

				if objects[0].Identity().Name == squallmodels.ExternalServiceIdentity.Name {

					if expectedExternalService1 == nil {
						expectedExternalService1 = objects[0].(*squallmodels.ExternalService)
					} else {
						expectedExternalService2 = objects[0].(*squallmodels.ExternalService)
					}

				}

				if objects[0].Identity().Name == squallmodels.FilePathIdentity.Name {
					expectedFilePath = objects[0].(*squallmodels.FilePath)
				}

				if objects[0].Identity().Name == squallmodels.APIAuthorizationPolicyIdentity.Name {
					expectedAPIAuthorization = objects[0].(*squallmodels.APIAuthorizationPolicy)
				}

				return nil
			})

			err := createContent(manipulator, topNamespace, namespace, content)
			So(err, ShouldBeNil)
			So(expectedNamespaceSession, ShouldEqual, "/coucou/java")
			So(len(expectedContent), ShouldEqual, 4)
			So(expectedExternalService1.Name, ShouldEqual, "externalService1")
			So(expectedExternalService2.Name, ShouldEqual, "externalService2")
			So(expectedFilePath.Name, ShouldEqual, "filePath1")
			So(expectedAPIAuthorization.Name, ShouldEqual, "apiAuthorizationPolicy1")
			So(expectedAPIAuthorization.AuthorizedNamespace, ShouldEqual, "/coucou/java/3")
			So(expectedAPIAuthorization.Subject[0][0], ShouldEqual, "$namespace=/coucou/java/4/5")
		})

		Convey("Given we the creation of the content is a failure", func() {

			manipulator.MockCreate(t, func(ctx *manipulate.Context, objects ...elemental.Identifiable) error {

				if objects[0].Identity().Name == squallmodels.FilePathIdentity.Name {
					return elemental.NewError("Invalid Entity", "", "", 500)
				}

				return nil
			})

			err := createContent(manipulator, topNamespace, namespace, content)
			So(err, ShouldNotBeNil)
		})
	})
}

var c *testing.T

func Test_importNamespaceContent(t *testing.T) {
	Convey("Given test data is prepared", t, func() {
		c = t
		manipulator := maniptest.NewTestManipulator()

		content := make(map[string]interface{})

		topNamespace := make(map[string]interface{})
		topNamespace["name"] = "apomux"

		topNamespaceContent := make(map[string]interface{})

		topNamespaceContent["externalservices"] = []interface{}{map[string]interface{}{"name": "externalService1"}, map[string]interface{}{"name": "externalService2"}}
		topNamespaceContent["filepaths"] = []interface{}{map[string]interface{}{"name": "filepath1"}}
		topNamespaceContent["apiauthorizationpolicies"] = []interface{}{map[string]interface{}{"name": "apiAuthorizationPolicy1", "authorizedNamespace": "/apomux/production", "subject": []interface{}{[]interface{}{"$namespace=/apomux/production"}}}}
		topNamespaceContent["namespacemappingpolicies"] = []interface{}{map[string]interface{}{"mappedNamespace": "/apomux/test", "subject": []interface{}{[]interface{}{"$namespace=/apomux/production/aporeto"}, []interface{}{"$namespace=/apomux/test"}}}}

		namespaceTest := map[string]interface{}{"name": "test", "content": make(map[string]interface{})}
		namespaceProductionAporeto := map[string]interface{}{"name": "aporeto", "content": make(map[string]interface{})}

		namespaceProductionContent := make(map[string]interface{})

		namespaceProductionContent["namespaces"] = []interface{}{namespaceProductionAporeto}
		namespaceProductionContent["filepaths"] = []interface{}{map[string]interface{}{"name": "filepath3"}}
		namespaceProductionContent["namespacemappingpolicies"] = []interface{}{map[string]interface{}{"mappedNamespace": "/apomux/production/aporeto", "subject": []interface{}{[]interface{}{"$namespace=/apomux/production/aporeto"}}}}

		namespaceProduction := map[string]interface{}{"name": "production", "content": namespaceProductionContent}
		topNamespaceContent["namespaces"] = []interface{}{namespaceTest, namespaceProduction}

		topNamespace["content"] = topNamespaceContent
		content["namespaces"] = []interface{}{topNamespace}

		Convey("Given importNamespaceContent is a success with the namespace /", func() {
			namespacesCreated := squallmodels.NamespacesList{}
			externalServicesCreated := squallmodels.ExternalServicesList{}
			filePathscreated := squallmodels.FilePathsList{}
			namespaceMappingPoliciesCreated := squallmodels.NamespaceMappingPoliciesList{}
			apiAuthorizationPoliciesCreated := squallmodels.APIAuthorizationPoliciesList{}

			manipulator.MockRetrieveMany(t, func(context *manipulate.Context, dest elemental.ContentIdentifiable) error {

				if dest.ContentIdentity().Name == squallmodels.NamespaceIdentity.Name {
					namespaces := dest.(*squallmodels.NamespacesList)
					dest = namespaces
					_ = dest
					return nil
				}
				return nil
			})

			manipulator.MockCreate(t, func(ctx *manipulate.Context, objects ...elemental.Identifiable) error {

				if objects[0].Identity().Name == squallmodels.NamespaceIdentity.Name {
					namespace := objects[0].(*squallmodels.Namespace)
					if ctx.Namespace == "/" {
						namespace.Name = "/" + namespace.Name
					} else {
						namespace.Name = ctx.Namespace + "/" + namespace.Name
					}
					namespace.Namespace = ctx.Namespace
					namespacesCreated = append(namespacesCreated, namespace)
				}

				if objects[0].Identity().Name == squallmodels.ExternalServiceIdentity.Name {
					externalService := objects[0].(*squallmodels.ExternalService)
					externalService.Namespace = ctx.Namespace
					externalServicesCreated = append(externalServicesCreated, externalService)
				}

				if objects[0].Identity().Name == squallmodels.FilePathIdentity.Name {
					filePath := objects[0].(*squallmodels.FilePath)
					filePath.Namespace = ctx.Namespace
					filePathscreated = append(filePathscreated, filePath)
				}

				if objects[0].Identity().Name == squallmodels.NamespaceMappingPolicyIdentity.Name {
					namespaceMappingPolicy := objects[0].(*squallmodels.NamespaceMappingPolicy)
					namespaceMappingPolicy.Namespace = ctx.Namespace
					namespaceMappingPoliciesCreated = append(namespaceMappingPoliciesCreated, namespaceMappingPolicy)
				}

				if objects[0].Identity().Name == squallmodels.APIAuthorizationPolicyIdentity.Name {
					apiAuthorizationPolicy := objects[0].(*squallmodels.APIAuthorizationPolicy)
					apiAuthorizationPolicy.Namespace = ctx.Namespace
					apiAuthorizationPoliciesCreated = append(apiAuthorizationPoliciesCreated, apiAuthorizationPolicy)
				}

				return nil
			})

			err := Import(manipulator, "/", content, false)

			So(err, ShouldBeNil)
			So(len(namespacesCreated), ShouldEqual, 4)
			So(len(externalServicesCreated), ShouldEqual, 2)
			So(len(filePathscreated), ShouldEqual, 2)
			So(len(apiAuthorizationPoliciesCreated), ShouldEqual, 1)
			So(len(namespaceMappingPoliciesCreated), ShouldEqual, 2)

			So(apiAuthorizationPoliciesCreated[0].Name, ShouldEqual, "apiAuthorizationPolicy1")
			So(apiAuthorizationPoliciesCreated[0].AuthorizedNamespace, ShouldEqual, "/apomux/production")
			So(apiAuthorizationPoliciesCreated[0].Namespace, ShouldEqual, "/apomux")
			So(apiAuthorizationPoliciesCreated[0].Subject[0][0], ShouldEqual, "$namespace=/apomux/production")

			So(namespacesCreated[0].Name, ShouldEqual, "/apomux")
			So(namespacesCreated[0].Namespace, ShouldEqual, "/")

			So(namespacesCreated[1].Name, ShouldEqual, "/apomux/test")
			So(namespacesCreated[1].Namespace, ShouldEqual, "/apomux")

			So(namespacesCreated[2].Name, ShouldEqual, "/apomux/production")
			So(namespacesCreated[2].Namespace, ShouldEqual, "/apomux")

			So(namespacesCreated[3].Name, ShouldEqual, "/apomux/production/aporeto")
			So(namespacesCreated[3].Namespace, ShouldEqual, "/apomux/production")

			So(externalServicesCreated[0].Name, ShouldEqual, "externalService1")
			So(externalServicesCreated[0].Namespace, ShouldEqual, "/apomux")

			So(externalServicesCreated[1].Name, ShouldEqual, "externalService2")
			So(externalServicesCreated[1].Namespace, ShouldEqual, "/apomux")

			So(filePathscreated[1].Name, ShouldEqual, "filepath1")
			So(filePathscreated[1].Namespace, ShouldEqual, "/apomux")

			So(filePathscreated[0].Name, ShouldEqual, "filepath3")
			So(filePathscreated[0].Namespace, ShouldEqual, "/apomux/production")

			So(namespaceMappingPoliciesCreated[0].MappedNamespace, ShouldEqual, "/apomux/production/aporeto")
			So(namespaceMappingPoliciesCreated[0].Namespace, ShouldEqual, "/apomux/production")
			So(namespaceMappingPoliciesCreated[0].Subject[0][0], ShouldEqual, "$namespace=/apomux/production/aporeto")

			So(namespaceMappingPoliciesCreated[1].MappedNamespace, ShouldEqual, "/apomux/test")
			So(namespaceMappingPoliciesCreated[1].Namespace, ShouldEqual, "/apomux")
			So(namespaceMappingPoliciesCreated[1].Subject[0][0], ShouldEqual, "$namespace=/apomux/production/aporeto")
			So(namespaceMappingPoliciesCreated[1].Subject[1][0], ShouldEqual, "$namespace=/apomux/test")
		})

		Convey("Given importNamespaceContent is a success with the namespace /level", func() {
			namespacesCreated := squallmodels.NamespacesList{}
			externalServicesCreated := squallmodels.ExternalServicesList{}
			filePathscreated := squallmodels.FilePathsList{}
			namespaceMappingPoliciesCreated := squallmodels.NamespaceMappingPoliciesList{}
			apiAuthorizationPoliciesCreated := squallmodels.APIAuthorizationPoliciesList{}

			manipulator.MockRetrieveMany(t, func(context *manipulate.Context, dest elemental.ContentIdentifiable) error {

				if dest.ContentIdentity().Name == squallmodels.NamespaceIdentity.Name {
					namespaces := dest.(*squallmodels.NamespacesList)
					dest = namespaces
					_ = dest
					return nil
				}
				return nil
			})

			manipulator.MockCreate(t, func(ctx *manipulate.Context, objects ...elemental.Identifiable) error {

				if objects[0].Identity().Name == squallmodels.NamespaceIdentity.Name {
					namespace := objects[0].(*squallmodels.Namespace)
					if ctx.Namespace == "/" {
						namespace.Name = "/" + namespace.Name
					} else {
						namespace.Name = ctx.Namespace + "/" + namespace.Name
					}
					namespace.Namespace = ctx.Namespace
					namespacesCreated = append(namespacesCreated, namespace)
				}

				if objects[0].Identity().Name == squallmodels.ExternalServiceIdentity.Name {
					externalService := objects[0].(*squallmodels.ExternalService)
					externalService.Namespace = ctx.Namespace
					externalServicesCreated = append(externalServicesCreated, externalService)
				}

				if objects[0].Identity().Name == squallmodels.FilePathIdentity.Name {
					filePath := objects[0].(*squallmodels.FilePath)
					filePath.Namespace = ctx.Namespace
					filePathscreated = append(filePathscreated, filePath)
				}

				if objects[0].Identity().Name == squallmodels.NamespaceMappingPolicyIdentity.Name {
					namespaceMappingPolicy := objects[0].(*squallmodels.NamespaceMappingPolicy)
					namespaceMappingPolicy.Namespace = ctx.Namespace
					namespaceMappingPoliciesCreated = append(namespaceMappingPoliciesCreated, namespaceMappingPolicy)
				}

				if objects[0].Identity().Name == squallmodels.APIAuthorizationPolicyIdentity.Name {
					apiAuthorizationPolicy := objects[0].(*squallmodels.APIAuthorizationPolicy)
					apiAuthorizationPolicy.Namespace = ctx.Namespace
					apiAuthorizationPoliciesCreated = append(apiAuthorizationPoliciesCreated, apiAuthorizationPolicy)
				}

				return nil
			})

			err := Import(manipulator, "/level", content, false)

			So(err, ShouldBeNil)
			So(len(namespacesCreated), ShouldEqual, 4)
			So(len(externalServicesCreated), ShouldEqual, 2)
			So(len(filePathscreated), ShouldEqual, 2)
			So(len(apiAuthorizationPoliciesCreated), ShouldEqual, 1)
			So(len(namespaceMappingPoliciesCreated), ShouldEqual, 2)

			So(apiAuthorizationPoliciesCreated[0].Name, ShouldEqual, "apiAuthorizationPolicy1")
			So(apiAuthorizationPoliciesCreated[0].AuthorizedNamespace, ShouldEqual, "/level/apomux/production")
			So(apiAuthorizationPoliciesCreated[0].Namespace, ShouldEqual, "/level/apomux")
			So(apiAuthorizationPoliciesCreated[0].Subject[0][0], ShouldEqual, "$namespace=/level/apomux/production")

			So(namespacesCreated[0].Name, ShouldEqual, "/level/apomux")
			So(namespacesCreated[0].Namespace, ShouldEqual, "/level")

			So(namespacesCreated[1].Name, ShouldEqual, "/level/apomux/test")
			So(namespacesCreated[1].Namespace, ShouldEqual, "/level/apomux")

			So(namespacesCreated[2].Name, ShouldEqual, "/level/apomux/production")
			So(namespacesCreated[2].Namespace, ShouldEqual, "/level/apomux")

			So(namespacesCreated[3].Name, ShouldEqual, "/level/apomux/production/aporeto")
			So(namespacesCreated[3].Namespace, ShouldEqual, "/level/apomux/production")

			So(externalServicesCreated[0].Name, ShouldEqual, "externalService1")
			So(externalServicesCreated[0].Namespace, ShouldEqual, "/level/apomux")

			So(externalServicesCreated[1].Name, ShouldEqual, "externalService2")
			So(externalServicesCreated[1].Namespace, ShouldEqual, "/level/apomux")

			So(filePathscreated[1].Name, ShouldEqual, "filepath1")
			So(filePathscreated[1].Namespace, ShouldEqual, "/level/apomux")

			So(filePathscreated[0].Name, ShouldEqual, "filepath3")
			So(filePathscreated[0].Namespace, ShouldEqual, "/level/apomux/production")

			So(namespaceMappingPoliciesCreated[0].MappedNamespace, ShouldEqual, "/level/apomux/production/aporeto")
			So(namespaceMappingPoliciesCreated[0].Namespace, ShouldEqual, "/level/apomux/production")
			So(namespaceMappingPoliciesCreated[0].Subject[0][0], ShouldEqual, "$namespace=/level/apomux/production/aporeto")

			So(namespaceMappingPoliciesCreated[1].MappedNamespace, ShouldEqual, "/level/apomux/test")
			So(namespaceMappingPoliciesCreated[1].Namespace, ShouldEqual, "/level/apomux")
			So(namespaceMappingPoliciesCreated[1].Subject[0][0], ShouldEqual, "$namespace=/level/apomux/production/aporeto")
			So(namespaceMappingPoliciesCreated[1].Subject[1][0], ShouldEqual, "$namespace=/level/apomux/test")
		})

		Convey("Given importNamespaceContent is a success and previous namespace is deleted with namespace /", func() {

			namespacesDeleted := squallmodels.NamespacesList{}

			manipulator.MockRetrieveMany(t, func(context *manipulate.Context, dest elemental.ContentIdentifiable) error {

				if dest.ContentIdentity().Name == squallmodels.NamespaceIdentity.Name {
					namespaces := dest.(*squallmodels.NamespacesList)
					namespace := squallmodels.NewNamespace()
					namespace.Name = context.Filter.String()

					*namespaces = append(*namespaces, namespace)
					dest = namespaces
					_ = dest
					return nil
				}
				return nil
			})

			manipulator.MockCreate(t, func(ctx *manipulate.Context, objects ...elemental.Identifiable) error {
				return nil
			})

			manipulator.MockDelete(t, func(ctx *manipulate.Context, objects ...elemental.Identifiable) error {
				if objects[0].Identity().Name == squallmodels.NamespaceIdentity.Name {
					namespacesDeleted = append(namespacesDeleted, objects[0].(*squallmodels.Namespace))
				}

				return nil
			})

			err := Import(manipulator, "/", content, true)
			So(err, ShouldBeNil)
			So(len(namespacesDeleted), ShouldEqual, 4)
		})

		Convey("Given importNamespaceContent is a success and previous content is deleted with namespace /level", func() {

			networksAccessPolicy1 := squallmodels.NewNetworkAccessPolicy()
			networksAccessPolicy1.Name = "networksAccessPolicy1"

			fileAccessPolicy1 := squallmodels.NewFileAccessPolicy()
			fileAccessPolicy1.Name = "fileAccessPolicy1"

			var expectedDeletedFileAccess *squallmodels.FileAccessPolicy
			var expectedDeletedNetworkAccess *squallmodels.NetworkAccessPolicy

			manipulator.MockRetrieveMany(t, func(context *manipulate.Context, dest elemental.ContentIdentifiable) error {

				if dest.ContentIdentity().Name == squallmodels.NamespaceIdentity.Name {
					namespaces := dest.(*squallmodels.NamespacesList)
					namespace := squallmodels.NewNamespace()
					*namespaces = append(*namespaces, namespace)
					dest = namespaces
					_ = dest
					return nil
				}

				if dest.ContentIdentity().Name == squallmodels.NetworkAccessPolicyIdentity.Name {
					policies := dest.(*squallmodels.NetworkAccessPoliciesList)
					*policies = append(*policies, networksAccessPolicy1)
					dest = policies
					_ = dest
				}

				if dest.ContentIdentity().Name == squallmodels.FileAccessPolicyIdentity.Name {
					policies := dest.(*squallmodels.FileAccessPoliciesList)
					*policies = append(*policies, fileAccessPolicy1)
					dest = policies
					_ = dest
				}

				return nil
			})

			manipulator.MockCreate(t, func(ctx *manipulate.Context, objects ...elemental.Identifiable) error {
				return nil
			})

			manipulator.MockDelete(t, func(ctx *manipulate.Context, objects ...elemental.Identifiable) error {

				if objects[0].Identity().Name == squallmodels.FileAccessPolicyIdentity.Name {
					expectedDeletedFileAccess = objects[0].(*squallmodels.FileAccessPolicy)
				}
				if objects[0].Identity().Name == squallmodels.NetworkAccessPolicyIdentity.Name {
					expectedDeletedNetworkAccess = objects[0].(*squallmodels.NetworkAccessPolicy)
				}
				return nil
			})

			err := Import(manipulator, "/level", content, false)
			So(err, ShouldBeNil)
			So(expectedDeletedFileAccess.Name, ShouldEqual, "fileAccessPolicy1")
			So(expectedDeletedNetworkAccess.Name, ShouldEqual, "networksAccessPolicy1")
		})

		Convey("Given importNamespaceContent is a failure and previous namespace is deleted with namespace /", func() {

			namespacesDeleted := squallmodels.NamespacesList{}

			manipulator.MockRetrieveMany(t, func(context *manipulate.Context, dest elemental.ContentIdentifiable) error {

				if dest.ContentIdentity().Name == squallmodels.NamespaceIdentity.Name {
					namespaces := dest.(*squallmodels.NamespacesList)
					namespace := squallmodels.NewNamespace()
					namespace.Name = context.Filter.String()

					*namespaces = append(*namespaces, namespace)
					dest = namespaces
					_ = dest
					return nil
				}
				return nil
			})

			manipulator.MockDelete(t, func(ctx *manipulate.Context, objects ...elemental.Identifiable) error {
				if objects[0].Identity().Name == squallmodels.NamespaceIdentity.Name {
					namespacesDeleted = append(namespacesDeleted, objects[0].(*squallmodels.Namespace))
				}

				return elemental.NewError("Invalid Entity", "", "", 500)
			})

			err := Import(manipulator, "/", content, true)
			So(err, ShouldNotBeNil)
		})

		Convey("Given importNamespaceContent is a failure when retrieving content with namespace /level", func() {

			manipulator.MockRetrieveMany(t, func(context *manipulate.Context, dest elemental.ContentIdentifiable) error {

				if dest.ContentIdentity().Name == squallmodels.NamespaceIdentity.Name {
					namespaces := dest.(*squallmodels.NamespacesList)
					namespace := squallmodels.NewNamespace()
					*namespaces = append(*namespaces, namespace)
					dest = namespaces
					_ = dest
					return nil
				}

				if dest.ContentIdentity().Name == squallmodels.NetworkAccessPolicyIdentity.Name {
					policies := dest.(*squallmodels.NetworkAccessPoliciesList)
					dest = policies
					_ = dest
				}

				if dest.ContentIdentity().Name == squallmodels.FileAccessPolicyIdentity.Name {
					policies := dest.(*squallmodels.FileAccessPoliciesList)
					dest = policies
					_ = dest
					return elemental.NewError("Invalid Entity", "", "", 500)
				}

				return nil
			})

			err := Import(manipulator, "/level", content, false)
			So(err, ShouldNotBeNil)
		})

		Convey("Given importNamespaceContent is a failure when deleting the previous content with namespace /level", func() {

			networksAccessPolicy1 := squallmodels.NewNetworkAccessPolicy()
			networksAccessPolicy1.Name = "networksAccessPolicy1"

			fileAccessPolicy1 := squallmodels.NewFileAccessPolicy()
			fileAccessPolicy1.Name = "fileAccessPolicy1"

			var expectedDeletedFileAccess *squallmodels.FileAccessPolicy
			var expectedDeletedNetworkAccess *squallmodels.NetworkAccessPolicy

			manipulator.MockRetrieveMany(t, func(context *manipulate.Context, dest elemental.ContentIdentifiable) error {

				if dest.ContentIdentity().Name == squallmodels.NamespaceIdentity.Name {
					namespaces := dest.(*squallmodels.NamespacesList)
					namespace := squallmodels.NewNamespace()
					*namespaces = append(*namespaces, namespace)
					dest = namespaces
					_ = dest
					return nil
				}

				if dest.ContentIdentity().Name == squallmodels.NetworkAccessPolicyIdentity.Name {
					policies := dest.(*squallmodels.NetworkAccessPoliciesList)
					*policies = append(*policies, networksAccessPolicy1)
					dest = policies
					_ = dest
				}

				if dest.ContentIdentity().Name == squallmodels.FileAccessPolicyIdentity.Name {
					policies := dest.(*squallmodels.FileAccessPoliciesList)
					*policies = append(*policies, fileAccessPolicy1)
					dest = policies
					_ = dest
				}

				return nil
			})

			manipulator.MockCreate(t, func(ctx *manipulate.Context, objects ...elemental.Identifiable) error {
				return nil
			})

			manipulator.MockDelete(t, func(ctx *manipulate.Context, objects ...elemental.Identifiable) error {

				if objects[0].Identity().Name == squallmodels.FileAccessPolicyIdentity.Name {
					expectedDeletedFileAccess = objects[0].(*squallmodels.FileAccessPolicy)
				}
				if objects[0].Identity().Name == squallmodels.NetworkAccessPolicyIdentity.Name {
					expectedDeletedNetworkAccess = objects[0].(*squallmodels.NetworkAccessPolicy)
					return elemental.NewError("Invalid Entity", "", "", 500)
				}
				return nil
			})

			err := Import(manipulator, "/level", content, false)
			So(err, ShouldNotBeNil)
		})

		Convey("Given importNamespaceContent got an error when retrieving a namespace", func() {
			manipulator.MockRetrieveMany(t, func(context *manipulate.Context, dest elemental.ContentIdentifiable) error {
				return elemental.NewError("Invalid Entity", "", "", 500)
			})

			err := Import(manipulator, "/level", content, false)

			So(err, ShouldNotBeNil)
		})
	})
}
