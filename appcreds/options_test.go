package appcreds

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestOptions(t *testing.T) {

	Convey("calling newConfig should work", t, func() {
		cfg := newConfig()
		So(cfg.subnets, ShouldBeNil)
		So(cfg.maxValidity, ShouldEqual, 0)
	})

	Convey("calling OptionSubnets should work", t, func() {
		cfg := newConfig()
		OptionSubnets([]string{"1.2.3.4/4"})(&cfg)
		So(cfg.subnets, ShouldResemble, []string{"1.2.3.4/4"})
	})

	Convey("calling OptionMaxValidity should work", t, func() {
		cfg := newConfig()
		OptionMaxValidity(3 * time.Minute)(&cfg)
		So(cfg.maxValidity, ShouldEqual, 3*time.Minute)
	})
}
