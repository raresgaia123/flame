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

// A DesignsApiController binds http requests to an api service and writes the service results to the http response
type DesignsApiController struct {
	service DesignsApiServicer
}

// NewDesignsApiController creates a default api controller
func NewDesignsApiController(s DesignsApiServicer) Router {
	return &DesignsApiController{service: s}
}

// Routes returns all of the api route for the DesignsApiController
func (c *DesignsApiController) Routes() Routes {
	return Routes{
		{
			"CreateDesign",
			strings.ToUpper("Post"),
			"/{user}/designs",
			c.CreateDesign,
		},
		{
			"GetDesign",
			strings.ToUpper("Get"),
			"/{user}/designs/{designId}",
			c.GetDesign,
		},
		{
			"GetDesigns",
			strings.ToUpper("Get"),
			"/{user}/designs",
			c.GetDesigns,
		},
	}
}

// CreateDesign - Create a new design template.
func (c *DesignsApiController) CreateDesign(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user := params["user"]

	designInfo := &DesignInfo{}
	if err := json.NewDecoder(r.Body).Decode(&designInfo); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	result, err := c.service.CreateDesign(r.Context(), user, *designInfo)
	// If an error occurred, encode the error with the status code
	if err != nil {
		EncodeJSONResponse(err.Error(), &result.Code, w)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// GetDesign - Get design template information
func (c *DesignsApiController) GetDesign(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user := params["user"]

	designId := params["designId"]

	result, err := c.service.GetDesign(r.Context(), user, designId)
	// If an error occurred, encode the error with the status code
	if err != nil {
		EncodeJSONResponse(err.Error(), &result.Code, w)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// GetDesigns - Get list of all the designs created by the user.
func (c *DesignsApiController) GetDesigns(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	query := r.URL.Query()
	user := params["user"]

	limit, err := parseInt32Parameter(query.Get("limit"), false)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	result, err := c.service.GetDesigns(r.Context(), user, limit)
	// If an error occurred, encode the error with the status code
	if err != nil {
		EncodeJSONResponse(err.Error(), &result.Code, w)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}
