package nsutils

import (
	"fmt"
	"strings"
)

// NamespaceSetter is an interface that allows to set a namespace
type NamespaceSetter interface {
	SetNamespace(string)
}

// NamespaceGetter is an interface that allows to get a namespace
type NamespaceGetter interface {
	GetNamespace() string
}

const nsSeparator = "/"

// ParentNamespaceFromString returns the parent namespace of a namespace
// It returns empty it the string is invalid
func ParentNamespaceFromString(namespace string) (string, error) {

	if namespace == nsSeparator {
		return "", nil
	}

	index := strings.LastIndex(namespace, nsSeparator)

	switch index {
	case -1:
		return "", fmt.Errorf("Invalid namespace name")
	case 0:
		return namespace[:index+1], nil
	default:
		return namespace[:index], nil
	}
}

// IsNamespaceRelatedToNamespace returns true if the given namespace is related to the given parent
func IsNamespaceRelatedToNamespace(ns string, parent string) bool {
	return IsNamespaceParentOfNamespace(ns, parent) || IsNamespaceChildrenOfNamespace(ns, parent) || (ns == parent && ns != "" && parent != "")
}

// IsNamespaceParentOfNamespace returns true if the given namespace is a parent of the given parent
func IsNamespaceParentOfNamespace(ns string, child string) bool {

	if len(ns) == 0 {
		return false
	}

	if ns == child {
		return false
	}

	if ns[len(ns)-1] != '/' {
		ns = ns + "/"
	}

	return strings.HasPrefix(child, ns)
}

// IsNamespaceChildrenOfNamespace returns true of the given ns is a children of the given parent.
func IsNamespaceChildrenOfNamespace(ns string, parent string) bool {

	if len(parent) == 0 {
		return false
	}

	if ns == parent {
		return false
	}

	if parent[len(parent)-1] != '/' {
		parent = parent + "/"
	}

	return strings.HasPrefix(ns, parent)
}

// NamespaceAncestorsNames returns the list of fully qualified namespaces
// in the hierarchy of a given namespace. It returns an empty
// array for the root namespace
func NamespaceAncestorsNames(namespace string) []string {

	if namespace == nsSeparator {
		return []string{}
	}

	parts := strings.Split(namespace, nsSeparator)
	sep := nsSeparator
	namespaces := []string{}

	for i := len(parts) - 1; i >= 2; i-- {
		namespaces = append(namespaces, sep+strings.Join(parts[1:i], sep))
	}

	namespaces = append(namespaces, sep)

	return namespaces
}
