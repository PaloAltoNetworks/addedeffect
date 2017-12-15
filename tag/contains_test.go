package tag

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTag_Contains(t *testing.T) {

	Convey("Given I have a list of tag values", t, func() {

		lst := []string{"a=b", "c=d"}

		Convey("When I check if the list contains 'a=b'", func() {

			ok := Contains(lst, "a=b")

			Convey("Then it should be found", func() {
				So(ok, ShouldBeTrue)
			})
		})

		Convey("When I check if the list contains 'nope=b'", func() {

			ok := Contains(lst, "nope=b")

			Convey("Then it should be found", func() {
				So(ok, ShouldBeFalse)
			})
		})
	})
}
