// © 2024 NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package transformers

import "fmt"

// ToODPS maps datasets to an ODPS v1.0 structure.
func ToODPS(datasets []Dataset) map[string]interface{} {
	var apiList []map[string]interface{}
	for _, ds := range datasets {
		apiList = append(apiList, map[string]interface{}{
			"id":          ds.ID,
			"title":       ds.Shortname,
			"description": fmt.Sprintf("Dataset type: %s", ds.Type),
			"version":     "v1",
			"contact": map[string]interface{}{
				"name":  "Support Open Data Hub",
				"email": ContactEmail,
			},
			"endpoints": []map[string]interface{}{
				{
					"url":           ds.ApiUrl,
					"methods":       []string{"GET"},
					"formats":       []string{"application/json"},
					"documentation": ContactWebsite,
				},
			},
			"license": map[string]interface{}{
				"name": "CC BY 4.0",
				"url":  "https://creativecommons.org/licenses/by/4.0/",
			},
		})
	}

	return map[string]interface{}{
		"odps": "1.0",
		"catalog": map[string]interface{}{
			"title":       "API Catalog",
			"description": "A catalog of APIs provided by " + OrganizationName + ".",
			"publisher": map[string]interface{}{
				"name": OrganizationName,
				"url":  OrganizationURL,
			},
			"apis": apiList,
		},
	}
}
