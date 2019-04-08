package appcreds

import (
	"context"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"go.aporeto.io/elemental"
	"go.aporeto.io/gaia"
	"go.aporeto.io/manipulate"
	"go.aporeto.io/manipulate/maniptest"
	"go.aporeto.io/tg/tglib"
)

func TestAppCred_New(t *testing.T) {

	Convey("Given I have a manipulator", t, func() {

		m := maniptest.NewTestManipulator()

		var expectedCSR string
		m.MockCreate(t, func(ctx manipulate.Context, object elemental.Identifiable) error {

			if ctx.Namespace() != "/ns" {
				panic("expected ns to be /ns")
			}

			ac := object.(*gaia.AppCredential)
			ac.ID = "ID"
			ac.Namespace = "/ns"
			ac.Credentials = gaia.NewCredential()
			ac.Credentials.APIURL = "https://labas"
			ac.Credentials.Name = ac.Name
			ac.Credentials.Namespace = ac.Namespace

			expectedCSR = ac.CSR

			return nil
		})

		Convey("When I call New", func() {

			c, err := New(context.Background(), m, "/ns", "name", []string{"@auth:role=role1"}, nil)

			Convey("Then err should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then the cred should be correct", func() {
				So(c.Name, ShouldEqual, "name")
				So(c.ID, ShouldEqual, "ID")
				So(c.Namespace, ShouldEqual, "/ns")
				So(c.Credentials.CertificateKey, ShouldNotBeEmpty)
			})

			Convey("When I verify the csr", func() {

				csrs, err := tglib.LoadCSRs([]byte(expectedCSR))

				Convey("Then err should be nil", func() {
					So(err, ShouldBeNil)
				})

				Convey("Then csr should be correct", func() {
					So(len(csrs), ShouldEqual, 1)
				})
			})
		})
	})

	Convey("Given I have a manipulator that fails at creation", t, func() {

		m := maniptest.NewTestManipulator()

		m.MockCreate(t, func(ctx manipulate.Context, object elemental.Identifiable) error {
			return fmt.Errorf("boom")
		})

		Convey("When I call New", func() {

			c, err := New(context.Background(), m, "/ns", "name", []string{"@auth:role=role1"}, nil)

			Convey("Then err should not be nil", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "boom")
			})

			Convey("Then the cred should be nilt", func() {
				So(c, ShouldBeNil)
			})
		})
	})
}

func TestAppCred_NewWithAppCredential(t *testing.T) {

	Convey("Given I have a manipulator", t, func() {

		m := maniptest.NewTestManipulator()

		var expectedCSR string
		m.MockCreate(t, func(ctx manipulate.Context, object elemental.Identifiable) error {

			if ctx.Namespace() != "/ns" {
				panic("expected ns to be /ns")
			}

			ac := object.(*gaia.AppCredential)
			ac.ID = "ID"
			ac.Namespace = "/ns"
			ac.Credentials = gaia.NewCredential()
			ac.Credentials.APIURL = "https://labas"
			ac.Credentials.Name = ac.Name
			ac.Credentials.Namespace = ac.Namespace

			expectedCSR = ac.CSR

			return nil
		})

		Convey("When I call NewWithAppCredential", func() {

			template := gaia.NewAppCredential()
			template.Name = "name"
			template.Description = "description"
			template.Protected = true
			template.Metadata = []string{"random=tag"}
			template.Roles = []string{"role=test"}
			template.Namespace = "/ns"

			c, err := NewWithAppCredential(context.Background(), m, template)

			Convey("Then credential should have template information", func() {
				So(c.Name, ShouldEqual, template.Name)
				So(c.Description, ShouldEqual, template.Description)
				So(c.Protected, ShouldEqual, template.Protected)
				So(c.Metadata, ShouldResemble, template.Metadata)
				So(c.Roles, ShouldResemble, template.Roles)
				So(c.Namespace, ShouldEqual, template.Namespace)
			})

			Convey("Then err should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then the cred should be correct", func() {
				So(c.Name, ShouldEqual, "name")
				So(c.ID, ShouldEqual, "ID")
				So(c.Namespace, ShouldEqual, "/ns")
				So(c.Credentials.CertificateKey, ShouldNotBeEmpty)
			})

			Convey("When I verify the csr", func() {

				csrs, err := tglib.LoadCSRs([]byte(expectedCSR))

				Convey("Then err should be nil", func() {
					So(err, ShouldBeNil)
				})

				Convey("Then csr should be correct", func() {
					So(len(csrs), ShouldEqual, 1)
				})
			})
		})

	})

	Convey("Given I have a manipulator that fails at creation", t, func() {

		m := maniptest.NewTestManipulator()

		m.MockCreate(t, func(ctx manipulate.Context, object elemental.Identifiable) error {
			return fmt.Errorf("boom")
		})

		Convey("When I call New", func() {

			template := gaia.NewAppCredential()
			template.Name = "name"
			template.Description = "description"
			template.Protected = true
			template.Metadata = []string{"random=tag"}
			template.Roles = []string{"role=test"}
			template.Namespace = "/ns"

			c, err := NewWithAppCredential(context.Background(), m, template)

			Convey("Then err should not be nil", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "boom")
			})

			Convey("Then the cred should be nilt", func() {
				So(c, ShouldBeNil)
			})
		})
	})
}

func TestAppCred_Renew(t *testing.T) {

	Convey("Given I have a manipulator", t, func() {

		m := maniptest.NewTestManipulator()

		ac := gaia.NewAppCredential()
		ac.ID = "ID"
		ac.Name = "name"
		ac.Namespace = "/ns"
		ac.Roles = []string{"@auth:role=role1"}

		var expectedCSR string
		m.MockUpdate(t, func(ctx manipulate.Context, object elemental.Identifiable) error {

			if ctx.Namespace() != "/ns" {
				panic("expected ns to be /ns")
			}

			ac := object.(*gaia.AppCredential)
			ac.Credentials = gaia.NewCredential()
			ac.Credentials.APIURL = "https://labas"
			ac.Credentials.Name = ac.Name
			ac.Credentials.Namespace = ac.Namespace

			expectedCSR = ac.CSR

			return nil
		})

		Convey("When I call Renew", func() {

			c, err := Renew(context.Background(), m, ac)

			Convey("Then err should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then the cred should be correct", func() {
				So(c.Name, ShouldEqual, "name")
				So(c.ID, ShouldEqual, "ID")
				So(c.Namespace, ShouldEqual, "/ns")
				So(c.Credentials.CertificateKey, ShouldNotBeEmpty)
			})

			Convey("When I verify the csr", func() {

				csrs, err := tglib.LoadCSRs([]byte(expectedCSR))

				Convey("Then err should be nil", func() {
					So(err, ShouldBeNil)
				})

				Convey("Then csr should be correct", func() {
					So(len(csrs), ShouldEqual, 1)
				})
			})
		})
	})

	Convey("Given I have a manipulator that fails at update", t, func() {

		m := maniptest.NewTestManipulator()

		ac := gaia.NewAppCredential()
		ac.ID = "ID"
		ac.Name = "name"
		ac.Namespace = "/ns"
		ac.Roles = []string{"@auth:role=role1"}

		m.MockUpdate(t, func(ctx manipulate.Context, object elemental.Identifiable) error {
			return fmt.Errorf("paf")
		})

		Convey("When I call New", func() {

			c, err := Renew(context.Background(), m, ac)

			Convey("Then err should not be nil", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "paf")
			})

			Convey("Then the cred should be nilt", func() {
				So(c, ShouldBeNil)
			})
		})
	})
}
