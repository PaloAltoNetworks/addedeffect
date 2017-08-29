package tracer

import (
	"crypto/x509"
	"time"

	"github.com/aporeto-inc/addedeffect/discovery"
	"github.com/opentracing/opentracing-go"

	jaeger "github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

// CloseRecorderHandler is the type of recorder closer handler
type CloseRecorderHandler func()

// ConfigureTracer configure the tracer for opentracing with the given platform and cert pool
func ConfigureTracer(pf *discovery.PlatformInfo, rootCAPool *x509.CertPool, serviceName string, insecureSkipVerify bool) (CloseRecorderHandler, opentracing.Tracer, error) {
	return func() {}, nil, nil
}

// ConfigureJaegerTracer returns a jaeger backed opentracing tracer.
func ConfigureJaegerTracer(pf *discovery.PlatformInfo, serviceName string) (CloseRecorderHandler, opentracing.Tracer, error) {

	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  pf.JaegerService,
		},
	}

	tracer, close, err := cfg.New(serviceName, jaegercfg.Logger(jaeger.NullLogger))
	if err != nil {
		return nil, nil, err
	}

	closer := func() {
		close.Close() // nolint: errcheck
	}

	opentracing.InitGlobalTracer(tracer)
	return closer, tracer, nil
}
