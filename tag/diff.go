package tag

// Diff returns the list of tags that has been added and removed
// between given elemental.Identifiables.
func Diff(new, old []string) (added []string, removed []string) {

	for _, t := range new {
		if !Contains(old, t) {
			added = append(added, t)
		}
	}

	for _, t := range old {
		if !Contains(new, t) {
			removed = append(removed, t)
		}
	}

	return
}
