// © 2024 NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package transformers

import (
	"fmt"
	"time"
)

// ToDCAT maps a slice of datasets to a DCAT‑AP 3.0 compliant catalog.
// It uses qualified properties (e.g., dct:title, dct:description, dct:type),
// language‑tagged values, and adds mandatory metadata (such as dct:identifier, dct:issued, and dct:modified).
func ToDCAT(datasets []Dataset) map[string]interface{} {
	now := time.Now().Format("2006-01-02")
	var datasetList []map[string]interface{}
	for _, ds := range datasets {
		datasetList = append(datasetList, map[string]interface{}{
			"@type":          "dcat:Dataset",
			"@id":            ds.Self,
			"dct:identifier": ds.ID,
			// Mandatory property: dct:type
			"dct:type": map[string]string{
				"en": "dcat:Dataset",
			},
			"dct:title": map[string]string{
				"en": ds.Shortname,
			},
			"dct:description": map[string]string{
				"en": fmt.Sprintf("Dataset type: %s", ds.Type),
			},
			"dct:issued":   ds.FirstImport,
			"dct:modified": ds.LastChange,
			"distribution": []map[string]interface{}{
				{
					"@type":          "dcat:Distribution",
					"@id":            ds.ApiUrl,
					"dct:identifier": ds.ApiUrl, // using API URL as identifier
					"dct:type": map[string]string{
						"en": "dcat:Distribution",
					},
					"dct:title": map[string]string{
						"en": ds.Shortname + " API Endpoint",
					},
					"dct:format": "application/json",
					"accessURL":  ds.ApiUrl,
				},
			},
		})
	}

	return map[string]interface{}{
		"@context": map[string]interface{}{
			"dcat": "https://www.w3.org/ns/dcat#",
			"dct":  "http://purl.org/dc/terms/",
			"foaf": "http://xmlns.com/foaf/0.1/",
			"xsd":  "http://www.w3.org/2001/XMLSchema#",
			// Define language containers:
			"dct:title": map[string]interface{}{
				"@id":       "dct:title",
				"@container": "@language",
			},
			"dct:description": map[string]interface{}{
				"@id":       "dct:description",
				"@container": "@language",
			},
			"dct:issued": map[string]interface{}{
				"@id":   "dct:issued",
				"@type": "xsd:date",
			},
			"dct:modified": map[string]interface{}{
				"@id":   "dct:modified",
				"@type": "xsd:date",
			},
		},
		"@type": "dcat:Catalog",
		"@id":   BaseURL + "api-catalog",
		// Mandatory property for catalog:
		"dct:type": map[string]string{
			"en": "dcat:Catalog",
		},
		"dct:identifier": "catalog-001",
		"dct:title": map[string]string{
			"en": OrganizationName + " API Catalog",
		},
		"dct:description": map[string]string{
			"en": "A catalog of APIs provided by " + OrganizationName + ".",
		},
		"dct:issued":   now,
		"dct:modified": now,
		"publisher": map[string]interface{}{
			"@type":          "foaf:Organization",
			"dct:identifier": "org-001",
			"dct:title": map[string]string{
				"en": OrganizationName,
			},
			"homepage": OrganizationURL,
		},
		"dataset": datasetList,
	}
}
