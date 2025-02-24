// © 2024 NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"opendatahub.com/dataset-catalog-api/transformers"
)

func ODPSGinHandler(c *gin.Context) {
	ds, err := fetchDatasets(1)
	if err != nil || len(ds) == 0 {
		c.String(http.StatusNotFound, "No data found")
		return
	}
	output := transformers.ToODPS(ConvertDatasets(ds))
	c.JSON(http.StatusOK, output)
}
