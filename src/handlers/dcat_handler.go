// © 2024 NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"opendatahub.com/dataset-catalog-api/transformers"
)

func DcatGinHandler(c *gin.Context) {
	ds, err := fetchDatasets(1)
	if err != nil || len(ds) == 0 {
		c.String(http.StatusNotFound, "No data found")
		return
	}
	output := transformers.ToDCAT(ConvertDatasets(ds))
	format := c.Query("format")
	if format == "yaml" {
		yamlData, err := yaml.Marshal(output)
    if err != nil {
      c.String(http.StatusInternalServerError, "Error marshaling YAML")
      return
    }
    c.Data(http.StatusOK, "text/plain; charset=utf-8", yamlData)
	} else {
	  c.JSON(http.StatusOK, output)
	}
}
