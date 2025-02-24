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

// ODPS30GinHandler handles the listing endpoint for ODPS30.
// GET /odps30?page={n} returns a paginated list (10 items per page) of dataset endpoints.
// Each endpoint object includes: uuid, datasetName, originalUrl, and url.
// Default output is YAML; use ?format=json for JSON.
func ODPS30GinHandler(c *gin.Context) {
	// Ensure that pagination always starts at 1.
	pageStr := c.Query("page")
	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		} else {
			c.String(http.StatusNotFound, "No data found")
			return
		}
	}

	// Fetch the datasets for the requested page.
	resp, err := fetchDatasetsResponse(page)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error fetching data")
		return
	}

	// Calculate total pages from total records.
	totalItems := resp.TotalResults
	totalPages := int(math.Ceil(float64(totalItems) / float64(pageSize)))

	// If the requested page is greater than totalPages, return no data.
	if page > totalPages {
		c.String(http.StatusNotFound, "No data found")
		return
	}

	// Build an array of objects with uuid, datasetName, originalUrl, and internal URL.
	var endpoints []map[string]interface{}
	for _, ds := range resp.Items {
		item := map[string]interface{}{
			"uuid":        ds.ID,
			"datasetName": ds.Shortname,
			"originalUrl": ds.ApiUrl, // Assuming ApiUrl contains the external API URL
			"url":         transformers.BaseURL + "odps30/" + ds.ID,
		}
		endpoints = append(endpoints, item)
	}

	output := map[string]interface{}{
		"current_page":  resp.CurrentPage,
		"total_pages":   totalPages,
		"totalRecord":   totalItems,
		"endpoints":     endpoints,
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



// ODPS30DetailGinHandler handles the detail endpoint for ODPS30.
// GET /odps30/:uuid returns detailed information for the dataset with the given UUID.
// Default output is YAML; use ?format=json for JSON.
func ODPS30DetailGinHandler(c *gin.Context) {
	datasetID := c.Param("uuid")
	if datasetID == "" {
		c.String(http.StatusBadRequest, "Missing dataset ID")
		return
	}
	log.Printf("ODPS30 detail endpoint requested for dataset ID: %s", datasetID)
	found := searchDatasetByID(datasetID)
	if found == nil {
		c.String(http.StatusNotFound, "Dataset not found")
		return
	}
	conv := ConvertDatasets([]transformers.Dataset{*found})
	output := transformers.ToODPS30(conv)
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
