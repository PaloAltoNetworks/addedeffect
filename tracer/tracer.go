package tracer

import (
	"time"

	"github.com/aporeto-inc/addedeffect/discovery"
	"github.com/opentracing/opentracing-go"

	jaeger "github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

// CloseRecorderHandler is the type of recorder closer handler
type CloseRecorderHandler func()

// ConfigureTracerWithURL returns a jaeger backed opentracing tracer from an URL.
func ConfigureTracerWithURL(tracerURL string, serviceName string) (CloseRecorderHandler, error) {

	if tracerURL == "" {
		return nil, nil
	}

	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  tracerURL,
		},
	}

	tracer, close, err := cfg.New(serviceName, jaegercfg.Logger(jaeger.NullLogger))
	if err != nil {
		return nil, err
	}

	opentracing.InitGlobalTracer(tracer)

	return func() { close.Close() }, nil // nolint: errcheck
}

// ConfigureTracer returns a jaeger backed opentracing tracer.
func ConfigureTracer(pf *discovery.PlatformInfo, serviceName string) (CloseRecorderHandler, error) {

	return ConfigureTracerWithURL(pf.OpenTracingService, serviceName)
}
