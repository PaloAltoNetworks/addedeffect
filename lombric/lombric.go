package lombric

import (
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Configurable is the interface of a configuration.
type Configurable interface {
	Prefix() string
	RequiredParameters() []string
}

// CidCommunicator is an extension to Configurable that asks for
// an initial ca to talk to cid.
type CidCommunicator interface {
	SetInitialCAPool(pool *x509.CertPool)
}

// Initialize does all the basic job of bindings
func Initialize(conf Configurable) {

	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		panic("Unable to bind flags: " + err.Error())
	}

	viper.SetEnvPrefix(conf.Prefix())
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	checkRequired(conf.RequiredParameters()...)

	if err := viper.Unmarshal(conf); err != nil {
		panic("Unable to unmarshal configuration: " + err.Error())
	}

	if c, ok := conf.(CidCommunicator); ok {

		pool, err := x509.SystemCertPool()
		if err != nil {
			panic("Unable to load system CA pool: " + err.Error())
		}

		if path := viper.GetString("cid-cacert"); path != "" {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				panic("Unable to read cid CA file: " + err.Error())
			}
			pool.AppendCertsFromPEM(data)
		}

		c.SetInitialCAPool(pool)
	}
}

func checkRequired(keys ...string) {

	var fail bool
	for _, key := range keys {

		if !viper.IsSet(key) || reflect.DeepEqual(viper.Get(key), reflect.Zero(reflect.TypeOf(viper.Get(key))).Interface()) {
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
