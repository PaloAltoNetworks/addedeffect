package lombric

import (
	"fmt"
	"os"
	"reflect"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// CheckRequired is a helper to check if all the required
// parameters in viper are set.
func CheckRequired(keys ...string) {

	var fail bool
	for _, key := range keys {

		if reflect.DeepEqual(viper.Get(key), reflect.Zero(reflect.TypeOf(viper.Get(key))).Interface()) {
			fmt.Printf("Error: Parameter '--%s' is required.\n", key)
			fail = true
		}
	}

	if fail {
		fmt.Println()
		pflag.Usage()
		os.Exit(1)
	}
}
