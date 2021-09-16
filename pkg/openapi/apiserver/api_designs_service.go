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
	"strconv"

	"go.uber.org/zap"

	"wwwin-github.cisco.com/eti/fledge/pkg/openapi"
	"wwwin-github.cisco.com/eti/fledge/pkg/restapi"
	"wwwin-github.cisco.com/eti/fledge/pkg/util"
)

// DesignsApiService is a service that implents the logic for the DesignsApiServicer
// This service should implement the business logic for every endpoint for the DesignsApi API.
// Include any external packages or services that will be required by this service.
type DesignsApiService struct {
}

// NewDesignsApiService creates a default api service
func NewDesignsApiService() openapi.DesignsApiServicer {
	return &DesignsApiService{}
}

// CreateDesign - Create a new design template.
func (s *DesignsApiService) CreateDesign(ctx context.Context, user string, designInfo openapi.DesignInfo) (openapi.ImplResponse, error) {
	//TODO input validation
	zap.S().Debugf("New design request received for user: %s | designInfo: %v", user, designInfo)

	// create controller request
	uriMap := map[string]string{
		"user": user,
	}
	url := restapi.CreateURL(Host, Port, restapi.CreateDesignEndPoint, uriMap)

	// send post request
	code, _, err := restapi.HTTPPost(url, designInfo, "application/json")

	// response to the user
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), fmt.Errorf("create new design request failed")
	}

	if err = restapi.CheckStatusCode(code); err != nil {
		return openapi.Response(code, nil), err
	}

	return openapi.Response(http.StatusCreated, nil), nil
}

// GetDesign - Get design template information
func (s *DesignsApiService) GetDesign(ctx context.Context, user string, designId string) (openapi.ImplResponse, error) {
	//TODO input validation
	zap.S().Debugf("Get design template information for user: %s | designId: %s", user, designId)

	//create controller request
	uriMap := map[string]string{
		"user":     user,
		"designId": designId,
	}
	url := restapi.CreateURL(Host, Port, restapi.GetDesignEndPoint, uriMap)

	//send get request
	code, responseBody, err := restapi.HTTPGet(url)

	//response to the user
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), fmt.Errorf("get design template information request failed")
	}

	if err = restapi.CheckStatusCode(code); err != nil {
		return openapi.Response(code, nil), err
	}

	resp := openapi.Design{}
	err = util.ByteToStruct(responseBody, &resp)

	return openapi.Response(http.StatusOK, resp), err
}

// GetDesigns - Get list of all the designs created by the user.
func (s *DesignsApiService) GetDesigns(ctx context.Context, user string, limit int32) (openapi.ImplResponse, error) {
	//TODO input validation
	zap.S().Debugf("get list of designs for user: %s | limit: %d", user, limit)

	//create controller request
	//construct URL
	uriMap := map[string]string{
		"user":  user,
		"limit": strconv.Itoa(int(limit)),
	}
	url := restapi.CreateURL(Host, Port, restapi.GetDesignsEndPoint, uriMap)

	//send get request
	code, responseBody, err := restapi.HTTPGet(url)

	//response to the user
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), fmt.Errorf("get design template information request failed")
	}

	if err = restapi.CheckStatusCode(code); err != nil {
		return openapi.Response(code, nil), err
	}

	var resp []openapi.DesignInfo
	err = util.ByteToStruct(responseBody, &resp)
	return openapi.Response(http.StatusOK, resp), err
}