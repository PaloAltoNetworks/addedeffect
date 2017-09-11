package envopt

import (
	"os"
	"regexp"
	"strings"
)

// Parse parses the environment variables to find any environment that matches docopt usage.
//
// It will use the given prefix to match the variables.
// Options are translated like:
//      --option-a: PREFIX_OPTION_A
// If the usage needs a value, it will add in the os.Args:
//      --options-a=${PREFIX_OPTION_A}
// If not, it will simply add:
//      --options-a
// And the value of environment variable will be simply ignored.
func Parse(prefix string, doc string) error {

	args := extractArguments(extractUsage(doc))

	done := map[string]bool{}

	for _, flag := range args {

		parts := strings.Split(flag, "=")
		option := parts[0]

		// Check if we already set the value
		if done[option] {
			continue
		}
		done[option] = true

		hasValue := hasValue(flag)

		env := prefix + "_" + strings.ToUpper(strings.Replace(strings.Replace(option, "--", "", -1), "-", "_", -1))

		if e := os.Getenv(env); e != "" {

			// if os.Args already contains the option, we skip
			if isAlreadySet(option) {
				continue
			}

			for _, v := range strings.Split(e, ",") {
				if hasValue {
					os.Args = append(os.Args, option+`=`+v)
				} else {
					if e == "false" {
						continue
					}
					os.Args = append(os.Args, option)
				}
			}
		}
	}

	return nil
}

func extractUsage(doc string) string {

	p := regexp.MustCompile(`(?im)^([^\n]*usage:[^\n]*\n?(?:[ \t].*?(?:\n|$))*)`)
	s := p.FindAllString(doc, -1)

	for i, v := range s {
		s[i] = strings.TrimSpace(v)
	}

	return s[0]
}

func extractArguments(usage string) []string {

	p := regexp.MustCompile(`--[\w-=<>]+`)
	return p.FindAllString(usage, -1)
}

func hasValue(option string) bool {

	p := regexp.MustCompile(`=(.*)`)
	s := p.FindStringSubmatch(option)

	return len(s) > 0
}

func isAlreadySet(option string) bool {

	for _, arg := range os.Args {
		k := strings.Split(arg, "=")
		if k[0] == option {
			return true
		}
	}

	return false
}
