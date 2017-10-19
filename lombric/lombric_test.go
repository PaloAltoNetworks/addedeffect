package lombric

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	. "github.com/smartystreets/goconvey/convey"
)

const usage = `      --a-bool                                This is a boolean [required] (default true)
      --a-bool-nodef                          This is a no def boolean
      --a-duration duration                   This is a duration [required] (default 10s)
      --a-duration-nodef duration             This is a no def duration
      --a-integer int                         This is a number [required] (default 42)
      --a-integer-nodef int                   This is a no def number
      --a-string string                       This is a string [required] (default "hello")
      --a-string-nodef string                 This is a no def string
      --a-string-slice stringSlice            This is a string slice [required] (default [a,b,c])
      --a-string-slice-from-var stringSlice   This is a no def string slice populated from var
      --a-string-slice-nodef stringSlice      This is a no def string slice
`

type testConf struct {
	AString      string        `mapstructure:"a-string"        desc:"This is a string"       required:"true" default:"hello"`
	ABool        bool          `mapstructure:"a-bool"          desc:"This is a boolean"      required:"true" default:"true"`
	ADuration    time.Duration `mapstructure:"a-duration"      desc:"This is a duration"     required:"true" default:"10s"`
	AInteger     int           `mapstructure:"a-integer"       desc:"This is a number"       required:"true" default:"42"`
	AStringSlice []string      `mapstructure:"a-string-slice"  desc:"This is a string slice" required:"true" default:"a,b,c"`

	AStringNoDef      string        `mapstructure:"a-string-nodef"        desc:"This is a no def string"`
	ABoolNoDef        bool          `mapstructure:"a-bool-nodef"          desc:"This is a no def boolean"`
	ADurationNoDef    time.Duration `mapstructure:"a-duration-nodef"      desc:"This is a no def duration"`
	AIntegerNoDef     int           `mapstructure:"a-integer-nodef"       desc:"This is a no def number"`
	AStringSliceNoDef []string      `mapstructure:"a-string-slice-nodef"  desc:"This is a no def string slice"`

	AnotherStringSliceNoDef []string `mapstructure:"a-string-slice-from-var"  desc:"This is a no def string slice populated from var"`
}

// Prefix return the configuration prefix.
func (c *testConf) Prefix() string { return "lombric" }

func TestLombric_Initialize(t *testing.T) {

	Convey("Given have a conf", t, func() {

		conf := &testConf{}
		Initialize(conf)

		Convey("Then the flags should be correctly set", func() {
			os.Setenv("LOMBRIC_A_STRING_SLICE_FROM_VAR", "x y z") // nolint: errcheck

			So(viper.GetString("a-string"), ShouldEqual, "hello")
			So(viper.GetBool("a-bool"), ShouldEqual, true)
			So(viper.GetDuration("a-duration"), ShouldEqual, 10*time.Second)
			So(viper.GetInt("a-integer"), ShouldEqual, 42)
			So(viper.GetStringSlice("a-string-slice"), ShouldResemble, []string{"a", "b", "c"})

			So(viper.GetString("a-string-nodef"), ShouldEqual, "")
			So(viper.GetBool("a-bool-nodef"), ShouldEqual, false)
			So(viper.GetDuration("a-duration-nodef"), ShouldEqual, 0)
			So(viper.GetInt("a-integer-nodef"), ShouldEqual, 0)
			So(viper.GetStringSlice("a-string-slice-nodef"), ShouldResemble, []string{})

			So(viper.GetStringSlice("a-string-slice-from-var"), ShouldResemble, []string{"x", "y", "z"})

			So(strings.Replace(pflag.CommandLine.FlagUsages(), " ", "", -1), ShouldEqual, strings.Replace(usage, " ", "", -1))
		})

	})
}
