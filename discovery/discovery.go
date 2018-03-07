package discovery

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	jsoniter "github.com/json-iterator/go"
)

// A PlatformInfo describes the Aporeto platform services.
type PlatformInfo struct {
	PublicAPIURL       string `json:"public-api"`
	PublicCladURL      string `json:"public-clad"`
	WutaiURL           string `json:"wutai"`
	SquallURL          string `json:"squall"`
	CladURL            string `json:"clad"`
	MidgardURL         string `json:"midgard"`
	ZackURL            string `json:"zack"`
	VinceURL           string `json:"vince"`
	JunonURL           string `json:"junon"`
	YuffieURL          string `json:"yuffie"`
	BarretURL          string `json:"barret"`
	HighwindURL        string `json:"highwind"`
	GeoIPURL           string `json:"geoipURL"`
	PubSubService      string `json:"pubsub"`
	MongoURL           string `json:"mongo"`
	InfluxDBURL        string `json:"influxdb"`
	OpenTracingService string `json:"openTracingService"`

	License string `json:"license"`
}

func (p *PlatformInfo) String() string {

	return fmt.Sprintf(
		"<platform: wutai:%s squall:%s midgard:%s zack:%s vince:%s junon:%s yuffie:%s opentracing:%s>",
		p.WutaiURL,
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
		zap.String("squall", p.SquallURL),
		zap.String("midgard", p.MidgardURL),
		zap.String("zack", p.ZackURL),
		zap.String("vince", p.VinceURL),
		zap.String("junon", p.JunonURL),
		zap.String("yuffie", p.YuffieURL),
		zap.String("barret", p.BarretURL),
		zap.String("opentracing", p.OpenTracingService),
		zap.String("mongo", p.MongoURL),
		zap.String("influxdb", p.InfluxDBURL),
		zap.String("geoip", p.GeoIPURL),
		zap.String("nats", p.PubSubService),
	}
}

// Discover retrieves the Platform Information from a Cid URL. In case of communication error it will retry
// every 3 seconds until the given context is canceled.
func Discover(ctx context.Context, cidURL string, rootCAPool *x509.CertPool, skip bool) (*PlatformInfo, error) {

	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:            rootCAPool,
				InsecureSkipVerify: skip,
			},
		},
	}

	var resp *http.Response

	for {
		req, err := http.NewRequest(http.MethodGet, cidURL, nil)
		if err != nil {
			return nil, fmt.Errorf("unable to create request %s: %s", cidURL, err)
		}

		resp, err = client.Do(req.WithContext(ctx))
		if err == nil {
			break
		}

		zap.L().Warn("Unable to send request to cid. Retrying", zap.Error(err))

		select {
		case <-time.After(3 * time.Second):
		case <-ctx.Done():
			return nil, fmt.Errorf("discovery aborted per context: %s", ctx.Err())
		}
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unable to retrieve system info: status code %d", resp.StatusCode)
	}

	defer resp.Body.Close() // nolint: errcheck
	info := &PlatformInfo{}
	if err := jsoniter.ConfigCompatibleWithStandardLibrary.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("unable to decode system info: %s", err)
	}

	return info, nil
}
