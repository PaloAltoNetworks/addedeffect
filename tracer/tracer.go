package tracer

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/aporeto-inc/addedeffect/discovery"
	"github.com/opentracing/opentracing-go"

	zipkin "github.com/openzipkin/zipkin-go-opentracing"
)

// ConfigureTracer configure the tracer for opentracing with the given platform and cert pool
func ConfigureTracer(pf *discovery.PlatformInfo, rootCAPool *x509.CertPool, serviceName string) {

	if pf.ZipkinURL == "" {
		return
	}

	logrus.WithFields(logrus.Fields{
		"services": pf.ZipkinURL,
	}).Info("Connecting to zipkin...")

	httpClientOption := zipkin.HTTPClient(&http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: rootCAPool,
			},
		},
	})

	collector, e := zipkin.NewHTTPCollector(pf.ZipkinURL+"/api/v1/spans", httpClientOption)
	if e != nil {
		logrus.WithFields(logrus.Fields{
			"error":   e.Error(),
			"service": pf.ZipkinURL,
		}).Fatal("Unable to connect to zipkin server")
	}
	recorder := zipkin.NewRecorder(collector, false, "0.0.0.0:0", serviceName)
	tracer, e := zipkin.NewTracer(recorder, zipkin.ClientServerSameSpan(true), zipkin.TraceID128Bit(true))

	if e != nil {
		logrus.WithFields(logrus.Fields{
			"error":   e.Error(),
			"service": pf.ZipkinURL,
		}).Fatal("Unable to create tracer")
	}

	opentracing.InitGlobalTracer(tracer)
}
