package tagutils

import (
	"net/http"
	"strings"

	"go.aporeto.io/elemental"
	"go.aporeto.io/gaia"

	gaiaconstants "go.aporeto.io/gaia/constants"
)

const (
	prefixDynamicTag  = "$"
	prefixExpandedTag = "+"
	prefixMetadata    = "@"
)

// ValidateTagStrings validates the given string and check if it can be a valid value for a Tag.
func ValidateTagStrings(allowsReserved bool, strs ...string) error {

	errs := []error{}

	for _, s := range strs {

		if !allowsReserved && (strings.HasPrefix(s, prefixMetadata) || strings.HasPrefix(s, prefixDynamicTag) || strings.HasPrefix(s, prefixExpandedTag)) {
			errs = append(errs, elemental.NewError("Reserved Tag", "Tags with starting with an @, a $ or a + are reserved", "crud", http.StatusUnprocessableEntity))
			continue
		}

		t := &gaia.Tag{Value: s}
		if err := t.Validate(); err != nil {
			if e, ok := err.(elemental.Errors); ok {
				errs = append(errs, e...)
			} else {
				errs = append(errs, e)
			}
		}
	}

	if len(errs) > 0 {
		return elemental.NewErrors(errs...)
	}

	return nil
}

// ValidateMetadataStrings validates the given string and check if it can be a valid value for a Metadata.
func ValidateMetadataStrings(strs ...string) error {

	errs := []error{}

	for _, s := range strs {

		if strings.HasPrefix(s, gaiaconstants.AuthKey) {
			errs = append(errs, elemental.NewError("Invalid Metadata", "Prefix @auth: is reserved", "crud", http.StatusUnprocessableEntity))
			continue
		}

		t := &gaia.Tag{Value: s}
		if err := t.Validate(); err != nil {
			if e, ok := err.(elemental.Errors); ok {
				errs = append(errs, e...)
			} else {
				errs = append(errs, e)
			}
		}
	}

	if len(errs) > 0 {
		return elemental.NewErrors(errs...)
	}

	return nil
}
