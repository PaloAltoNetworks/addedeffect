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
