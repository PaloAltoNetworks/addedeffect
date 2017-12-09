package lombric

import (
	"os"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"
)

// const usage = `      --a-bool                                This is a boolean [required] (default true)
//       --a-bool-nodef                          This is a no def boolean
//       --a-duration duration                   This is a duration [required] (default 10s)
//       --a-duration-nodef duration             This is a no def duration
//       --a-integer int                         This is a number [required] (default 42)
//       --a-integer-nodef int                   This is a no def number
//       --a-string string                       This is a string [required] (default "hello")
//       --a-string-nodef string                 This is a no def string
//       --a-string-slice stringSlice            This is a string slice [required] (default [a,b,c])
//       --a-string-slice-from-var stringSlice   This is a no def string slice populated from var
//       --a-string-slice-nodef stringSlice      This is a no def string slice
// `

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
	ASecret                 string   `mapstructure:"a-secret-from-var"        desc:"This is a secret"       secret:"true"`

	embedTestConf `mapstructure:",squash" override:"embeded-string-a=outter1,embeded-ignored-string=-"`
}

type embedTestConf struct {
	EmbededStringA        string `mapstructure:"embeded-string-a"        desc:"This is a string"       required:"true" default:"inner1"`
	EmbededStringB        string `mapstructure:"embeded-string-b"        desc:"This is a string"       required:"true" default:"inner2"`
	EmbededIgnoredStringB string `mapstructure:"embeded-ignored-string"  desc:"This is a string"       required:"true" default:"inner3"`
}

// Prefix return the configuration prefix.
func (c *testConf) Prefix() string { return "lombric" }

func TestLombric_Initialize(t *testing.T) {

	Convey("Given have a conf", t, func() {

		conf := &testConf{}
		os.Setenv("LOMBRIC_A_STRING_SLICE_FROM_VAR", "x y z") // nolint: errcheck
		os.Setenv("LOMBRIC_A_SECRET_FROM_VAR", "secret")      // nolint: errcheck

		Initialize(conf)

		Convey("Then the flags should be correctly set", func() {

			So(conf.AString, ShouldEqual, "hello")
			So(conf.ABool, ShouldEqual, true)
			So(conf.ADuration, ShouldEqual, 10*time.Second)
			So(conf.AInteger, ShouldEqual, 42)
			So(conf.AStringSlice, ShouldResemble, []string{"a", "b", "c"})

			So(conf.AStringNoDef, ShouldEqual, "")
			So(conf.ABoolNoDef, ShouldEqual, false)
			So(conf.ADurationNoDef, ShouldEqual, 0)
			So(conf.AIntegerNoDef, ShouldEqual, 0)
			So(conf.AStringSliceNoDef, ShouldResemble, []string{})

			So(conf.EmbededStringA, ShouldEqual, "outter1")
			So(conf.EmbededStringB, ShouldEqual, "inner2")
			So(conf.EmbededIgnoredStringB, ShouldEqual, "")
			So(viper.AllKeys(), ShouldNotContain, "embeded-ignored-string")

			So(conf.AStringSliceNoDef, ShouldResemble, []string{})

			So(conf.AnotherStringSliceNoDef, ShouldResemble, []string{"x", "y", "z"})
			So(conf.ASecret, ShouldEqual, "secret")
			So(os.Getenv("LOMBRIC_A_SECRET_FROM_VAR"), ShouldEqual, "")

			// This test is disabled because here we have stringSlice and on concourse we get strings...
			// So(strings.Replace(pflag.CommandLine.FlagUsages(), " ", "", -1), ShouldEqual, strings.Replace(usage, " ", "", -1))
		})

	})
}
