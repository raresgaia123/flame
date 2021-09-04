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
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// A DesignCodesApiController binds http requests to an api service and writes the service results to the http response
type DesignCodesApiController struct {
	service DesignCodesApiServicer
}

// NewDesignCodesApiController creates a default api controller
func NewDesignCodesApiController(s DesignCodesApiServicer) Router {
	return &DesignCodesApiController{service: s}
}

// Routes returns all of the api route for the DesignCodesApiController
func (c *DesignCodesApiController) Routes() Routes {
	return Routes{
		{
			"CreateDesignCode",
			strings.ToUpper("Post"),
			"/{user}/designs/{designId}/codes",
			c.CreateDesignCode,
		},
		{
			"GetDesignCode",
			strings.ToUpper("Get"),
			"/{user}/designs/{designId}/codes/{version}",
			c.GetDesignCode,
		},
		{
			"UpdateDesignCode",
			strings.ToUpper("Put"),
			"/{user}/designs/{designId}/codes/{version}",
			c.UpdateDesignCode,
		},
	}
}

// CreateDesignCode - Upload a new design code
func (c *DesignCodesApiController) CreateDesignCode(w http.ResponseWriter, r *http.Request) {
	var maxMemory int64 = 32 << 20 // 32MB
	if err := r.ParseMultipartForm(maxMemory); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	params := mux.Vars(r)
	user := params["user"]

	designId := params["designId"]

	fileName := r.FormValue("fileName")
	fileVer := r.FormValue("fileVer")

	fileData, err := ReadFormFileToTempFile(r, "fileData")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	result, err := c.service.CreateDesignCode(r.Context(), user, designId, fileName, fileVer, fileData)
	// If an error occurred, encode the error with the status code
	if err != nil {
		EncodeJSONResponse(err.Error(), &result.Code, w)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// GetDesignCode - Get a zipped design code file owned by user
func (c *DesignCodesApiController) GetDesignCode(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user := params["user"]

	designId := params["designId"]

	version := params["version"]

	result, err := c.service.GetDesignCode(r.Context(), user, designId, version)
	// If an error occurred, encode the error with the status code
	if err != nil {
		EncodeJSONResponse(err.Error(), &result.Code, w)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// UpdateDesignCode - Update a design code
func (c *DesignCodesApiController) UpdateDesignCode(w http.ResponseWriter, r *http.Request) {
	var maxMemory int64 = 32 << 20 // 32MB
	if err := r.ParseMultipartForm(maxMemory); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	params := mux.Vars(r)
	user := params["user"]

	designId := params["designId"]

	version := params["version"]

	fileName := r.FormValue("fileName")
	fileVer := r.FormValue("fileVer")

	fileData, err := ReadFormFileToTempFile(r, "fileData")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	result, err := c.service.UpdateDesignCode(r.Context(), user, designId, version, fileName, fileVer, fileData)
	// If an error occurred, encode the error with the status code
	if err != nil {
		EncodeJSONResponse(err.Error(), &result.Code, w)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}
