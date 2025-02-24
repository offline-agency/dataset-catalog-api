// © 2024 NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package transformers

import (
	"fmt"
)

func ToODPS31(datasets []Dataset) map[string]interface{} {
	if len(datasets) == 0 {
		return nil
	}
	ds := datasets[0]

	en := map[string]interface{}{
		"OutputFileFormats": []string{"JSON", "YAML"},
		"brandSlogan":       BrandSlogan,
		"categories":        ds.Category,
		"description":       ds.ApiDescription["en"],
		"logoURL":           ds.Self,
		"name":              ds.Shortname,
		"productID":         ds.ID,
		"productSeries":     ds.Shortname + " Series",
		"standards":         []string{"Standard-Dev"},
		"status":            "active",
		"tags":              []string{},
		"type":              ds.Type,
		"useCases": []map[string]interface{}{
			{
				"useCase": map[string]interface{}{
					"useCaseTitle":       "Discover Insights - example",
					"useCaseDescription": "description example",
					"useCaseURL":         ds.ApiUrl + "/usecase/insights",
				},
			},
		},
		"valueProposition": fmt.Sprintf("A tailored data product for %s data", ds.Type),
		"version":          "v1.0",
		"visibility":       "public",
	}

	dataAccess := map[string]interface{}{
		"authenticationMethod": "None",
		"documentationURL":     ds.SwaggerUrl,
		"format":               "JSON",
		"specification":        "OpenAPI",
		"type":                 "REST",
	}

	dataHolder := map[string]interface{}{
		"URL":              ds.Self,
		"addressCountry":   "IT",
		"addressLocality":  AddressLocality,
		"addressRegion":    AddressRegion,
		"aggregateRating":  "5 stars",
		"businessDomain":   "Data",
		"description":      BrandSlogan,
		"logoURL":          ds.Self,
		"parentOrganization": OrganizationName,
		"postalCode":       PostalCode,
		"ratingCount":      100,
		"slogan":           BrandSlogan,
		"streetAddress":    StreetAddress,
		"taxID":            TaxID,
		"telephone":        ContactPhoneNumber,
		"vatID":            VatID,
	}

	dataOps := map[string]interface{}{
		"build": map[string]interface{}{
			"checksum":                   "abc123",
			"deploymentDocumentationURL": ds.Self + "/deploy",
			"format":                     "docker",
			"hashType":                   "SHA256",
			"scriptURL":                  ds.Self + "/build.sh",
			"signatureType":              "PGP",
		},
		"data": map[string]interface{}{
			"schemaLocationURL": ds.Self + "/schema",
		},
		"infrastructure": map[string]interface{}{
			"containerTool":    "Docker",
			"platform":         "Kubernetes",
			"region":           "eu-south-1",
			"storageTechnology": "S3",
			"storageType":      "Object",
		},
		"lineage": map[string]interface{}{
			"dataLineageOutput": "LineageInfo",
			"dataLineageTool":   "LineageTool",
		},
	}

	SLA := []interface{}{
		map[string]interface{}{
			"dimension":    "Availability",
			"displaytitle": []interface{}{map[string]interface{}{"en": "Availability"}},
			"monitoring": map[string]interface{}{
				"reference": ds.Self + "/monitoring",
				"spec":      "SLA Spec",
				"type":      "Service Level",
			},
			"objective": 99.9,
			"unit":      "%",
		},
	}

	dataQuality := []interface{}{
		map[string]interface{}{
			"dimension":    "Accuracy",
			"displaytitle": []interface{}{map[string]interface{}{"en": "Accuracy"}},
			"monitoring": map[string]interface{}{
				"reference": ds.Self + "/quality",
				"spec":      "Quality Spec",
				"type":      "Quality",
			},
			"objective": 95.0,
			"unit":      "%",
		},
	}

	pricingPlans := map[string]interface{}{
		"en": []interface{}{
			map[string]interface{}{
				"billingDuration":        "Monthly",
				"maxTransactionQuantity": "1000",
				"name":                   "Free",
				"offering":               []string{"Basic"},
				"price":                  "0",
				"priceCurrency":          "EUR",
				"unit":                   "month",
			},
		},
	}

	support := map[string]interface{}{
		"documentationURL": ds.SwaggerUrl,
		"email":           ContactEmail,
		"emailServiceHours": "9-5",
		"phoneNumber":      ContactPhoneNumber,
		"phoneServiceHours": "9-5",
	}

	product := map[string]interface{}{
		"SLA":         SLA,
		"dataAccess":  dataAccess,
		"dataHolder":  dataHolder,
		"dataOps":     dataOps,
		"dataQuality": dataQuality,
		"en":          en,
		"license": map[string]interface{}{
			"governance": map[string]interface{}{
				"applicableLaws": "GDPR",
				"audit":          "Annual",
				"confidentiality": "High",
				"damages":         "None",
				"forceMajeure":    "Standard",
				"ownership":       OrganizationName,
				"warranties":      "None",
			},
			"scope": map[string]interface{}{
				"definition":      "Full access",
				"exclusive":       false,
				"geographicalArea": []string{"Global"},
				"language":        "en",
				"permanent":       true,
				"restrictions":    "None",
				"rights":          []string{"Read", "Write"},
			},
			"termination": map[string]interface{}{
				"continuityConditions":  "N/A",
				"terminationConditions": "Violation of terms",
			},
		},
		"pricingPlans":            pricingPlans,
		"recommendedDataProducts": []string{ds.Self + "/recommended/1", ds.Self + "/recommended/2"},
		"support":                 support,
	}

	details := map[string]interface{}{
		"summary":     ds.Shortname,
		"description": ds.ApiDescription["en"],
		"language":    "en",
		"metadata":    ds.Meta,
	}

	return map[string]interface{}{
		"schema":  "https://opendataproducts.org/v3.1/schema/odps.yaml",
		"version": "3.1",
		"product": product,
		"details": details,
		"dct:issued":   ds.FirstImport,
    "dct:modified": ds.LastChange,
	}
}
