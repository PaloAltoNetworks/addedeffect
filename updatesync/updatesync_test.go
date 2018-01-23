package updatesync

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/aporeto-inc/elemental"
	"github.com/aporeto-inc/gaia/squallmodels/v1/golang"
	"github.com/aporeto-inc/manipulate"
	"github.com/aporeto-inc/manipulate/maniptest"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAPI_UpdateSync(t *testing.T) {

	Convey("Given I have a manipulator an object and an update func", t, func() {

		var synced int
		m := maniptest.NewTestManipulator()
		o := squallmodels.NewProcessingUnit()
		o.Name = "name-original"
		o.Description = "desc-original"

		uf := func(obj elemental.Identifiable) {
			synced++
			obj.(*squallmodels.ProcessingUnit).Description = "synced!"
		}

		Convey("When I update my object there is no sync needed", func() {

			m.MockUpdate(t, func(ctx *manipulate.Context, objects ...elemental.Identifiable) error {
				return nil
			})

			err := UpdateSync(context.TODO(), m, nil, o, uf, 10)

			Convey("Then err should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then obj should have been updated", func() {
				So(o.Description, ShouldEqual, "synced!")
				So(o.Name, ShouldEqual, "name-original")
			})

			Convey("Then the updateFunc should have been called once", func() {
				So(synced, ShouldEqual, 1)
			})
		})

		Convey("When I update my object there is a sync needed", func() {

			m.MockUpdate(t, func(ctx *manipulate.Context, objects ...elemental.Identifiable) error {
				objects[0].(*squallmodels.ProcessingUnit).Name = fmt.Sprintf("sync%d", synced)

				if synced <= 3 {
					e := elemental.NewError("Read Only Error", "bloob", "subject", http.StatusUnprocessableEntity)
					e.Data = map[string]interface{}{"attribute": "updateTime"}
					return e
				}

				return nil
			})

			err := UpdateSync(context.TODO(), m, nil, o, uf, 10)

			Convey("Then err should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then obj should have been updated", func() {
				So(o.Name, ShouldEqual, "sync4")
				So(o.Description, ShouldEqual, "synced!")
			})

			Convey("Then the updateFunc should have been called four times", func() {
				So(synced, ShouldEqual, 4)
			})
		})

		Convey("When I update my object but there is an error right away", func() {

			m.MockUpdate(t, func(ctx *manipulate.Context, objects ...elemental.Identifiable) error {
				return elemental.NewError("Not Read Only Error", "bloob", "subject", http.StatusInternalServerError)
			})

			err := UpdateSync(context.TODO(), m, nil, o, uf, 10)

			Convey("Then err should not be nil", func() {
				So(err, ShouldNotBeNil)
			})

			// TODO: that would be ideal, but this involves deep copy,
			// and I'm not confident about it.
			// Convey("Then obj should not have been updated", func() {
			// 	So(o.Name, ShouldEqual, "name-original")
			// 	So(o.Description, ShouldEqual, "desc-original")
			// })

			Convey("Then the updateFunc should have been called once", func() {
				So(synced, ShouldEqual, 1)
			})
		})

		Convey("When I update my object there is a sync needed but if fails more than the number of tries", func() {

			m.MockUpdate(t, func(ctx *manipulate.Context, objects ...elemental.Identifiable) error {
				e := elemental.NewError("Read Only Error", "bloob", "subject", http.StatusUnprocessableEntity)
				e.Data = map[string]interface{}{"attribute": "updateTime"}
				return e
			})

			err := UpdateSync(context.TODO(), m, nil, o, uf, 2)

			Convey("Then err should not be nil", func() {
				So(err, ShouldNotBeNil)
			})

			Convey("Then the updateFunc should have been called four times", func() {
				So(synced, ShouldEqual, 2)
			})
		})

		Convey("When I update my object there is a sync needed but the retrieve fails", func() {

			m.MockUpdate(t, func(ctx *manipulate.Context, objects ...elemental.Identifiable) error {
				objects[0].(*squallmodels.ProcessingUnit).Name = fmt.Sprintf("sync%d", synced)

				if synced <= 3 {
					e := elemental.NewError("Read Only Error", "bloob", "subject", http.StatusUnprocessableEntity)
					e.Data = map[string]interface{}{"attribute": "updateTime"}
					return e
				}

				return nil
			})

			m.MockRetrieve(t, func(ctx *manipulate.Context, objects ...elemental.Identifiable) error {
				return fmt.Errorf("boom")
			})

			err := UpdateSync(context.TODO(), m, nil, o, uf, 10)

			Convey("Then err should not be nil", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "boom")
			})

			Convey("Then the updateFunc should have been called four times", func() {
				So(synced, ShouldEqual, 1)
			})
		})
	})
}
