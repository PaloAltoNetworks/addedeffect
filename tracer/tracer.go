// Copyright 2019 Aporeto Inc.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tracer

import (
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

// CloseRecorderHandler is the type of recorder closer handler
type CloseRecorderHandler func()

// ConfigureTracerWithURL returns a jaeger backed opentracing tracer from an URL.
func ConfigureTracerWithURL(tracerURL string, serviceName string) (CloseRecorderHandler, error) {

	if tracerURL == "" {
		return nil, nil
	}

	tracer, closeFunc, err := jaegercfg.Configuration{
		ServiceName: serviceName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  tracerURL,
		},
	}.NewTracer()

	if err != nil {
		return nil, err
	}

	opentracing.InitGlobalTracer(tracer)

	return func() { _ = closeFunc.Close() }, nil // nolint: errcheck
}
