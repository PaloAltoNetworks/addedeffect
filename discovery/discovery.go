package discovery

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// A PlatformInfo describes the Aporeto platform services.
type PlatformInfo struct {
	SquallURL          string   `json:"squall"`
	MidgardURL         string   `json:"midgard"`
	ZackURL            string   `json:"zack"`
	VinceURL           string   `json:"vince"`
	KairosDBURL        string   `json:"kairosdb"`
	PubSubServices     []string `json:"pubsub"`
	CassandraServices  []string `json:"cassandra"`
	MongoServices      []string `json:"mongo"`
	GoogleClientID     string   `json:"googleClientID"`
	GrayLogServer      string   `json:"graylog"`
	GrayLogID          string   `json:"graylogID"`
	CACert             string   `json:"CACert"`
	CACertKey          string   `json:"CACertKey"`
	ServicesCert       string   `json:"servicesCert"`
	ServicesCertKey    string   `json:"servicesCertKey"`
	ZackClientCert     string   `json:"zackClientCert"`
	ZackClientCertKey  string   `json:"zackClientCertKey"`
	VinceClientCert    string   `json:"vinceClientCert"`
	VinceClientCertKey string   `json:"vinceClientCertKey"`
}

// ServicesKeyPair decodes the services certificates using the given password.
func (p *PlatformInfo) ServicesKeyPair(password string) (tls.Certificate, error) {

	return tls.X509KeyPair([]byte(p.ServicesCert), []byte(p.ServicesCertKey))
}

// ZackClientKeyPair decodes the zack client certificates using the given password.
func (p *PlatformInfo) ZackClientKeyPair(password string) (tls.Certificate, error) {

	return tls.X509KeyPair([]byte(p.ZackClientCert), []byte(p.ZackClientCertKey))
}

// VinceClientKeyPair decodes the vince client certificates using the given password.
func (p *PlatformInfo) VinceClientKeyPair(password string) (tls.Certificate, error) {

	return tls.X509KeyPair([]byte(p.VinceClientCert), []byte(p.VinceClientCertKey))
}

func (p *PlatformInfo) String() string {

	return fmt.Sprintf(
		"<platform: squall:%s midgard:%s zack:%s vince:%s graylog:%s logid:%s>",
		p.SquallURL,
		p.MidgardURL,
		p.ZackURL,
		p.VinceURL,
		p.GrayLogServer,
		p.GrayLogID,
	)
}

// RootCAPool returns the a CA pool using the system certificates + the custom CA.
func (p *PlatformInfo) RootCAPool() (*x509.CertPool, error) {

	pool, err := x509.SystemCertPool()
	if err != nil {
		return nil, err
	}

	pool.AppendCertsFromPEM([]byte(p.CACert))

	return pool, nil
}

// ClientCAPool returns a a CA pool using only the custom CA.
func (p *PlatformInfo) ClientCAPool() (*x509.CertPool, error) {

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM([]byte(p.CACert))

	return pool, nil
}

// DiscoverPlatform retrieves the Platform Information from a Squall URL.
func DiscoverPlatform(cidURL string) (*PlatformInfo, error) {

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, cidURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Unable to create request %s: %s", cidURL, err)
	}
	req.Close = true

	try := 0
	var resp *http.Response

	for {
		resp, err = client.Do(req)
		if err == nil {
			break
		}

		<-time.After(3 * time.Second)
		try++
		if try > 20 {
			return nil, fmt.Errorf("Unable retrieve platform info after 1m. Aborting. error: %s", err)
		}
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unable to retrieve system info: status code %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	info := &PlatformInfo{}
	if err = json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("Unable to decode system info: %s", err)
	}

	return info, nil
}
