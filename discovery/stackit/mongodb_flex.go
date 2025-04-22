// Copyright 2020 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package stackit

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/discovery/refresh"
	"github.com/prometheus/prometheus/discovery/targetgroup"
	"io"
	"log/slog"
	"net/http"
)

const (
	stackitMongoAPIEndpoint         = "https://mongodb-flex-service.api.stackit.cloud"
	stackitMongoPrometheusProxyHost = "mongodb-prom-proxy.api.stackit.cloud"
	stackitMongoRunningState        = "READY"
	stackitLabelMongoDBFlexID       = stackitLabelPrefix + "mongodb_flex_id"
	stackitLabelMongoDBFlexName     = stackitLabelPrefix + "mongodb_flex_name"
	stackitLabelMongoDBFlexStatus   = stackitLabelPrefix + "mongodb_flex_status"
	stackitLabelMongoDBFlexRegion   = stackitLabelPrefix + "mongodb_flex_region"
)

type mongoDBFlexDiscovery struct {
	*refresh.Discovery
	httpClient  *http.Client
	logger      *slog.Logger
	apiEndpoint string
	project     string
	region      string
}

// newMongoDBFlexDiscovery creates a new MongoDB Flex service discovery instance.
// It sets up the HTTP client, authentication, and endpoint configuration using shared logic.
func newMongoDBFlexDiscovery(conf *SDConfig, logger *slog.Logger) (*mongoDBFlexDiscovery, error) {
	base, err := setupDiscoveryBase(
		conf,
		"STACKIT MongoDBFlex API",
		func(_ *SDConfig) string {
			return stackitMongoAPIEndpoint
		},
	)
	if err != nil {
		return nil, err
	}

	return &mongoDBFlexDiscovery{
		project:     conf.Project,
		region:      conf.Region,
		apiEndpoint: base.apiEndpoint,
		httpClient:  base.httpClient,
		logger:      logger,
	}, nil
}

func (p *mongoDBFlexDiscovery) refresh(ctx context.Context) ([]*targetgroup.Group, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/v2/projects/%s/regions/%s/instances", p.apiEndpoint, p.project, p.region),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	res, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		errorMessage, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("unexpected status code: %d, message: %s", res.StatusCode, string(errorMessage))
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	var mongoDBFlexListResponse *MongoDBFlexListResponse

	if err := json.Unmarshal(body, &mongoDBFlexListResponse); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	if mongoDBFlexListResponse == nil || mongoDBFlexListResponse.Items == nil || len(*mongoDBFlexListResponse.Items) == 0 {
		return []*targetgroup.Group{{Source: "stackit", Targets: []model.LabelSet{}}}, nil
	}

	targets := make([]model.LabelSet, 0, len(*mongoDBFlexListResponse.Items))
	for _, instance := range *mongoDBFlexListResponse.Items {
		if instance.Status != stackitMongoRunningState {
			continue
		}

		labels := model.LabelSet{
			stackitLabelRole:              model.LabelValue(RoleMongoDbFlex),
			stackitLabelProject:           model.LabelValue(p.project),
			stackitLabelMongoDBFlexID:     model.LabelValue(instance.Id),
			stackitLabelMongoDBFlexName:   model.LabelValue(instance.Name),
			stackitLabelMongoDBFlexStatus: model.LabelValue(instance.Status),
			stackitLabelMongoDBFlexRegion: model.LabelValue(p.region),
			model.InstanceLabel:           model.LabelValue(instance.Id),
			model.SchemeLabel:             model.LabelValue("https"),
			model.MetricsPathLabel:        model.LabelValue(fmt.Sprintf("/v2/projects/%s/regions/%s/instances/%s/metrics", p.project, p.region, instance.Id)),
			model.AddressLabel:            stackitMongoPrometheusProxyHost,
		}

		targets = append(targets, labels)
	}

	return []*targetgroup.Group{{Source: "stackit", Targets: targets}}, nil
}
