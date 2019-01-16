package tokenutils

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// Snip snips the given token from the given error.
func Snip(err error, token string) error {

	if len(token) == 0 || err == nil {
		return err
	}

	return fmt.Errorf("%s",
		strings.Replace(
			err.Error(),
			token,
			"[snip]",
			-1),
	)
}

// UnsecureClaimsMap decodes the claims in the given JWT token without
// verifying its validity. Only use or trust this after proper validation.
func UnsecureClaimsMap(token string) (claims map[string]interface{}, err error) {

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid jwt: not enough segments")
	}

	data, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid jwt: %s", err)
	}

	claims = map[string]interface{}{}
	if err := json.Unmarshal(data, &claims); err != nil {
		return nil, fmt.Errorf("invalid jwt: %s", err)
	}

	return claims, nil
}

// SigAlg returns the signature used by the token
func SigAlg(token string) (string, error) {

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "", errors.New("invalid jwt: not enough segments")
	}

	data, err := base64.RawStdEncoding.DecodeString(parts[0])
	if err != nil {
		return "", fmt.Errorf("invalid jwt: %s", err)
	}

	header := struct {
		Alg string `json:"alg"`
	}{}

	if err := json.Unmarshal(data, &header); err != nil {
		return "", fmt.Errorf("invalid jwt: %s", err)
	}

	return header.Alg, nil
}

// ExtractQuota extracts the eventual quota from a token.
// Not that the token is not verified in the process,
// you the verification must be done before trusting
// the extracted quota value.
func ExtractQuota(token string) (int, error) {

	claims, err := UnsecureClaimsMap(token)
	if err != nil {
		return 0, err
	}

	quota, ok := claims["quota"]
	if !ok {
		return 0, nil
	}

	q, ok := quota.(float64)
	if !ok {
		return 0, fmt.Errorf("invalid quota claim type")
	}

	return int(q), nil
}
