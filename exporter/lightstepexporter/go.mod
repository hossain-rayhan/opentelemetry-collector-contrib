module github.com/open-telemetry/opentelemetry-collector-contrib/exporter/lightstepexporter

go 1.14

require (
	github.com/census-instrumentation/opencensus-proto v0.2.1
	github.com/golang/protobuf v1.3.5
	github.com/lightstep/opentelemetry-exporter-go v0.6.2
	github.com/shirou/gopsutil v2.20.4+incompatible // indirect
	github.com/stretchr/testify v1.5.1
	go.opentelemetry.io/collector v0.5.0
	go.opentelemetry.io/otel v0.6.0
	go.uber.org/zap v1.14.0
)

replace github.com/apache/thrift => github.com/apache/thrift v0.0.0-20161221203622-b2a4d4ae21c7
