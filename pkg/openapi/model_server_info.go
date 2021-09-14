// Copyright (c) 2021 Cisco Systems, Inc. and its affiliates
// All rights reserved
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

/*
 * Fledge REST API
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

import "strconv"

// ServerInfo - server information
type ServerInfo struct {
	Name string `json:"name,omitempty"`

	Ip string `json:"ip,omitempty"`

	Port int32 `json:"port,omitempty"`

	Uuid string `json:"uuid,omitempty"`

	Tags string `json:"tags,omitempty"`

	Role string `json:"role,omitempty"`

	State string `json:"state,omitempty"`

	// TODO: FIXME - remove all of the variables in the below or refactor openapi spec
	Command []string `json:"command,omitempty"`

	// required by the controller to check what type of notification to send.
	// Ideally K8 will provide a new node and system will be able to determine it
	IsExistingNode bool `yaml:"is_existing_node" json:"is_existing_node"`

	//required by the controller to check if anything related to the node got updated.
	// Example - schema design change impacted this node so a notification is required to be sent.
	IsUpdated bool `yaml:"is_updated" json:"is_updated"`
}

// TODO: FIXME - remove this after openapi refactoring is done
func (s *ServerInfo) GetAddress() string {
	return s.Ip + ":" + strconv.Itoa(int(s.Port))
}
