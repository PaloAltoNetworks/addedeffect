package envopt

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

const usage = `
Usage: test -h | --help
       test -v | --version
       test init
        --option-a=<value>
        --option-b
        [--option-c=<value>]
        [--option-d]
        [--option-i]
        [--option-j]
       test init --option-e=<value> --option-f [--option-g=<value>] [--option-h]
`

func TestCommon_ProcessCreate(t *testing.T) {

	Convey("Given I have no os.Args and no environment variables", t, func() {

		os.Args = []string{}

		Convey("When I call Parse on the usage", func() {

			_ = Parse("TEST", usage)

			Convey("Then os.Args should be empty", func() {
				So(len(os.Args), ShouldEqual, 0)
			})
		})
	})

	Convey("Given I have no os.Args and some environment variables", t, func() {

		os.Args = []string{}

		Convey("When I call Parse on the usage", func() {

			_ = os.Setenv("TEST_OPTION_A", "value1")
			_ = os.Setenv("TEST_OPTION_B", "1")
			_ = os.Setenv("TEST_OPTION_C", "value2")
			_ = os.Setenv("TEST_OPTION_D", "2")

			_ = os.Setenv("TEST_OPTION_E", "value1")
			_ = os.Setenv("TEST_OPTION_F", "1")
			_ = os.Setenv("TEST_OPTION_G", "value2")
			_ = os.Setenv("TEST_OPTION_H", "2")

			_ = os.Setenv("TEST_OPTION_I", "false")
			_ = os.Setenv("TEST_OPTION_J", "true")

			_ = Parse("TEST", usage)

			Convey("Then os.Args should not be empty", func() {
				So(len(os.Args), ShouldEqual, 9)
			})

			Convey("Then os.Args should have the correct flag for option-a", func() {
				So(os.Args, ShouldContain, "--option-a=value1")
			})

			Convey("Then os.Args should have the correct flag for option-b", func() {
				So(os.Args, ShouldContain, "--option-b")
			})

			Convey("Then os.Args should have the correct flag for option-c", func() {
				So(os.Args, ShouldContain, "--option-c=value2")
			})

			Convey("Then os.Args should have the correct flag for option-d", func() {
				So(os.Args, ShouldContain, "--option-d")
			})

			Convey("Then os.Args should have the correct flag for option-e", func() {
				So(os.Args, ShouldContain, "--option-e=value1")
			})

			Convey("Then os.Args should have the correct flag for option-f", func() {
				So(os.Args, ShouldContain, "--option-f")
			})

			Convey("Then os.Args should have the correct flag for option-g", func() {
				So(os.Args, ShouldContain, "--option-g=value2")
			})

			Convey("Then os.Args should have the correct flag for option-h", func() {
				So(os.Args, ShouldContain, "--option-h")
			})

			Convey("Then os.Args should have the correct flag for option-i", func() {
				So(os.Args, ShouldNotContain, "--option-i")
			})

			Convey("Then os.Args should have the correct flag for option-j", func() {
				So(os.Args, ShouldContain, "--option-j")
			})
		})
	})

	Convey("Given I have some os.Args that already contains an option set in environment", t, func() {

		os.Args = []string{"--option-x=a", "--option-y"}
		_ = os.Setenv("TEST_OPTION_X", "not-a")
		_ = os.Setenv("TEST_OPTION_Y", "1")

		Convey("When I use Parse", func() {

			_ = Parse("TEST", usage)

			Convey("Then the original options should remain unchanged", func() {
				So(os.Args, ShouldContain, "--option-x=a")
				So(os.Args, ShouldContain, "--option-y")
			})
		})
	})
}
