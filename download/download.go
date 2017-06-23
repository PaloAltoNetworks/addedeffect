package download

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/blang/semver"
	uuid "github.com/satori/go.uuid"
)

// A Manifest represents a Download Manifest
type Manifest map[string]Component

// A Component represents a downloadable component.
type Component struct {
	Latest   string             `json:"latest"`
	Versions map[string]Version `json:"versions"`
}

// A Version represents a particular version of a Component.
type Version struct {
	Version   string             `json:"version"`
	ChangeLog string             `json:"changelog"`
	Variants  map[string]Variant `json:"variants"`
}

// A Variant represents a variant of a Version.
type Variant struct {
	URL       string `json:"url"`
	Signature string `json:"signature"`
}

// RetrieveManifest fetch the manifest at the given URL.
func RetrieveManifest(url string) (Manifest, error) {

	resp, err := http.Get(url + "?nocache=" + uuid.NewV4().String())
	if err != nil {
		return Manifest{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return Manifest{}, fmt.Errorf("Unable to download manifest: %s", resp.Status)
	}

	manifest := Manifest{}
	defer resp.Body.Close() // nolint: errcheck
	if err = json.NewDecoder(resp.Body).Decode(&manifest); err != nil {
		return Manifest{}, err
	}

	return manifest, nil
}

// Binary downloads and saves the binary at the given url to the given dest with the given mode.
func Binary(url string, dest string, mode os.FileMode) error {

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("Unable to find the request binary: %s", resp.Status)
	}

	defer resp.Body.Close() // nolint: errcheck
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(dest, bytes, mode)
}

// IsOutdated checks if the given current is outdated relatively to the second using semver.
func IsOutdated(current, available string) (bool, error) {

	semVerRemote, err := semver.Make(available)
	if err != nil {
		return false, err
	}

	semVerCurrent, err := semver.Make(strings.Replace(current, "v", "", 1))
	if err != nil {
		return false, err
	}

	return semVerRemote.Compare(semVerCurrent) > 0, nil
}
