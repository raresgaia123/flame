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

package apiserver

import (
	"context"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"wwwin-github.cisco.com/eti/fledge/pkg/openapi"
	"wwwin-github.cisco.com/eti/fledge/pkg/restapi"
	"wwwin-github.cisco.com/eti/fledge/pkg/util"
)

// DesignSchemasApiService is a service that implents the logic for the DesignSchemasApiServicer
// This service should implement the business logic for every endpoint for the DesignSchemasApi API.
// Include any external packages or services that will be required by this service.
type DesignSchemasApiService struct {
}

// NewDesignSchemasApiService creates a default api service
func NewDesignSchemasApiService() openapi.DesignSchemasApiServicer {
	return &DesignSchemasApiService{}
}

// CreateDesignSchema - Update a design schema
func (s *DesignSchemasApiService) CreateDesignSchema(ctx context.Context, user string, designId string,
	designSchema openapi.DesignSchema) (openapi.ImplResponse, error) {
	//TODO input validation
	zap.S().Debugf("Create design schema request received for designId: %v", designId)

	//create controller request
	uriMap := map[string]string{
		"user":     user,
		"designId": designId,
	}
	url := restapi.CreateURL(Host, Port, restapi.CreateDesignSchemaEndPoint, uriMap)

	//send get request
	code, resp, err := restapi.HTTPPost(url, designSchema, "application/json")
	zap.S().Debugf("code: %d, resp: %s", code, string(resp))

	// response to the user
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), fmt.Errorf("error while updating/inserting design schema")
	}

	if err = restapi.CheckStatusCode(code); err != nil {
		return openapi.Response(code, nil), err
	}

	return openapi.Response(http.StatusOK, nil), err
}

// GetDesignSchema - Get a design schema owned by user
func (s *DesignSchemasApiService) GetDesignSchema(ctx context.Context, user string, designId string,
	version string) (openapi.ImplResponse, error) {
	//TODO input validation
	zap.S().Debugf("Get design schema details for user: %s | designId: %s | version: %s", user, designId, version)

	//create controller request
	uriMap := map[string]string{
		"user":     user,
		"designId": designId,
		"version":  version,
	}
	url := restapi.CreateURL(Host, Port, restapi.GetDesignSchemaEndPoint, uriMap)

	//send get request
	code, responseBody, err := restapi.HTTPGet(url)

	// response to the user
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), fmt.Errorf("get design schema details request failed")
	}

	if err = restapi.CheckStatusCode(code); err != nil {
		return openapi.Response(code, nil), err
	}

	resp := openapi.DesignSchema{}
	err = util.ByteToStruct(responseBody, &resp)

	return openapi.Response(http.StatusOK, resp), err
}

// GetDesignSchemas - Get all design schemas in a design
func (s *DesignSchemasApiService) GetDesignSchemas(ctx context.Context, user string, designId string) (openapi.ImplResponse, error) {
	//TODO input validation
	zap.S().Debugf("Get design schema details for user: %s | designId: %s", user, designId)

	//create controller request
	uriMap := map[string]string{
		"user":     user,
		"designId": designId,
	}
	url := restapi.CreateURL(Host, Port, restapi.GetDesignSchemasEndPoint, uriMap)

	//send get request
	code, responseBody, err := restapi.HTTPGet(url)

	// response to the user
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), fmt.Errorf("get design schema details request failed")
	}

	if err = restapi.CheckStatusCode(code); err != nil {
		return openapi.Response(code, nil), err
	}

	var resp []openapi.DesignSchema
	err = util.ByteToStruct(responseBody, &resp)

	return openapi.Response(http.StatusOK, resp), err
}

// UpdateDesignSchema - Update a schema for a given design
func (s *DesignSchemasApiService) UpdateDesignSchema(ctx context.Context, user string, designId string, version string,
	designSchema openapi.DesignSchema) (openapi.ImplResponse, error) {
	zap.S().Debugf("Update design schema request received for designId: %v", designId)

	//create controller request
	uriMap := map[string]string{
		"user":     user,
		"designId": designId,
		"version":  version,
	}
	url := restapi.CreateURL(Host, Port, restapi.UpdateDesignSchemaEndPoint, uriMap)

	//send put request
	code, resp, err := restapi.HTTPPut(url, designSchema, "application/json")
	zap.S().Debugf("code: %d, resp: %s", code, string(resp))

	// response to the user
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), fmt.Errorf("error while updating design schema")
	}

	if err = restapi.CheckStatusCode(code); err != nil {
		return openapi.Response(code, nil), err
	}

	return openapi.Response(http.StatusOK, nil), err
}