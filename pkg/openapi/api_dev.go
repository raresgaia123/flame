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

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// A DevApiController binds http requests to an api service and writes the service results to the http response
type DevApiController struct {
	service DevApiServicer
}

// NewDevApiController creates a default api controller
func NewDevApiController(s DevApiServicer) Router {
	return &DevApiController{service: s}
}

// Routes returns all of the api route for the DevApiController
func (c *DevApiController) Routes() Routes {
	return Routes{
		{
			"JobNodes",
			strings.ToUpper("Post"),
			"/{user}/nodes",
			c.JobNodes,
		},
		{
			"UpdateJobNodes",
			strings.ToUpper("Put"),
			"/{user}/nodes",
			c.UpdateJobNodes,
		},
	}
}

// JobNodes - Nodes information for the job
func (c *DevApiController) JobNodes(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user := params["user"]

	jobNodes := &JobNodes{}
	if err := json.NewDecoder(r.Body).Decode(&jobNodes); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	result, err := c.service.JobNodes(r.Context(), user, *jobNodes)
	// If an error occurred, encode the error with the status code
	if err != nil {
		EncodeJSONResponse(err.Error(), &result.Code, w)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// UpdateJobNodes - Update or add new nodes information for the job
func (c *DevApiController) UpdateJobNodes(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user := params["user"]

	jobNodes := &JobNodes{}
	if err := json.NewDecoder(r.Body).Decode(&jobNodes); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	result, err := c.service.UpdateJobNodes(r.Context(), user, *jobNodes)
	// If an error occurred, encode the error with the status code
	if err != nil {
		EncodeJSONResponse(err.Error(), &result.Code, w)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}