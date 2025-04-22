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

// ServerListResponse Response object for server list request.
// https://docs.api.eu01.stackit.cloud/documentation/iaas/version/v1#tag/Servers/operation/v1ListServersInProject
type ServerListResponse struct {
	Items *[]Server `json:"items"`
}

type Server struct {
	AffinityGroup    string                 `json:"affinityGroup"`
	AvailabilityZone string                 `json:"availabilityZone"`
	ID               string                 `json:"id"`
	Labels           map[string]interface{} `json:"labels"`
	MachineType      string                 `json:"machineType"`
	Name             string                 `json:"name"`
	Nics             []ServerNetwork        `json:"nics"`
	PowerStatus      string                 `json:"powerStatus"`
	Status           string                 `json:"status"`
}

// ServerNetwork Describes the object that matches servers to its networks.
type ServerNetwork struct {
	NetworkName string  `json:"networkName"`
	IPv4        *string `json:"ipv4,omitempty"`
	IPv6        *string `json:"ipv6,omitempty"`
	PublicIP    *string `json:"publicIp,omitempty"`
}

// PostgresFlexListResponse Response object for PostgresFlex list request.
// https://docs.api.eu01.stackit.cloud/documentation/postgres-flex-service/version/v2#tag/instance
type PostgresFlexListResponse struct {
	Items *[]PostgresFlex `json:"items"`
}

type PostgresFlex struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

// MongoDBFlexListResponse Response object for MongoDBFlex list request.
// https://docs.api.eu01.stackit.cloud/documentation/mongodb-flex-service/version/v2#tag/instance
type MongoDBFlexListResponse struct {
	Items *[]MongoDBFlex `json:"items"`
}

type MongoDBFlex struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}
