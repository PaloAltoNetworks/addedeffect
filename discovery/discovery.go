package discovery

import (
	"encoding/json"
	"net/http"

	"github.com/Sirupsen/logrus"
)

// A PlatformInfo describes the Aporeto platform services.
type PlatformInfo struct {
	MidgardURL    string `json:"midgardURL"`
	PubsubService string `json:"pubsubService"`
	ZackURL       string `json:"zackURL"`
}

// RetrievePlatformInfo retrieves
func RetrievePlatformInfo(squallURL string) PlatformInfo {

	resp, err := http.Get(squallURL + "/systeminfos")

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"package":   "main",
			"squallURL": squallURL,
			"error":     err.Error(),
		}).Fatal("Unable to create request")
	}

	if resp.StatusCode != http.StatusOK {
		logrus.WithFields(logrus.Fields{
			"package":   "main",
			"squallURL": squallURL,
			"code":      resp.StatusCode,
		}).Fatal("Unable to get system info")
	}

	defer resp.Body.Close()
	info := PlatformInfo{}
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		logrus.WithFields(logrus.Fields{
			"package":   "main",
			"squallURL": squallURL,
			"code":      resp.StatusCode,
		}).Fatal("Unable to decode system info")
	}

	return info
}
