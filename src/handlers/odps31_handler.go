// © 2024 NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package handlers

import (
	"math"
	"net/http"
	"strconv"

	"log"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"opendatahub.com/dataset-catalog-api/transformers"
)

// ODPS31GinHandler handles the listing endpoint for ODPS31.
// GET /odps31?page={n} returns a paginated list (10 items per page) of dataset endpoints.
// Each endpoint object includes: uuid, datasetName, originalUrl, and url.
// Default output is YAML; use ?format=json for JSON.
func ODPS31GinHandler(c *gin.Context) {
	pageStr := c.Query("page")
	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	resp, err := fetchDatasetsResponse(page)
	if err != nil || resp == nil || len(resp.Items) == 0 {
		c.String(http.StatusNotFound, "No data found")
		return
	}

	totalItems := resp.TotalResults
	totalPages := int(math.Ceil(float64(totalItems) / float64(pageSize)))

	var endpoints []map[string]interface{}
	for _, ds := range resp.Items {
		item := map[string]interface{}{
			"uuid":        ds.ID,
			"datasetName": ds.Shortname,
			"originalUrl": ds.ApiUrl,
			"url":         transformers.BaseURL + "odps31/" + ds.ID,
		}
		endpoints = append(endpoints, item)
	}

	output := map[string]interface{}{
		"current_page": resp.CurrentPage,
		"total_pages":  totalPages,
		"endpoints":    endpoints,
	}

	format := c.Query("format")
	if format == "json" {
		c.JSON(http.StatusOK, output)
	} else {
		yamlData, err := yaml.Marshal(output)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error marshaling YAML")
			return
		}
		c.Data(http.StatusOK, "text/plain; charset=utf-8", yamlData)
	}
}

// ODPS31DetailGinHandler handles the detail endpoint for ODPS31.
// GET /odps31/:uuid returns detailed information for the dataset with the given UUID.
// Default output is YAML; use ?format=json for JSON.
func ODPS31DetailGinHandler(c *gin.Context) {
	datasetID := c.Param("uuid")
	if datasetID == "" {
		c.String(http.StatusBadRequest, "Missing dataset ID")
		return
	}
	log.Printf("ODPS31 detail endpoint requested for dataset ID: %s", datasetID)
	found := searchDatasetByID(datasetID)
	if found == nil {
		c.String(http.StatusNotFound, "Dataset not found")
		return
	}
	conv := ConvertDatasets([]transformers.Dataset{*found})
	output := transformers.ToODPS31(conv)
	format := c.Query("format")
	if format == "json" {
		c.JSON(http.StatusOK, output)
	} else {
		yamlData, err := yaml.Marshal(output)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error marshaling YAML")
			return
		}
		c.Data(http.StatusOK, "text/plain; charset=utf-8", yamlData)
	}
}
