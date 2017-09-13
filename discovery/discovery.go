package discovery

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// A PlatformInfo describes the Aporeto platform services.
type PlatformInfo struct {
	CidURL                string   `json:"cid,omitempty"`
	CidPublicURL          string   `json:"cidPublic,omitempty"`
	CladURL               string   `json:"clad,omitempty"`
	CladPublicURL         string   `json:"cladPublic,omitempty"`
	SquallURL             string   `json:"squall,omitempty"`
	SquallPublicURL       string   `json:"squallPublic,omitempty"`
	MidgardURL            string   `json:"midgard,omitempty"`
	MidgardPublicURL      string   `json:"midgardPublic,omitempty"`
	ZackURL               string   `json:"zack,omitempty"`
	ZackPublicURL         string   `json:"zackPublic,omitempty"`
	VinceURL              string   `json:"vince,omitempty"`
	VincePublicURL        string   `json:"vincePublic,omitempty"`
	JunonURL              string   `json:"junon,omitempty"`
	JunonPublicURL        string   `json:"junonPublic,omitempty"`
	YuffieURL             string   `json:"yuffie,omitempty"`
	BarretURL             string   `json:"barret,omitempty"`
	PubSubServices        []string `json:"pubsub,omitempty"`
	CassandraServices     []string `json:"cassandra,omitempty"`
	MongoServices         []string `json:"mongo,omitempty"`
	InfluxDBURL           string   `json:"influxdb,omitempty"`
	GoogleClientID        string   `json:"googleClientID,omitempty"`
	OpenTracingService    string   `json:"openTracingService,omitempty"`
	GeoIPURL              string   `json:"geoipURL,omitempty"`
	CACert                string   `json:"CACert,omitempty"`
	ServicesCert          string   `json:"servicesCert,omitempty"`
	ServicesCertKey       string   `json:"servicesCertKey,omitempty"`
	PublicServicesCert    string   `json:"publicServicesCert,omitempty"`
	PublicServicesCertKey string   `json:"publicServicesCertKey,omitempty"`
	ZackClientCert        string   `json:"zackClientCert,omitempty"`
	ZackClientCertKey     string   `json:"zackClientCertKey,omitempty"`
	VinceClientCert       string   `json:"vinceClientCert,omitempty"`
	VinceClientCertKey    string   `json:"vinceClientCertKey,omitempty"`
	SquallClientCert      string   `json:"squallClientCert,omitempty"`
	SquallClientCertKey   string   `json:"squallClientCertKey,omitempty"`
	TidusClientCert       string   `json:"tidusClientCert,omitempty"`
	TidusClientCertKey    string   `json:"tidusClientCertKey,omitempty"`
	MarleneClientCert     string   `json:"marleneClientCert,omitempty"`
	MarleneClientCertKey  string   `json:"marleneClientCertKey,omitempty"`
	TifaClientCert        string   `json:"tifaClientCert,omitempty"`
	TifaClientCertKey     string   `json:"tifaClientCertKey,omitempty"`
	BarretClientCert      string   `json:"barretClientCert,omitempty"`
	BarretClientCertKey   string   `json:"barretClientCertKey,omitempty"`
	MidgardClientCert     string   `json:"midgardClientCert,omitempty"`
	MidgardClientCertKey  string   `json:"midgardClientCertKey,omitempty"`
	GaiaVersion           string   `json:"gaiaVersion,omitempty"`
	SystemVersion         string   `json:"systemVersion,omitempty"`
	DownloadManifestURL   string   `json:"downloadManifestURL,omitempty"`
	ApoctlLinuxURL        string   `json:"apoctlLinuxURL,omitempty"`
	ApoctlWindowsURL      string   `json:"apoctlWindowsURL,omitempty"`
	ApoctlDarwinURL       string   `json:"apoctlDarwinURL,omitempty"`
	EnforcerdDockerImage  string   `json:"enforcerdDockerImage,omitempty"`
	KubesquallDockerImage string   `json:"kubesquallDockerImage,omitempty"`
	EnforcerdURL          string   `json:"enforcerdURL,omitempty"`
}

// ServicesKeyPair decodes the services certificates using the given password.
func (p *PlatformInfo) ServicesKeyPair(password string) ([]tls.Certificate, error) {

	ret := []tls.Certificate{}

	if p.PublicServicesCert != "" && p.PublicServicesCertKey != "" {
		externalKeyPair, err := loadCertificates([]byte(p.PublicServicesCert), []byte(p.PublicServicesCertKey), password)
		if err != nil {
			return nil, err
		}
		ret = append(ret, externalKeyPair)
	}

	internalKeyPair, err := loadCertificates([]byte(p.ServicesCert), []byte(p.ServicesCertKey), password)
	if err != nil {
		return nil, err
	}
	ret = append(ret, internalKeyPair)

	return ret, nil
}

// ZackClientKeyPair decodes the zack client certificates using the given password.
func (p *PlatformInfo) ZackClientKeyPair(password string) (tls.Certificate, error) {

	return loadCertificates([]byte(p.ZackClientCert), []byte(p.ZackClientCertKey), password)
}

// VinceClientKeyPair decodes the vince client certificates using the given password.
func (p *PlatformInfo) VinceClientKeyPair(password string) (tls.Certificate, error) {

	return loadCertificates([]byte(p.VinceClientCert), []byte(p.VinceClientCertKey), password)
}

// SquallClientKeyPair decodes the squall client certificates using the given password.
func (p *PlatformInfo) SquallClientKeyPair(password string) (tls.Certificate, error) {

	return loadCertificates([]byte(p.SquallClientCert), []byte(p.SquallClientCertKey), password)
}

// TidusClientKeyPair decodes the tidus client certificates using the given password.
func (p *PlatformInfo) TidusClientKeyPair(password string) (tls.Certificate, error) {

	return loadCertificates([]byte(p.TidusClientCert), []byte(p.TidusClientCertKey), password)
}

// TifaClientKeyPair decodes the tifa client certificates using the given password.
func (p *PlatformInfo) TifaClientKeyPair(password string) (tls.Certificate, error) {

	return loadCertificates([]byte(p.TifaClientCert), []byte(p.TifaClientCertKey), password)
}

// MidgardClientKeyPair decodes the midgard client certificates using the given password.
func (p *PlatformInfo) MidgardClientKeyPair(password string) (tls.Certificate, error) {

	return loadCertificates([]byte(p.MidgardClientCert), []byte(p.MidgardClientCertKey), password)
}

// MarleneClientKeyPair decodes the marlene client certificates using the given password.
func (p *PlatformInfo) MarleneClientKeyPair(password string) (tls.Certificate, error) {

	return loadCertificates([]byte(p.MarleneClientCert), []byte(p.MarleneClientCertKey), password)
}

// BarretClientKeyPair decodes the marlene client certificates using the given password.
func (p *PlatformInfo) BarretClientKeyPair(password string) (tls.Certificate, error) {

	return loadCertificates([]byte(p.BarretClientCert), []byte(p.BarretClientCertKey), password)
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
		zap.Strings("mongo", p.MongoServices),
		zap.Strings("cassandra", p.CassandraServices),
		zap.String("influxdb", p.InfluxDBURL),
		zap.String("geoip", p.GeoIPURL),
		zap.Strings("nats", p.PubSubServices),
		zap.String("system-version", p.SystemVersion),
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
		zap.String("system-version", p.SystemVersion),
	}
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
