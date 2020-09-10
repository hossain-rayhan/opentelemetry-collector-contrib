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

package awsecscontainermetricsreceiver

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/exporter/exportertest"
	"go.uber.org/zap"
)

type fakeRestClient struct {
}

func (f fakeRestClient) EndpointResponse() ([]byte, []byte, error) {
	taskStats, err := ioutil.ReadFile("../testdata/task_stats.json")
	if err != nil {
		return nil, nil, err
	}
	taskMetadata, err := ioutil.ReadFile("../testdata/task_metadata.json")
	if err != nil {
		return nil, nil, err
	}
	return taskStats, taskMetadata, nil
}

func TestReceiver(t *testing.T) {
	cfg := createDefaultConfig().(*Config)
	metricsReceiver, err := New(
		zap.NewNop(),
		cfg,
		exportertest.NewNopMetricsExporter(),
		&fakeRestClient{},
	)

	require.NoError(t, err)
	require.NotNil(t, metricsReceiver)

	r := metricsReceiver.(*awsEcsContainerMetricsReceiver)
	ctx := context.Background()

	err = r.Start(ctx, componenttest.NewNopHost())
	require.NoError(t, err)

	err = r.Shutdown(ctx)
	require.NoError(t, err)
}

func TestCollectDataFromEndpoint(t *testing.T) {
	cfg := createDefaultConfig().(*Config)
	metricsReceiver, err := New(
		zap.NewNop(),
		cfg,
		new(exportertest.SinkMetricsExporter),
		&fakeRestClient{},
	)

	require.NoError(t, err)
	require.NotNil(t, metricsReceiver)

	r := metricsReceiver.(*awsEcsContainerMetricsReceiver)
	ctx := context.Background()

	err = r.collectDataFromEndpoint(ctx, "")
	require.NoError(t, err)
}
