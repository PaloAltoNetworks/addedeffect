package tag

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTagging_Diff(t *testing.T) {

	Convey("Given I have a two set of tags", t, func() {

		s1 := []string{"a=a", "b=b", "e=e"}
		s2 := []string{"a=a", "b=b", "c=c", "d=d"}

		Convey("When I use Diff", func() {

			added, removed := Diff(s1, s2)

			Convey("Then the added list should be correct", func() {
				So(added, ShouldResemble, []string{"e=e"})
			})

			Convey("Then the removed list should be correct", func() {
				So(removed, ShouldResemble, []string{"c=c", "d=d"})
			})
		})
	})
}
