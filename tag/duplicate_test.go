package tag

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTag_RemoveDuplicates(t *testing.T) {

	Convey("Given I have a list of tag values", t, func() {

		lst := []string{"a=b", "c=d", "c=d", "a=b"}

		Convey("When I use RemoveDuplicateTagStrings", func() {

			r := RemoveDuplicates(lst)

			Convey("Then the tags values list should not have any duplicate", func() {
				So(r, ShouldResemble, []string{"a=b", "c=d"})
			})
		})
	})
}
