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
       test init --option-e=<value> --option-f [--option-g=<value>] [--option-h]
`

func TestCommon_ProcessCreate(t *testing.T) {

	Convey("Given I have no os.Args and no environment variables", t, func() {

		os.Args = []string{}

		Convey("When I call Parse on the usage", func() {

			Parse("TEST", usage)

			Convey("Then os.Args should be empty", func() {
				So(len(os.Args), ShouldEqual, 0)
			})
		})
	})

	Convey("Given I have no os.Args and some environment variables", t, func() {

		os.Args = []string{}

		Convey("When I call Parse on the usage", func() {

			os.Setenv("TEST_OPTION_A", "value1")
			os.Setenv("TEST_OPTION_B", "1")
			os.Setenv("TEST_OPTION_C", "value2")
			os.Setenv("TEST_OPTION_D", "2")

			os.Setenv("TEST_OPTION_E", "value1")
			os.Setenv("TEST_OPTION_F", "1")
			os.Setenv("TEST_OPTION_G", "value2")
			os.Setenv("TEST_OPTION_H", "2")

			Parse("TEST", usage)

			Convey("Then os.Args should not be empty", func() {
				So(len(os.Args), ShouldEqual, 8)
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
		})
	})
}
