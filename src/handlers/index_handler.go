// © 2024 NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// IndexHandler renders an index HTML page listing all available endpoints.
func IndexHandler(c *gin.Context) {
	// Define a list of endpoint paths.
	endpoints := []string{
		"/dcat",
		"/odps",
		"/odps30",
		"/odps31",
	}
	c.HTML(http.StatusOK, "index.html", gin.H{
		"endpoints": endpoints,
	})
}
