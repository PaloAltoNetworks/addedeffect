package lombric

import (
	"net"
	"os"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/viper"
)

type testConf struct {
	ABool                   bool          `mapstructure:"a-bool"                    desc:"This is a boolean"            required:"true" default:"true"`
	ABoolNoDef              bool          `mapstructure:"a-bool-nodef"              desc:"This is a no def boolean"     `
	ABoolSlice              []bool        `mapstructure:"a-bool-slice"              desc:"This is a bool slice"         required:"true" default:"true,false,true"`
	ADuration               time.Duration `mapstructure:"a-duration"                desc:"This is a duration"           required:"true" default:"10s"`
	ADurationNoDef          time.Duration `mapstructure:"a-duration-nodef"          desc:"This is a no def duration"    `
	AInteger                int           `mapstructure:"a-integer"                 desc:"This is a number"             required:"true" default:"42"`
	AIntegerNoDef           int           `mapstructure:"a-integer-nodef"           desc:"This is a no def number"      `
	AIntSlice               []int         `mapstructure:"a-int-slice"               desc:"This is a int slice"          required:"true" default:"1,2,3"`
	AnEnum                  string        `mapstructure:"a-enum"                    desc:"This is an enum"              allowed:"a,b,c" default:"a"`
	AnIPSlice               []net.IP      `mapstructure:"a-ip-slice"                desc:"This is an ip slice"          default:"127.0.0.1,192.168.100.1"`
	AnotherStringSliceNoDef []string      `mapstructure:"a-string-slice-from-var"   desc:"This is a no def string"      `
	ASecret                 string        `mapstructure:"a-secret-from-var"         desc:"This is a secret"             secret:"true"`
	AString                 string        `mapstructure:"a-string"                  desc:"This is a string"             required:"true" default:"hello"`
	AStringNoDef            string        `mapstructure:"a-string-nodef"            desc:"This is a no def string"      `
	AStringSlice            []string      `mapstructure:"a-string-slice"            desc:"This is a string slice"       required:"true" default:"a,b,c"`
	AStringSliceNoDef       []string      `mapstructure:"a-string-slice-nodef"      desc:"This is a no def string slice"`

	embedTestConf `mapstructure:",squash" override:"embedded-string-a=outter1,embedded-ignored-string=-"`
}

type embedTestConf struct {
	EmbeddedStringA        string `mapstructure:"embedded-string-a"        desc:"This is a string"       required:"true" default:"inner1"`
	EmbeddedStringB        string `mapstructure:"embedded-string-b"        desc:"This is a string"       required:"true" default:"inner2"`
	EmbeddedIgnoredStringB string `mapstructure:"embedded-ignored-string"  desc:"This is a string"       required:"true" default:"inner3"`
}

// Prefix return the configuration prefix.
func (c *testConf) Prefix() string { return "lombric" }
func (c *testConf) PrintVersion()  {}

func TestLombric_Initialize(t *testing.T) {

	Convey("Given have a conf", t, func() {

		conf := &testConf{}
		os.Setenv("LOMBRIC_A_STRING_SLICE_FROM_VAR", "x y z") // nolint: errcheck
		os.Setenv("LOMBRIC_A_SECRET_FROM_VAR", "secret")      // nolint: errcheck

		Initialize(conf)

		Convey("Then the flags should be correctly set", func() {

			So(conf.ABool, ShouldEqual, true)
			So(conf.ABoolNoDef, ShouldEqual, false)
			So(conf.ABoolSlice, ShouldResemble, []bool{true, false, true})
			So(conf.ADuration, ShouldEqual, 10*time.Second)
			So(conf.ADurationNoDef, ShouldEqual, 0)
			So(conf.AInteger, ShouldEqual, 42)
			So(conf.AIntegerNoDef, ShouldEqual, 0)
			So(conf.AIntSlice, ShouldResemble, []int{1, 2, 3})
			So(conf.AnIPSlice, ShouldResemble, []net.IP{net.IPv4(127, 0, 0, 1), net.IPv4(192, 168, 100, 1)})
			So(conf.AnotherStringSliceNoDef, ShouldResemble, []string{"x", "y", "z"})
			So(conf.ASecret, ShouldEqual, "secret")
			So(conf.AString, ShouldEqual, "hello")
			So(conf.AStringNoDef, ShouldEqual, "")
			So(conf.AStringSlice, ShouldResemble, []string{"a", "b", "c"})
			So(conf.AStringSliceNoDef, ShouldResemble, []string(nil))
			So(conf.EmbeddedIgnoredStringB, ShouldEqual, "")
			So(conf.EmbeddedStringA, ShouldEqual, "outter1")
			So(conf.EmbeddedStringB, ShouldEqual, "inner2")
			So(os.Getenv("LOMBRIC_A_SECRET_FROM_VAR"), ShouldEqual, "")
			So(viper.AllKeys(), ShouldNotContain, "embedded-ignored-string")
		})
	})
}

func TestBadDefaults(t *testing.T) {

	Convey("Given I have struct with bad default duration", t, func() {

		c := &struct {
			A time.Duration `mapstructure:"BadDefaultDuration" desc:"" default:"toto"`
		}{}

		Convey("Then calling Initialize should panic", func() {
			So(func() { Initialize(c) }, ShouldPanicWith, "Unable to parse duration from: toto")
		})
	})

	Convey("Given I have struct with bad default int", t, func() {

		c := &struct {
			A int `mapstructure:"badDefaultInt" desc:"" default:"toto"`
		}{}

		Convey("Then calling Initialize should panic", func() {
			So(func() { Initialize(c) }, ShouldPanicWith, "Unable to parse int from: toto")
		})
	})

	Convey("Given I have struct with unsuported type", t, func() {

		c := &struct {
			A float64 `mapstructure:"badFloat" desc:"" default:"toto"`
		}{}

		Convey("Then calling Initialize should panic", func() {
			So(func() { Initialize(c) }, ShouldPanicWith, "Unsupported type: float64")
		})
	})

	Convey("Given I have struct with bad default bool slice", t, func() {

		c := &struct {
			A []bool `mapstructure:"badBools" desc:"" default:"a,b,c"`
		}{}

		Convey("Then calling Initialize should panic", func() {
			So(func() { Initialize(c) }, ShouldPanic)
			// So(func() { Initialize(c) }, ShouldPanicWith, `default value must a bool got: 'a'`)
		})
	})

	Convey("Given I have struct with bad default int slice", t, func() {

		c := &struct {
			A []int `mapstructure:"badInts" desc:"" default:"a,b,c"`
		}{}

		Convey("Then calling Initialize should panic", func() {
			So(func() { Initialize(c) }, ShouldPanic)
			// So(func() { Initialize(c) }, ShouldPanicWith, "default value must be an int. got 'a'")
		})
	})

	Convey("Given I have struct with bad default int slice", t, func() {

		c := &struct {
			A []net.IP `mapstructure:"badIPs" desc:"" default:"a,b,c"`
		}{}

		Convey("Then calling Initialize should panic", func() {
			So(func() { Initialize(c) }, ShouldPanic)
			// So(func() { Initialize(c) }, ShouldPanicWith, "default value must be an int. got 'a'")
		})
	})

	Convey("Given I have struct with unsuported slice", t, func() {

		c := &struct {
			A []float64 `mapstructure:"badFloats" desc:"" default:"a,b,c"`
		}{}

		Convey("Then calling Initialize should panic", func() {
			So(func() { Initialize(c) }, ShouldPanicWith, "Unsupported type: float64")
		})
	})
}
