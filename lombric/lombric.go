package lombric

import (
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const enabledKey = "true"

// Configurable is the interface of a configuration.
type Configurable interface {
	Prefix() string
}

// CidCommunicator is an extension to Configurable that asks for
// an initial ca to talk to cid.
type CidCommunicator interface {
	SetInitialCAPool(pool *x509.CertPool)
}

// VersionPrinter is an extension to Configurable that can print its version.
type VersionPrinter interface {
	PrintVersion()
}

// Initialize does all the basic job of bindings
func Initialize(conf Configurable) {

	requiredFlags, secretFlags, allowedValues := installFlags(conf)

	pflag.VisitAll(func(f *pflag.Flag) {
		var v interface{}
		var err error
		switch f.Value.Type() {
		case "stringSlice":
			v, err = pflag.CommandLine.GetStringSlice(f.Name)
		case "boolSlice":
			v, err = pflag.CommandLine.GetBoolSlice(f.Name)
		case "intSlice":
			v, err = pflag.CommandLine.GetIntSlice(f.Name)
		case "ipSlice":
			v, err = pflag.CommandLine.GetIPSlice(f.Name)
		}

		if err != nil {
			panic("Unable to trick viper with the defaults: %s" + err.Error())
		}

		viper.SetDefault(f.Name, v)
	})

	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		panic("Unable to bind flags: " + err.Error())
	}

	viper.SetEnvPrefix(conf.Prefix())
	viper.SetTypeByDefaultValue(true)
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	if vp, ok := conf.(VersionPrinter); ok && viper.GetBool("version") {
		vp.PrintVersion()
		os.Exit(0)
	}

	checkRequired(requiredFlags...)
	checkAllowedValues(allowedValues)

	if err := viper.Unmarshal(conf); err != nil {
		panic("Unable to unmarshal configuration: " + err.Error())
	}

	if c, ok := conf.(CidCommunicator); ok {

		pool, err := x509.SystemCertPool()
		if err != nil {
			panic("Unable to load system CA pool: " + err.Error())
		}

		path := viper.GetString("cid-cacert")
		if path == "" {
			path = viper.GetString("api-cacert")
		}

		if path != "" {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				panic("Unable to read cid CA file: " + err.Error())
			}
			pool.AppendCertsFromPEM(data)
		}

		c.SetInitialCAPool(pool)
	}

	// Clean up all secrets
	for _, key := range secretFlags {
		env := strings.Replace(strings.ToUpper(conf.Prefix()+"_"+key), "-", "_", -1)
		if err := os.Unsetenv(env); err != nil {
			panic("Unable to unset secret env variable " + env)
		}
	}
}

func deepFields(ift reflect.Type) ([]reflect.StructField, []string) {

	fields := make([]reflect.StructField, 0)
	overrides := []string{}

	for i := 0; i < ift.NumField(); i++ {
		field := ift.Field(i)

		switch field.Type.Kind() {
		case reflect.Struct:
			if overrideString := field.Tag.Get("override"); overrideString != "" {
				overrides = append(overrides, overrideString)
			}

			f, o := deepFields(field.Type)

			overrides = append(overrides, o...)
			fields = append(fields, f...)

		default:
			fields = append(fields, field)
		}
	}

	return fields, overrides
}

func installFlags(conf Configurable) (requiredFlags []string, secretFlags []string, allowedValues map[string][]string) {

	t := reflect.ValueOf(conf).Elem().Type()

	fields, overrides := deepFields(t)

	defaultOverrides := map[string]string{}
	allowedValues = map[string][]string{}

	for _, raw := range overrides {
		for _, innerOverride := range strings.Split(raw, ",") {
			parts := strings.SplitN(innerOverride, "=", 2)
			defaultOverrides[parts[0]] = parts[1]
		}
	}

	for _, field := range fields {

		key := field.Tag.Get("mapstructure")
		if key == "" || key == "-" {
			continue
		}

		description := field.Tag.Get("desc")

		var def string
		if o, ok := defaultOverrides[key]; ok {
			if o == "-" {
				continue
			}
			def = o
		} else {
			def = field.Tag.Get("default")
		}

		if field.Tag.Get("secret") == enabledKey {
			secretFlags = append(secretFlags, key)
		}

		if field.Tag.Get("required") == enabledKey {
			requiredFlags = append(requiredFlags, key)
			description += " [required]"
		}

		if field.Type.Kind() != reflect.Slice {

			switch field.Type.Name() {

			case "bool":
				pflag.Bool(key, def == enabledKey, description)

			case "string":
				if allowed := field.Tag.Get("allowed"); allowed != "" {
					allowedValues[key] = strings.Split(allowed, ",")
					description += fmt.Sprintf(" [allowed: %s]", allowed)
				}
				pflag.String(key, def, description)

			case "Duration":
				if def == "" {
					pflag.Duration(key, 0, description)
					break
				}
				d, err := time.ParseDuration(def)
				if err != nil {
					panic("Unable to parse duration from: " + def)
				}
				pflag.Duration(key, d, description)

			case "int":
				if def == "" {
					pflag.Int(key, 0, description)
					break
				}
				d, err := strconv.Atoi(def)
				if err != nil {
					panic("Unable to parse int from: " + def)
				}
				pflag.Int(key, d, description)

			default:
				panic("Unsupported type: " + field.Type.Name())
			}

		} else {

			switch field.Type.Elem().Name() {

			case "string":
				sdef := strings.Split(def, ",")
				pflag.StringSlice(key, sdef, description)

			default:
				panic("Unsupported type: " + field.Type.Name())
			}
		}
	}

	if _, ok := conf.(VersionPrinter); ok {
		pflag.BoolP("version", "v", false, "Display the version")
	}

	pflag.Parse()

	return requiredFlags, secretFlags, allowedValues
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

func checkAllowedValues(allowedValues map[string][]string) {

	var fail bool

	for key, values := range allowedValues {

		if !stringInSlice(viper.GetString(key), values) {
			fmt.Printf("Error: Parameter '--%s' must be one of %s. '%s' is invalid.\n", key, values, viper.GetString(key))
			fail = true
		}
	}

	if fail {
		fmt.Println()
		pflag.Usage()
		os.Exit(1)
	}
}

func stringInSlice(str string, list []string) bool {

	for _, s := range list {
		if s == str {
			return true
		}
	}

	return false
}
