package lombric

import (
	"errors"
	"fmt"
	"os"
	"reflect"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func checkRequired(keys ...string) error {

	var failed bool
	for _, key := range keys {
		if !viper.IsSet(key) || reflect.DeepEqual(viper.Get(key), reflect.Zero(reflect.TypeOf(viper.Get(key))).Interface()) {
			fmt.Printf("Error: Parameter '--%s' is required.\n", key)
			failed = true
		}
	}

	if failed {
		return errors.New("missing required parameter")
	}

	return nil
}

func checkAllowedValues(allowedValues map[string][]string) error {

	var failed bool
	for key, values := range allowedValues {

		if !stringInSlice(viper.GetString(key), values) {
			fmt.Printf("Error: Parameter '--%s' must be one of %s. '%s' is invalid.\n", key, values, viper.GetString(key))
			failed = true
		}
	}

	if failed {
		return errors.New("wrong allowed values")
	}

	return nil
}

func fail() {
	fmt.Println()
	pflag.Usage()
	os.Exit(1)
}

func stringInSlice(str string, list []string) bool {

	for _, s := range list {
		if s == str {
			return true
		}
	}

	return false
}
