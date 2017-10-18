package discovery

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/aporeto-inc/tg/tglib"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// A PlatformInfo describes the Aporeto platform services.
type PlatformInfo struct {
	CidURL       string `json:"cid,omitempty"`
	CidPublicURL string `json:"cidPublic,omitempty"`

	CladURL       string `json:"clad,omitempty"`
	CladPublicURL string `json:"cladPublic,omitempty"`

	SquallURL       string `json:"squall,omitempty"`
	SquallPublicURL string `json:"squallPublic,omitempty"`

	MidgardURL       string `json:"midgard,omitempty"`
	MidgardPublicURL string `json:"midgardPublic,omitempty"`

	ZackURL       string `json:"zack,omitempty"`
	ZackPublicURL string `json:"zackPublic,omitempty"`

	VinceURL       string `json:"vince,omitempty"`
	VincePublicURL string `json:"vincePublic,omitempty"`

	JunonURL       string `json:"junon,omitempty"`
	JunonPublicURL string `json:"junonPublic,omitempty"`

	YuffieURL string `json:"yuffie,omitempty"`
	BarretURL string `json:"barret,omitempty"`

	HighwindURL       string `json:"highwind,omitempty"`
	HighwindPublicURL string `json:"highwindPublic,omitempty"`

	GeoIPURL string `json:"geoipURL,omitempty"`

	PubSubServices []string `json:"pubsub,omitempty"`
	MongoURL       string   `json:"mongo,omitempty"`
	InfluxDBURL    string   `json:"influxdb,omitempty"`

	GoogleClientID     string `json:"googleClientID,omitempty"`
	OpenTracingService string `json:"openTracingService,omitempty"`

	CACert                      string `json:"CACert,omitempty"`
	SystemCACert                string `json:"systemCACert,omitempty"`
	PublicServicesCert          string `json:"publicServicesCert,omitempty"`
	PublicServicesCertKey       string `json:"publicServicesCertKey,omitempty"`
	IssuingServiceClientCert    string `json:"issuingServiceClientCert,omitempty"`
	IssuingServiceClientCertKey string `json:"issuingServiceClientCertKey,omitempty"`
	DownloadManifestURL         string `json:"downloadManifestURL,omitempty"`

	License string `json:"license,omitempty"`
}

// IssuingServiceClientCertPair decodes the initial issuing client certificates using the given password.
func (p *PlatformInfo) IssuingServiceClientCertPair(password string) (tls.Certificate, error) {

	keyBlock, err := tglib.DecryptPrivateKeyPEM([]byte(p.IssuingServiceClientCertKey), password)
	if err != nil {
		return tls.Certificate{}, err
	}

	return tls.X509KeyPair([]byte(p.IssuingServiceClientCert), pem.EncodeToMemory(keyBlock))
}

// PublicServicesCertPair decodes the initial public server certificates using the given password.
func (p *PlatformInfo) PublicServicesCertPair(password string) (tls.Certificate, error) {

	keyBlock, err := tglib.DecryptPrivateKeyPEM([]byte(p.PublicServicesCertKey), password)
	if err != nil {
		return tls.Certificate{}, err
	}

	return tls.X509KeyPair([]byte(p.PublicServicesCert), pem.EncodeToMemory(keyBlock))
}

func (p *PlatformInfo) String() string {

	return fmt.Sprintf(
		"<platform: cid:%s squall:%s midgard:%s zack:%s vince:%s junon:%s yuffie:%s opentracing:%s>",
		p.CidURL,
		p.SquallURL,
		p.MidgardURL,
		p.ZackURL,
		p.VinceURL,
		p.JunonURL,
		p.YuffieURL,
		p.OpenTracingService,
	)
}

