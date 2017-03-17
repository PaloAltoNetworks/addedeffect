package namespace

import (
	"github.com/aporeto-inc/elemental"
	squallmodels "github.com/aporeto-inc/gaia/squallmodels/current/golang"
)

var exportNamespacesObjects = map[string]elemental.Identity{
	// Policies
	"namespaceMappingPolicies":       squallmodels.NamespaceMappingPolicyIdentity,
	"networksAccessPolicies":         squallmodels.NetworkAccessPolicyIdentity,
	"fileAccessPolicies":             squallmodels.FileAccessPolicyIdentity,
	"APIAuthorizationPolicies":       squallmodels.APIAuthorizationPolicyIdentity,
	"enforcerProfileMappingPolicies": squallmodels.EnforcerProfileMappingPolicyIdentity,

	// Others
	"namespaces":            squallmodels.NamespaceIdentity,
	"externalServices":      squallmodels.ExternalServiceIdentity,
	"filePaths":             squallmodels.FilePathIdentity,
	"enforcerProfiles":      squallmodels.EnforcerProfileIdentity,
	"dependencyMapViews":    squallmodels.DependencyMapViewIdentity,
	"integrationIdentities": squallmodels.IntegrationIdentity,
}
