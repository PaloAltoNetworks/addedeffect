package tag

import "fmt"

// Split splits the given tag string into key and value without additional control.
// It will return the key, value and eventual error.
func Split(tag string) (key string, value string, err error) {

	l := len(tag)
	if l < 3 {
		err = fmt.Errorf("Invalid tag: invalid length '%s'", tag)
		return
	}

	if tag[0] == '=' {
		err = fmt.Errorf("Invalid tag: missing key '%s'", tag)
		return
	}

	if tag[l-1] == '=' {
		err = fmt.Errorf("Invalid tag: missing value '%s'", tag)
		return
	}

	for i := 0; i < l; i++ {
		if tag[i] == '=' {
			key = tag[:i]
			value = tag[i+1:]
			return
		}
	}

	err = fmt.Errorf("Invalid tag: missing equal symbol '%s'", tag)
	return
}

// SplitPtr splits the given tag string into key and value without additional control.
// It will populate the given string pointer key and value and will return an eventual error.
// This function is to be used when you really need optimization for a large tag splitting
// operation as it will save the allocation of the key and value strings.
func SplitPtr(tag string, key *string, value *string) (err error) {

	l := len(tag)
	if l < 3 {
		err = fmt.Errorf("Invalid tag: invalid length '%s'", tag)
		return
	}

	if tag[0] == '=' {
		err = fmt.Errorf("Invalid tag: missing key '%s'", tag)
		return
	}

	if tag[l-1] == '=' {
		err = fmt.Errorf("Invalid tag: missing value '%s'", tag)
		return
	}

	for i := 0; i < l; i++ {
		if tag[i] == '=' {
			*key = tag[:i]
			*value = tag[i+1:]
			return
		}
	}

	err = fmt.Errorf("Invalid tag: missing equal symbol '%s'", tag)
	return
}
