package discovery

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"net/http"
)

// A PlatformInfo describes the Aporeto platform services.
type PlatformInfo struct {
	MidgardURL           string `json:"midgardURL"`
	PubSubService        string `json:"pubSubService"`
	ZackURL              string `json:"zackURL"`
	KairosDBURL          string `json:"kairosDBURL"`
	CassandraService     string `json:"cassandraService"`
	VinceURL             string `json:"vinceURL"`
	GoogleClientID       string `json:"googleClientID"`
	CertificateAuthority string `json:"certificateAuthority"`
}

// RetrievePlatformInfo retrieves the Platform Information from a Squall URL.
func RetrievePlatformInfo(squallURL string, CAPool *x509.CertPool) (*PlatformInfo, error) {

	skip := false
	if CAPool == nil {
		skip = true
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:            CAPool,
				InsecureSkipVerify: skip,
			},
		},
	}

	req, err := http.NewRequest(http.MethodGet, squallURL+"/systeminfos", nil)
	if err != nil {
		return nil, fmt.Errorf("Unable to create request %s: %s", squallURL+"/systeminfos", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Unable to send request %s: %s", squallURL+"/systeminfos", err)
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
