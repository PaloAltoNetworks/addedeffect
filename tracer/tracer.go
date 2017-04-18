package tracer

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"

	"github.com/aporeto-inc/addedeffect/discovery"
	"github.com/opentracing/opentracing-go"

	zipkin "github.com/openzipkin/zipkin-go-opentracing"
)

// ConfigureTracer configure the tracer for opentracing with the given platform and cert pool
func ConfigureTracer(pf *discovery.PlatformInfo, rootCAPool *x509.CertPool, serviceName string) (zipkin.Collector, zipkin.SpanRecorder, opentracing.Tracer, error) {

	if pf.ZipkinURL == "" {
		return nil, nil, nil, nil
	}

	httpClientOption := zipkin.HTTPClient(&http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: rootCAPool,
			},
		},
	})

	collector, e := zipkin.NewHTTPCollector(pf.ZipkinURL+"/api/v1/spans", httpClientOption)

	if e != nil {
		return nil, nil, nil, e
	}

	recorder := zipkin.NewRecorder(collector, false, "0.0.0.0:0", serviceName)
	tracer, e := zipkin.NewTracer(recorder, zipkin.ClientServerSameSpan(true), zipkin.TraceID128Bit(true))

	if e != nil {
		return nil, nil, nil, e
	}

	opentracing.InitGlobalTracer(tracer)
	return collector, recorder, tracer, nil
}
