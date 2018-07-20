package download

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"github.com/blang/semver"
	"go.aporeto.io/addedeffect/retry"
	"go.uber.org/zap"
)

// A Manifest represents a Download Manifest
type Manifest map[string]Component

// RetrieveManifest fetch the manifest at the given URL.
func RetrieveManifest(ctx context.Context, url string) (Manifest, error) {

	out, err := retry.Retry(
		ctx,
		func() (interface{}, error) {
			resp, err := http.Get(fmt.Sprintf("%s?nocache=%d", url, rand.Int()))
			if err != nil {
				return nil, err
			}

			if resp.StatusCode != http.StatusOK {
				return nil, fmt.Errorf("Unable to download manifest: %s", resp.Status)
			}

			return resp, nil
		},
		func(err error) error {
			zap.L().Warn("Unable to download manifest. retrying in 3s", zap.Error(err))
			return nil
		},
	)

	if err != nil {
		return Manifest{}, err
	}

	resp := out.(*http.Response)

	manifest := Manifest{}
	defer resp.Body.Close() // nolint: errcheck
	if err = json.NewDecoder(resp.Body).Decode(&manifest); err != nil {
		return Manifest{}, err
	}

	return manifest, nil
}

// A Component represents a downloadable component.
type Component struct {
	Latest   string             `json:"latest"`
	Versions map[string]Version `json:"versions"`
}

// NewComponent returns a new Component.
func NewComponent(latest string) Component {
	return Component{
		Latest:   latest,
		Versions: map[string]Version{},
	}
}

// A Version represents a particular version of a Component.
type Version struct {
	Version   string             `json:"version"`
	ChangeLog string             `json:"changelog"`
	Variants  map[string]Variant `json:"variants"`
}

// NewVersion returns a new Version.
func NewVersion(version, changelog string) Version {
	return Version{
		Version:   version,
		ChangeLog: changelog,
		Variants:  map[string]Variant{},
	}
}

// A Variant represents a variant of a Version.
type Variant struct {
	URL       string `json:"url"`
	Signature string `json:"signature"`
}

// NewVariant returns a new Variant.
func NewVariant(url, signature string) Variant {
	return Variant{
		URL:       url,
		Signature: signature,
	}
}

// Binary downloads and saves the binary at the given url to the given dest with the given mode.
func Binary(ctx context.Context, url string, dest string, mode os.FileMode, signature string) error {

	out, err := retry.Retry(
		ctx,
		func() (interface{}, error) {
			resp, err := http.Get(url)
			if err != nil {
				return nil, err
			}

			if resp.StatusCode != 200 {
				return nil, fmt.Errorf("Unable to find the request binary: %s", resp.Status)
			}

			return resp, nil
		},
		func(err error) error {
			zap.L().Warn("Unable to download binary. retrying in 3s", zap.Error(err))
			return nil
		},
	)

	if err != nil {
		return err
	}

	resp := out.(*http.Response)

	defer resp.Body.Close() // nolint: errcheck
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if signature != "" {

		h := sha1.New()
		if _, err = h.Write(data); err != nil {
			return err
		}

		if fmt.Sprintf("%x", h.Sum(nil)) != signature {
			return fmt.Errorf("Inavlid signature")
		}
	}

	return ioutil.WriteFile(dest, data, mode)
}

// IsOutdated checks if the given current is outdated relatively to the second using semver.
func IsOutdated(current, available string) (bool, error) {

	semVerRemote, err := semver.Make(strings.Replace(available, "v", "", 1))
	if err != nil {
		return false, err
	}

	semVerCurrent, err := semver.Make(strings.Replace(current, "v", "", 1))
	if err != nil {
		return false, err
	}

	return semVerRemote.Compare(semVerCurrent) > 0, nil
}
