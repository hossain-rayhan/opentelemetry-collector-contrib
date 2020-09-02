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

package kubeletstatsreceiver

import (
	"errors"
	"fmt"
	"time"

	"go.opentelemetry.io/collector/config/configmodels"
	"go.opentelemetry.io/collector/config/confignet"
	"k8s.io/client-go/kubernetes"

	"github.com/open-telemetry/opentelemetry-collector-contrib/internal/k8sconfig"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/kubeletstatsreceiver/kubelet"
)

var _ configmodels.Receiver = (*Config)(nil)

type Config struct {
	configmodels.ReceiverSettings `mapstructure:",squash"`
	kubelet.ClientConfig          `mapstructure:",squash"`
	confignet.TCPAddr             `mapstructure:",squash"`
	CollectionInterval            time.Duration `mapstructure:"collection_interval"`

	// ExtraMetadataLabels contains list of extra metadata that should be taken from /pods endpoint
	// and put as extra labels on metrics resource.
	// No additional metadata is fetched by default, so there are no extra calls to /pods endpoint.
	// Supported values include container.id and k8s.volume.type.
	ExtraMetadataLabels []kubelet.MetadataLabel `mapstructure:"extra_metadata_labels"`

	// MetricGroupsToCollect provides a list of metrics groups to collect metrics from.
	// "container", "pod", "node" and "volume" are the only valid groups.
	MetricGroupsToCollect []kubelet.MetricGroup `mapstructure:"metric_groups"`

	// Configuration of the Kubernetes API client.
	K8sAPIConfig *k8sconfig.APIConfig `mapstructure:"k8s_api_config"`
}

// getReceiverOptions returns receiverOptions is the config is valid,
// otherwise it will return an error.
func (cfg *Config) getReceiverOptions() (*receiverOptions, error) {
	err := kubelet.ValidateMetadataLabelsConfig(cfg.ExtraMetadataLabels)
	if err != nil {
		return nil, err
	}

	mgs, err := getMapFromSlice(cfg.MetricGroupsToCollect)
	if err != nil {
		return nil, err
	}

	var k8sAPIClient kubernetes.Interface
	if cfg.K8sAPIConfig != nil {
		k8sAPIClient, err = k8sconfig.MakeClient(*cfg.K8sAPIConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create K8s API client: %w", err)
		}
	}

	return &receiverOptions{
		name:                  cfg.Name(),
		collectionInterval:    cfg.CollectionInterval,
		extraMetadataLabels:   cfg.ExtraMetadataLabels,
		metricGroupsToCollect: mgs,
		k8sAPIClient:          k8sAPIClient,
	}, nil
}

// getMapFromSlice returns a set of kubelet.MetricGroup values from
// the provided list. Returns an err if invalid entries are encountered.
func getMapFromSlice(collect []kubelet.MetricGroup) (map[kubelet.MetricGroup]bool, error) {
	out := make(map[kubelet.MetricGroup]bool, len(collect))
	for _, c := range collect {
		if !kubelet.ValidMetricGroups[c] {
			return nil, errors.New("invalid entry in metric_groups")
		}
		out[c] = true
	}

	return out, nil
}
