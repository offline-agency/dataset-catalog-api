// © 2024 NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package transformers

import (
	"fmt"
)

func ToODPS30(datasets []Dataset) map[string]interface{} {
	if len(datasets) == 0 {
		return nil
	}
	ds := datasets[0]

	productEn := map[string]interface{}{
		"name":              ds.Shortname,
		"productID":         ds.ID,
		"valueProposition":  fmt.Sprintf("A tailored data product for %s data", ds.Type),
		"description":       ds.ApiDescription["en"],
		"productSeries":     ds.Shortname + " Series",
		"visibility":        "public",
		"status":            "active",
		"version":           "v1.0",
		"categories":        ds.Category,
		"standards":         []string{"Standard-Dev"},
		"tags":              ds.ODHTags,
		"brandSlogan":       BrandSlogan,
		"type":              ds.Type,
		"logoURL":           ds.Self,
		"OutputFileFormats": []string{"JSON", "YAML"},
		"useCases": []map[string]interface{}{
      {
        "useCase": map[string]interface{}{
          "useCaseTitle":       "Discover Insights - example",
          "useCaseDescription": "description example",
          "useCaseURL":         ds.ApiUrl + "/usecase/insights",
        },
      },
    },
	}

	product := map[string]interface{}{
		"en": productEn,
	}

	recommendedDataProducts := []string{
		ds.ApiUrl + "/recommended/1",
		ds.ApiUrl + "/recommended/2",
	}

	pricingPlansEn := []interface{}{
		map[string]interface{}{
			"name":                   "Free",
			"priceCurrency":          "EUR",
			"price":                  "0",
			"billingDuration":        "Monthly",
			"unit":                   "month",
			"maxTransactionQuantity": "1000",
			"offering":               []string{"Basic"},
		},
	}
	pricingPlans := map[string]interface{}{
		"en": pricingPlansEn,
	}

	dataOps := map[string]interface{}{
		"data": map[string]interface{}{
			"schemaLocationURL": ds.ApiUrl + "/schema",
		},
		"lineage": map[string]interface{}{
			"dataLineageTool":   "LineageTool",
			"dataLineageOutput": "LineageInfo",
		},
		"infrastructure": map[string]interface{}{
			"containerTool":    "Docker",
      "platform":         "Kubernetes",
      "region":           "eu-south-1",
      "storageTechnology": "S3",
      "storageType":      "Object",
		},
		"build": map[string]interface{}{
			"format":                     "docker",
			"hashType":                   "SHA256",
			"checksum":                   "abc123",
			"signatureType":              "PGP",
			"scriptURL":                  ds.ApiUrl + "/build.sh",
			"deploymentDocumentationURL": ds.ApiUrl + "/deploy",
		},
	}

	dataAccess := map[string]interface{}{
		"type":                 "REST",
		"authenticationMethod": "None",
		"specification":        "OpenAPI",
		"format":               "JSON",
		"documentationURL":     ds.ApiUrl + "/docs",
	}

	// Build SLA.
	SLA := []interface{}{
		map[string]interface{}{
			"dimension": "Availability",
			"displaytitle": []interface{}{
				map[string]interface{}{"en": "Availability"},
			},
			"objective": 99.9,
			"unit":      "%",
			"monitoring": map[string]interface{}{
				"type":      "Service Level",
				"reference": ds.ApiUrl + "/monitoring",
				"spec":      "SLA Spec",
			},
		},
	}

	support := map[string]interface{}{
		"phoneNumber":       ContactPhoneNumber,
		"phoneServiceHours": "9-5",
		"email":             ContactEmail,
		"emailServiceHours": "9-5",
		"documentationURL":  ds.SwaggerUrl,
	}

	dataQuality := []interface{}{
		map[string]interface{}{
			"dimension": "Accuracy",
			"displaytitle": []interface{}{
				map[string]interface{}{"en": "Accuracy"},
			},
			"objective": 95.0,
			"unit":      "%",
			"monitoring": map[string]interface{}{
				"type":      "Quality",
				"reference": ds.ApiUrl + "/quality",
				"spec":      "Quality Spec",
			},
		},
	}

	license := map[string]interface{}{
		"scope": map[string]interface{}{
			"definition":      "Full access",
			"language":        "en",
			"restrictions":    "None",
			"geographicalArea": []string{"Global"},
			"permanent":       true,
			"exclusive":       false,
			"rights":          []string{"Read", "Write"},
		},
		"termination": map[string]interface{}{
			"terminationConditions": "Violation of terms",
			"continuityConditions":  "N/A",
		},
		"governance": map[string]interface{}{
			"ownership":       OrganizationName,
			"damages":         "None",
			"confidentiality": "High",
			"applicableLaws":  "GDPR",
			"warranties":      "None",
			"audit":           "Annual",
			"forceMajeure":    "Standard",
		},
	}

	dataHolder := map[string]interface{}{
		"taxID":            TaxID,
		"vatID":            VatID,
		"businessDomain":   "Data",
		"logoURL":          ds.Self,
		"description":      BrandSlogan,
		"URL":              ds.Self,
		"telephone":        ContactPhoneNumber,
		"streetAddress":    StreetAddress,
		"postalCode":       PostalCode,
		"addressRegion":    AddressRegion,
		"addressLocality":  AddressLocality,
		"addressCountry":   "IT",
		"aggregateRating":  "5 stars",
		"ratingCount":      100,
		"slogan":           BrandSlogan,
		"parentOrganization": OrganizationName,
	}

	return map[string]interface{}{
		"schema":                   "https://opendataproducts.org/v3.0/schema/odps.yaml",
		"version":                  "dev",
		"product":                  product,
		"recommendedDataProducts":  recommendedDataProducts,
		"pricingPlans":             pricingPlans,
		"dataOps":                  dataOps,
		"dataAccess":               dataAccess,
		"SLA":                      SLA,
		"support":                  support,
		"dataQuality":              dataQuality,
		"license":                  license,
		"dataHolder":               dataHolder,
	}
}
