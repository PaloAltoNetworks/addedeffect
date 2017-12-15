package tag

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTag_Split(t *testing.T) {

	Convey("Given I have a tag a=b", t, func() {
		t := "a=b"

		Convey("When I call Split", func() {
			k, v, e := Split(t)

			Convey("Then e should be nil", func() {
				So(e, ShouldBeNil)
			})

			Convey("Then k should equal a", func() {
				So(k, ShouldEqual, "a")
			})

			Convey("Then v should equal b", func() {
				So(v, ShouldEqual, "b")
			})

		})
	})

	Convey("Given I have a tag a=b c", t, func() {
		t := "a=b c"

		Convey("When I call Split", func() {
			k, v, e := Split(t)

			Convey("Then e should be nil", func() {
				So(e, ShouldBeNil)
			})

			Convey("Then k should equal a", func() {
				So(k, ShouldEqual, "a")
			})

			Convey("Then v should equal b c", func() {
				So(v, ShouldEqual, "b c")
			})

		})
	})

	Convey("Given I have a tag a=b c=ddd", t, func() {
		t := "a=b c=ddd"

		Convey("When I call Split", func() {
			k, v, e := Split(t)

			Convey("Then e should be nil", func() {
				So(e, ShouldBeNil)
			})

			Convey("Then k should equal a", func() {
				So(k, ShouldEqual, "a")
			})

			Convey("Then v should equal b c=ddd", func() {
				So(v, ShouldEqual, "b c=ddd")
			})

		})
	})

	Convey("Given I have a tag a", t, func() {
		t := "a"

		Convey("When I call Split", func() {
			_, _, e := Split(t)

			Convey("Then e should not be nil", func() {
				So(e.Error(), ShouldEqual, "Invalid tag: invalid length 'a'")
			})

		})
	})

	Convey("Given I have a tag a", t, func() {
		t := "a="

		Convey("When I call Split", func() {
			_, _, e := Split(t)

			Convey("Then e should not be nil", func() {
				So(e.Error(), ShouldEqual, "Invalid tag: invalid length 'a='")
			})

		})
	})

	Convey("Given I have a tag a", t, func() {
		t := "abc"

		Convey("When I call Split", func() {
			_, _, e := Split(t)

			Convey("Then e should not be nil", func() {
				So(e.Error(), ShouldEqual, "Invalid tag: missing equal symbol 'abc'")
			})

		})
	})

	Convey("Given I have a tag a", t, func() {
		t := "abc="

		Convey("When I call Split", func() {
			_, _, e := Split(t)

			Convey("Then e should not be nil", func() {
				So(e.Error(), ShouldEqual, "Invalid tag: missing value 'abc='")
			})

		})
	})

	Convey("Given I have a tag a", t, func() {
		t := "=abc"

		Convey("When I call Split", func() {
			_, _, e := Split(t)

			Convey("Then e should not be nil", func() {
				So(e.Error(), ShouldEqual, "Invalid tag: missing key '=abc'")
			})

		})
	})
}

func TestTag_SplitPtr(t *testing.T) {

	Convey("Given I have a tag a=b", t, func() {

		var k, v string
		t := "a=b"

		Convey("When I call Split", func() {

			e := SplitPtr(t, &k, &v)

			Convey("Then e should be nil", func() {
				So(e, ShouldBeNil)
			})

			Convey("Then k should equal a", func() {
				So(k, ShouldEqual, "a")
			})

			Convey("Then v should equal b", func() {
				So(v, ShouldEqual, "b")
			})

		})
	})

	Convey("Given I have a tag a=b c", t, func() {

		var k, v string
		t := "a=b c"

		Convey("When I call Split", func() {

			e := SplitPtr(t, &k, &v)

			Convey("Then e should be nil", func() {
				So(e, ShouldBeNil)
			})

			Convey("Then k should equal a", func() {
				So(k, ShouldEqual, "a")
			})

			Convey("Then v should equal b c", func() {
				So(v, ShouldEqual, "b c")
			})

		})
	})

	Convey("Given I have a tag a=b c=ddd", t, func() {

		var k, v string
		t := "a=b c=ddd"

		Convey("When I call Split", func() {

			e := SplitPtr(t, &k, &v)

			Convey("Then e should be nil", func() {
				So(e, ShouldBeNil)
			})

			Convey("Then k should equal a", func() {
				So(k, ShouldEqual, "a")
			})

			Convey("Then v should equal b c=ddd", func() {
				So(v, ShouldEqual, "b c=ddd")
			})

		})
	})

	Convey("Given I have a tag a", t, func() {

		var k, v string
		t := "a"

		Convey("When I call Split", func() {

			e := SplitPtr(t, &k, &v)

			Convey("Then e should not be nil", func() {
				So(e.Error(), ShouldEqual, "Invalid tag: invalid length 'a'")
			})

		})
	})

	Convey("Given I have a tag a", t, func() {

		var k, v string
		t := "a="

		Convey("When I call Split", func() {

			e := SplitPtr(t, &k, &v)

			Convey("Then e should not be nil", func() {
				So(e.Error(), ShouldEqual, "Invalid tag: invalid length 'a='")
			})

		})
	})

	Convey("Given I have a tag a", t, func() {

		var k, v string
		t := "abc"

		Convey("When I call Split", func() {

			e := SplitPtr(t, &k, &v)

			Convey("Then e should not be nil", func() {
				So(e.Error(), ShouldEqual, "Invalid tag: missing equal symbol 'abc'")
			})

		})
	})

	Convey("Given I have a tag a", t, func() {

		var k, v string
		t := "abc="

		Convey("When I call Split", func() {

			e := SplitPtr(t, &k, &v)

			Convey("Then e should not be nil", func() {
				So(e.Error(), ShouldEqual, "Invalid tag: missing value 'abc='")
			})

		})
	})

	Convey("Given I have a tag a", t, func() {

		var k, v string
		t := "=abc"

		Convey("When I call Split", func() {

			e := SplitPtr(t, &k, &v)

			Convey("Then e should not be nil", func() {
				So(e.Error(), ShouldEqual, "Invalid tag: missing key '=abc'")
			})

		})
	})
}
