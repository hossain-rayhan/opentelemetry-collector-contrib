module github.com/open-telemetry/opentelemetry-collector-contrib/exporter/signalfxexporter

go 1.14

require (
	github.com/Azure/go-autorest/autorest/adal v0.9.0 // indirect
	github.com/census-instrumentation/opencensus-proto v0.3.0
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.4.2
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/common v0.0.0-00010101000000-000000000000
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/k8sclusterreceiver v0.0.0-00010101000000-000000000000
	github.com/signalfx/com_signalfx_metrics_protobuf v0.0.1
	github.com/stretchr/testify v1.6.1
	go.opentelemetry.io/collector v0.9.1-0.20200901221426-ec327358d634
	go.uber.org/zap v1.15.0
	google.golang.org/protobuf v1.25.0
)

replace github.com/open-telemetry/opentelemetry-collector-contrib/internal/common => ../../internal/common

replace github.com/open-telemetry/opentelemetry-collector-contrib/internal/k8sconfig => ../../internal/k8sconfig

replace github.com/open-telemetry/opentelemetry-collector-contrib/receiver/k8sclusterreceiver => ../../receiver/k8sclusterreceiver
