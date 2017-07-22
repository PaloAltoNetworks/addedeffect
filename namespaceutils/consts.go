package namespaceutils

import (
	"github.com/aporeto-inc/elemental"
	"github.com/aporeto-inc/gaia/squallmodels/v1/golang"
)

var exportNamespacesObjects = map[string]elemental.Identity{
	// Policies
	"namespaceMappingPolicies":       squallmodels.NamespaceMappingPolicyIdentity,
	"networksAccessPolicies":         squallmodels.NetworkAccessPolicyIdentity,
	"fileAccessPolicies":             squallmodels.FileAccessPolicyIdentity,
	"APIAuthorizationPolicies":       squallmodels.APIAuthorizationPolicyIdentity,
	"enforcerProfileMappingPolicies": squallmodels.EnforcerProfileMappingPolicyIdentity,

	// Others
	"namespaces":       squallmodels.NamespaceIdentity,
	"externalServices": squallmodels.ExternalServiceIdentity,
	"filePaths":        squallmodels.FilePathIdentity,
	"enforcerProfiles": squallmodels.EnforcerProfileIdentity,
}
