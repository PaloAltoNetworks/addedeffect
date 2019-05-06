// Copyright 2019 Aporeto Inc.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tokenutils

import (
	"fmt"
	"strings"

	"go.aporeto.io/elemental"
)

// AudienceAny represents any operation, identity or namespace.
const AudienceAny = "*"

// UnsecureAudience extracts the audience list from a token string without
// verifying its validity. Only use or trust this after proper validation.
func UnsecureAudience(token string, modelManager elemental.ModelManager) (AudiencesList, error) {

	claims, err := UnsecureClaimsMap(token)
	if err != nil {
		return nil, fmt.Errorf("unable to extract audience: %s", err)
	}

	audStr, ok := claims["aud"].(string)
	if !ok || audStr == "" {
		return nil, nil
	}

	return ParseAudience(audStr, modelManager)
}

// ParseAudience parses the audience string and returns an AudiencesList.
func ParseAudience(audString string, modelManager elemental.ModelManager) (AudiencesList, error) {

	// TODO: In order to not invalidate all
	// currently issued tokens, if the audience
	// doesn't start with the correct prefix
	// we assume there is no audience.
	// This must be removed after a little while.
	if !strings.HasPrefix(audString, "aud:") {
		return nil, nil
	}

	auds := strings.Split(audString, ";")
	out := make(AudiencesList, len(auds))

	var parts []string
	var err error

	for i, a := range auds {

		parts = strings.Split(a, ":")
		if len(parts) != 4 {
			return nil, fmt.Errorf("invalid audience '%s': invalid number of parts", a)
		}

		// Validate operation
		ops := strings.Split(parts[1], ",")
		for _, o := range ops {
			if o == AudienceAny {
				continue
			}
			if _, err = elemental.ParseOperation(o); err != nil {
				return nil, fmt.Errorf("invalid audience '%s': %s", a, err)
			}
		}

		// Validate identity
		idents := strings.Split(parts[2], ",")
		for _, ident := range idents {
			if ident == AudienceAny {
				continue
			}
			if modelManager.IdentityFromCategory(ident).IsEmpty() {
				return nil, fmt.Errorf("invalid audience '%s': invalid identity '%s'", a, ident)
			}
		}

		out[i] = Audience{
			Operations: ops,
			Identities: idents,
			Namespaces: strings.Split(parts[3], ","),
		}
	}

	return out, nil
}

// AudiencesList is a list of audiences.
type AudiencesList []Audience

// Verify verifies at least one audience in the list is valid for the given operation, identity and namespace.
func (a AudiencesList) Verify(operation elemental.Operation, identity elemental.Identity, namespace string) bool {

	for _, item := range a {
		if item.Verify(operation, identity, namespace) {
			return true
		}
	}

	return false
}

func (a AudiencesList) String() string {

	parts := make([]string, len(a))
	for i, item := range a {
		parts[i] = item.String()
	}

	return strings.Join(parts, ";")
}

// Audience represents a parsed audience string.
type Audience struct {
	Operations []string
	Identities []string
	Namespaces []string
}

// Verify verifies the audience is valid for the given operation, identity and namespace.
func (a Audience) Verify(operation elemental.Operation, identity elemental.Identity, namespace string) bool {

	var operationOK, identityOK, namespaceOK bool

	for _, o := range a.Operations {
		if o == AudienceAny || o == string(operation) {
			operationOK = true
			break
		}
	}

	for _, i := range a.Identities {
		if i == AudienceAny || i == identity.Category {
			identityOK = true
			break
		}
	}

	for _, n := range a.Namespaces {
		if n == AudienceAny || n == namespace {
			namespaceOK = true
			break
		}
	}

	return operationOK && identityOK && namespaceOK
}

func (a Audience) String() string {
	return fmt.Sprintf("aud:%s:%s:%s",
		strings.Join(a.Operations, ","),
		strings.Join(a.Identities, ","),
		strings.Join(a.Namespaces, ","),
	)
}
