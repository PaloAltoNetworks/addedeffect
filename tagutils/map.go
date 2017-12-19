package tagutils

import (
	"fmt"
	"sort"
)

// ToMap converts the given tag list into a map.
// If the tags array contains invalid tags, ToMap will panic
func ToMap(tags []string) map[string]string {

	out := map[string]string{}

	for _, t := range tags {

		k, v, err := Split(t)
		if err != nil {
			panic("Invalid tag '%s' passed to ToMap")
		}

		out[k] = v
	}

	return out
}

// FromMap converts a map to a tag list.
func FromMap(m map[string]string) []string {
	r := []string{}
	for k, v := range m {
		r = append(r, fmt.Sprintf("%s=%s", k, v))
	}
	sort.Strings(r)
	return r
}
