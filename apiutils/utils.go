package apiutils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// ServiceVersion holds the version of a servie
type ServiceVersion struct {
	Libs    map[string]string
	Version string
	Sha     string
}

// GetServiceVersions returns the version of the services.
func GetServiceVersions(api string) (map[string]ServiceVersion, error) {

	resp, err := http.Get(fmt.Sprintf("%s/_meta/version", api))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Bad response status: %s", resp.Status)
	}

	out := map[string]ServiceVersion{}

	defer resp.Body.Close() // nolint: errcheck
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}

	return out, nil
}

// GetPublicCA returns the public CA used by the api.
func GetPublicCA(api string) ([]byte, error) {

	resp, err := http.Get(fmt.Sprintf("%s/_meta/ca", api))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Bad response status: %s", resp.Status)
	}

	defer resp.Body.Close() // nolint: errcheck
	return ioutil.ReadAll(resp.Body)
}

// GetManifestURL returns the url of the manifest.
func GetManifestURL(api string) ([]byte, error) {

	resp, err := http.Get(fmt.Sprintf("%s/_meta/manifest", api))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Bad response status: %s", resp.Status)
	}

	defer resp.Body.Close() // nolint: errcheck
	return ioutil.ReadAll(resp.Body)
}
