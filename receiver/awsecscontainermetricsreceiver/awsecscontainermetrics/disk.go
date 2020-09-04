// Copyright 2020, OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package awsecscontainermetrics

import (
	"time"

	metricspb "github.com/census-instrumentation/opencensus-proto/gen-go/metrics/v1"
)

func diskMetrics(prefix string, stats *DiskStats, labelKeys []*metricspb.LabelKey, labelValues []*metricspb.LabelValue) []*metricspb.Metric {
	readBytes, writeBytes := extractStorageUsage(stats)

	return applyCurrentTime([]*metricspb.Metric{
		intGauge(prefix+"disk.storage_read_bytes", "Bytes", &readBytes, labelKeys, labelValues),
		intGauge(prefix+"disk.storage_write_bytes", "Bytes", &writeBytes, labelKeys, labelValues),
	}, time.Now())
}
