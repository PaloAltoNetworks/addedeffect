package tokensnip

import "github.com/aporeto-inc/addedeffect/tokenutils"

// Snip snips the given token from the given error.
func Snip(err error, token string) error {

	return tokenutils.Snip(err, token)
}
