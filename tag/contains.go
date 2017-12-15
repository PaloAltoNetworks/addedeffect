package tag

// Contains returns true if the given tag strings list contains the given tag string.
func Contains(tags []string, value string) bool {

	for _, i := range tags {
		if i == value {
			return true
		}
	}

	return false
}
