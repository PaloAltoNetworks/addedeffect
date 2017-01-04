package discovery

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// A PlatformInfo describes the Aporeto platform services.
type PlatformInfo struct {
	MidgardURL       string `json:"midgardURL"`
	PubSubService    string `json:"pubSubService"`
	ZackURL          string `json:"zackURL"`
	KairosDBURL      string `json:"kairosDBURL"`
	CassandraService string `json:"cassandraService"`
}

// RetrievePlatformInfo retrieves the Platform Information from a Squall URL.
func RetrievePlatformInfo(squallURL string) (*PlatformInfo, error) {

	resp, err := http.Get(squallURL + "/systeminfos")

	if err != nil {
		return nil, fmt.Errorf("Unable to create request: %d", resp.StatusCode)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unable to retrieve system info: status code %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	info := &PlatformInfo{}
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("Unable to decode system info: %s", err)
	}

	return info, nil
}
