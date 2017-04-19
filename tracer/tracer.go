package tracer

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"time"

	"github.com/aporeto-inc/addedeffect/discovery"
	"github.com/opentracing/opentracing-go"

	zipkin "github.com/openzipkin/zipkin-go-opentracing"
)

// CloseRecorderHandler is the type of recorder closer handler
type CloseRecorderHandler func()

// ConfigureTracer configure the tracer for opentracing with the given platform and cert pool
func ConfigureTracer(pf *discovery.PlatformInfo, rootCAPool *x509.CertPool, serviceName string, insecureSkipVerify bool) (CloseRecorderHandler, opentracing.Tracer, error) {

	if pf.ZipkinURL == "" {
		return nil, nil, nil
	}

	httpClientOption := zipkin.HTTPClient(&http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 100,
			TLSClientConfig: &tls.Config{
				RootCAs:            rootCAPool,
				InsecureSkipVerify: insecureSkipVerify,
			},
		},
	})

	collector, e := zipkin.NewHTTPCollector(pf.ZipkinURL+"/api/v1/spans", httpClientOption)

	if e != nil {
		return nil, nil, e
	}

	recorder := zipkin.NewRecorder(collector, false, "0.0.0.0:0", serviceName)
	tracer, e := zipkin.NewTracer(recorder, zipkin.ClientServerSameSpan(true), zipkin.TraceID128Bit(true))

	if e != nil {
		return nil, nil, e
	}

	closer := func() {
		collector.Close()
	}

	opentracing.InitGlobalTracer(tracer)
	return closer, tracer, nil
}