// Fields returns ready to be dump zap Fields fields.
func (p *PlatformInfo) Fields() []zapcore.Field {
	return []zapcore.Field{
		zap.String("cid", p.CidURL),
		zap.String("cidPublic", p.CidPublicURL),
		zap.String("clad", p.CladURL),
		zap.String("cladPublic", p.CladPublicURL),
		zap.String("squall", p.SquallURL),
		zap.String("squallPublic", p.SquallPublicURL),
		zap.String("midgard", p.MidgardURL),
		zap.String("midgardPublic", p.MidgardPublicURL),
		zap.String("zack", p.ZackURL),
		zap.String("zackPublic", p.ZackPublicURL),
		zap.String("vince", p.VinceURL),
		zap.String("vincePublic", p.VincePublicURL),
		zap.String("junon", p.JunonURL),
		zap.String("junonPublic", p.JunonPublicURL),
		zap.String("yuffie", p.YuffieURL),
		zap.String("barret", p.BarretURL),
		zap.String("opentracing", p.OpenTracingService),
		zap.String("mongo", p.MongoURL),
		zap.String("influxdb", p.InfluxDBURL),
		zap.String("geoip", p.GeoIPURL),
		zap.Strings("nats", p.PubSubServices),
	}
}

// PublicFields returns ready to be dump zap fields.
func (p *PlatformInfo) PublicFields() []zapcore.Field {
	return []zapcore.Field{
		zap.String("cid", p.CidURL),
		zap.String("clad", p.CladURL),
		zap.String("squall", p.SquallURL),
		zap.String("midgard", p.MidgardURL),
		zap.String("zack", p.ZackURL),
		zap.String("vince", p.VinceURL),
		zap.String("junon", p.JunonURL),
	}
}

// RootCAPool returns the a CA pool using the system certificates + the custom CA.
func (p *PlatformInfo) RootCAPool() (*x509.CertPool, error) {

	pool, err := x509.SystemCertPool()
	if err != nil {
		return nil, err
	}

	if ok := pool.AppendCertsFromPEM([]byte(p.CACert)); !ok {
		return nil, fmt.Errorf("Unable to create RootCAPool: cannot append public ca certificate: %s", p.CACert)
	}

	if ok := pool.AppendCertsFromPEM([]byte(p.SystemCACert)); !ok {
		return nil, fmt.Errorf("Unable to create RootCAPool: cannot append system ca certificate: %s", p.SystemCACert)
	}

	return pool, nil
}

// SystemCAPool returns a a CA pool using only the system CA.
func (p *PlatformInfo) SystemCAPool() (*x509.CertPool, error) {

	pool := x509.NewCertPool()

	if ok := pool.AppendCertsFromPEM([]byte(p.SystemCACert)); !ok {
		return nil, fmt.Errorf("Unable to create SystemCAPool: cannot append system ca certificate: %s", p.SystemCACert)
	}

	return pool, nil
}

// ClientCAPool returns a a CA pool using only the client CA and the system CA.
func (p *PlatformInfo) ClientCAPool() (*x509.CertPool, error) {

	pool := x509.NewCertPool()

	if ok := pool.AppendCertsFromPEM([]byte(p.CACert)); !ok {
		return nil, fmt.Errorf("Unable to create ClientCAPool: cannot append public ca certificate: %s", p.CACert)
	}

	if ok := pool.AppendCertsFromPEM([]byte(p.SystemCACert)); !ok {
		return nil, fmt.Errorf("Unable to create ClientCAPool: cannot append system ca certificate: %s", p.SystemCACert)
	}

	return pool, nil
}

// DiscoverPlatform retrieves the Platform Information from a Squall URL.
func DiscoverPlatform(cidURL string, rootCAPool *x509.CertPool, skip bool) (*PlatformInfo, error) {

	client := &http.Client{
		Timeout: 3 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:            rootCAPool,
				InsecureSkipVerify: skip,
			},
		},
	}

	req, err := http.NewRequest(http.MethodGet, cidURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Unable to create request %s: %s", cidURL, err)
	}

	try := 0
	var resp *http.Response

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	for {
		resp, err = client.Do(req)
		if err == nil {
			break
		}

		select {
		case <-time.After(3 * time.Second):
		case <-c:
			return nil, fmt.Errorf("Discovery aborted per os signal")
		}

		try++
		if try > 20 {
			return nil, fmt.Errorf("Unable retrieve platform info after 1m. Aborting. error: %s", err)
		}
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unable to retrieve system info: status code %d", resp.StatusCode)
	}

	defer resp.Body.Close() // nolint: errcheck
	info := &PlatformInfo{}
	if err = json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("Unable to decode system info: %s", err)
	}

	info.CidURL = cidURL

	return info, nil
}
