package tag

// RemoveDuplicates returns a tag strings list with all duplicates removed.
func RemoveDuplicates(tags []string) (result []string) {

	if len(tags) == 0 {
		return
	}

	seen := map[string]string{}
	for _, val := range tags {

		if _, ok := seen[val]; ok {
			continue
		}

		result = append(result, val)
		seen[val] = val
	}

	return result
}
