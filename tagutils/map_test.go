package tagutils

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTag_ToMap(t *testing.T) {

	Convey("Given given I have a tag list", t, func() {

		lst := []string{"@image=hello/world", "@private=ok", "private=nok"}

		Convey("When I convert it into a map", func() {

			m := ToMap(lst)

			Convey("Then the result should be correct", func() {
				So(len(m), ShouldEqual, 3)
				So(m["@image"], ShouldEqual, "hello/world")
				So(m["@private"], ShouldEqual, "ok")
				So(m["private"], ShouldEqual, "nok")
			})
		})
	})

	Convey("Given given I an invalid tag list", t, func() {

		lst := []string{"@imagehello/world", "@private=ok", "private=nok"}

		Convey("When I convert it into a map", func() {

			Convey("Then it should panic", func() {
				So(func() { ToMap(lst) }, ShouldPanic)
			})
		})
	})
}
