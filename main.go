package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

// Dataset represents a minimal structure from the external API.
type Dataset struct {
	ID        string `json:"Id"`
	Self      string `json:"Self"`
	Type      string `json:"Type"`
	Shortname string `json:"Shortname"`
	ApiUrl    string `json:"ApiUrl"`
}

// ODHResponse wraps the external API's response.
type ODHResponse struct {
	TotalResults int       `json:"TotalResults"`
	TotalPages   int       `json:"TotalPages"`
	CurrentPage  int       `json:"CurrentPage"`
	NextPage     string    `json:"NextPage"`
	Items        []Dataset `json:"Items"`
}

const pageSize = 10

// cacheItem holds cached data with an expiration time.
type cacheItem struct {
	data       []Dataset
	expiration time.Time
}

// datasetCache holds cached responses keyed by page number.
var (
	datasetCache = make(map[int]cacheItem)
	cacheMutex   sync.RWMutex
)

// fetchDatasets retrieves data from the external API and caches the result.
// If cached data is still valid (not expired), it returns the cached data.
func fetchDatasets(page int) ([]Dataset, error) {
	// Check if data for this page is in cache.
	cacheMutex.RLock()
	if item, found := datasetCache[page]; found {
		if time.Now().Before(item.expiration) {
			cacheMutex.RUnlock()
			return item.data, nil
		}
	}
	cacheMutex.RUnlock()

	// Data not in cache or expired; fetch from external API.
	url := fmt.Sprintf("https://tourism.api.opendatahub.com/v1/MetaData?pagenumber=%d&limit=%d", page, pageSize)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data ODHResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Println("Error decoding JSON:", err)
		return nil, err
	}

	if data.Items == nil || len(data.Items) == 0 {
		log.Println("No datasets found from API")
		return nil, nil
	}

	// Update the cache with the fetched data (using a TTL of 5 minutes).
	cacheMutex.Lock()
	datasetCache[page] = cacheItem{
		data:       data.Items,
		expiration: time.Now().Add(5 * time.Minute),
	}
	cacheMutex.Unlock()

	return data.Items, nil
}

// getDefaultDatasets returns default dataset(s) in case of errors.
func getDefaultDatasets() []Dataset {
	return []Dataset{
		{
			ID:        "default",
			Self:      "https://example.org/default",
			Type:      "unknown",
			Shortname: "Default Dataset",
			ApiUrl:    "https://example.org/api/default",
		},
	}
}

// toDCAT maps datasets to a DCAT-compliant structure.
func toDCAT(datasets []Dataset) map[string]interface{} {
	var datasetList []map[string]interface{}
	for _, ds := range datasets {
		datasetList = append(datasetList, map[string]interface{}{
			"@type":       "dcat:Dataset",
			"@id":         ds.Self,
			"title":       ds.Shortname,
			"description": fmt.Sprintf("Dataset type: %s", ds.Type),
			"distribution": []map[string]interface{}{
				{
					"@type":     "dcat:Distribution",
					"@id":       ds.ApiUrl,
					"title":     fmt.Sprintf("%s API Endpoint", ds.Shortname),
					"format":    "application/json",
					"accessURL": ds.ApiUrl,
				},
			},
		})
	}

	return map[string]interface{}{
		"@context":    "https://www.w3.org/ns/dcat",
		"@type":       "dcat:Catalog",
		"@id":         "https://example.org/api-catalog",
		"title":       "NOI SPA API Catalog",
		"description": "A catalog of APIs provided by NOI SPA.",
		"publisher": map[string]interface{}{
			"@type":    "foaf:Organization",
			"name":     "Noi SPA",
			"homepage": "https://noi.bz.it",
		},
		"dataset": datasetList,
	}
}

// toODPS maps datasets to an ODPS v1.0 structure.
func toODPS(datasets []Dataset) map[string]interface{} {
	var apiList []map[string]interface{}
	for _, ds := range datasets {
		apiList = append(apiList, map[string]interface{}{
			"id":          ds.ID,
			"title":       ds.Shortname,
			"description": fmt.Sprintf("Dataset type: %s", ds.Type),
			"version":     "v1", // Internal API version is set to v1
			"contact": map[string]interface{}{
				"name":  "Support Open Data Hub",
				"email": "help@opendatahub.com",
			},
			"endpoints": []map[string]interface{}{
				{
					"url":           ds.ApiUrl,
					"methods":       []string{"GET"},
					"formats":       []string{"application/json"},
					"documentation": "https://opendatahub.com/",
				},
			},
			"license": map[string]interface{}{
				"name": "CC BY 4.0",
				"url":  "https://creativecommons.org/licenses/by/4.0/",
			},
		})
	}

	return map[string]interface{}{
		"odps": "1.0", // This indicates that the catalog follows the ODPS 1.0 convention.
		"catalog": map[string]interface{}{
			"title":       "API Catalog",
			"description": "A catalog of APIs provided by NOI SPA.",
			"publisher": map[string]interface{}{
				"name": "Noi SPA",
				"url":  "https://noi.bz.it",
			},
			"apis": apiList,
		},
	}
}

// toODPS31 maps datasets to an ODPS 3.1 structure.
func toODPS31(datasets []Dataset) map[string]interface{} {
	// Build useCases and unique categories.
	var useCases []map[string]interface{}
	categoriesSet := make(map[string]bool)
	for _, ds := range datasets {
		useCases = append(useCases, map[string]interface{}{
			"useCase": map[string]interface{}{
				"useCaseTitle":       ds.Shortname,
				"useCaseDescription": fmt.Sprintf("Dataset type: %s", ds.Type),
				"useCaseURL":         ds.ApiUrl,
			},
		})
		if ds.Type != "" {
			categoriesSet[ds.Type] = true
		}
	}
	var categories []string
	for cat := range categoriesSet {
		categories = append(categories, cat)
	}

	product := map[string]interface{}{
		"SLA": []map[string]interface{}{
			{
				"dimension": "Availability",
				"displaytitle": []map[string]interface{}{
					{"en": "Availability SLA"},
				},
				"monitoring": map[string]interface{}{
					"reference": "https://example.org/monitoring",
					"spec":      "SLA Monitoring Spec",
					"type":      "Service Level",
				},
				"objective": 99.9,
				"unit":      "%",
			},
		},
		"dataAccess": map[string]interface{}{
			"authenticationMethod": "None",
			"documentationURL":     "https://example.org/docs",
			"format":               "JSON",
			"specification":        "OpenAPI",
			"type":                 "REST",
		},
		"dataHolder": map[string]interface{}{
			"URL":                "https://noi.bz.it",
			"addressCountry":     "IT",
			"addressLocality":    "Bolzano",
			"addressRegion":      "Alto Adige",
			"aggregateRating":    "4.5",
			"businessDomain":     "Tourism",
			"description":        "Data holder description",
			"logoURL":            "https://opendatahub.com/img/NOI_OPENDATAHUB_NEW_BK_nospace-01.svg",
			"parentOrganization": "Noi SPA",
			"postalCode":         "39100",
			"ratingCount":        100,
			"slogan":             "Open Data for All",
			"streetAddress":      "Via Volta 13/A",
			"taxID":              "02595720216",
			"telephone":          "+390471066600",
			"vatID":              "02595720216",
		},
		"dataOps": map[string]interface{}{
			"build": map[string]interface{}{
				"checksum":                   "abc123",
				"deploymentDocumentationURL": "https://opendatahub.com/deploy",
				"format":                     "docker",
				"hashType":                   "sha256",
				"scriptURL":                  "https://opendatahub.com/build.sh",
				"signatureType":              "PGP",
			},
			"data": map[string]interface{}{
				"schemaLocationURL": "https://opendatahub.com/schema",
			},
			"infrastructure": map[string]interface{}{
				"containerTool":    "Docker",
				"platform":         "Kubernetes",
				"region":           "eu-south-1",
				"storageTechnology": "S3",
				"storageType":      "object",
			},
			"lineage": map[string]interface{}{
				"dataLineageOutput": "Lineage info",
				"dataLineageTool":   "Lineage tool",
			},
		},
		"dataQuality": []map[string]interface{}{
			{
				"dimension": "Accuracy",
				"displaytitle": []map[string]interface{}{
					{"en": "Data Accuracy"},
				},
				"monitoring": map[string]interface{}{
					"reference": "https://opendatahub.com/quality",
					"spec":      "Quality Spec",
					"type":      "Quality",
				},
				"objective": 95.0,
				"unit":      "%",
			},
		},
		"en": map[string]interface{}{
			"OutputFileFormats": []string{"JSON", "YAML"},
			"brandSlogan":       "Open Data for All",
			"categories":        categories,
			"description":       "A collection of APIs/datasets from Open Data Hub.",
			"logoURL":           "https://opendatahub.com/img/NOI_OPENDATAHUB_NEW_BK_nospace-01.svg",
			"name":              "API Catalog",
			"productID":         "catalog-001",
			"productSeries":     "v1.0",
			"standards":         []string{"ODPS 3.1"},
			"status":            "active",
			"tags":              []string{"API", "Dataset"},
			"type":              "data",
			"useCases":          useCases,
			"valueProposition":  "Aggregated API Catalog",
			"version":           "v1",
			"visibility":        "public",
		},
		"license": map[string]interface{}{
			"governance": map[string]interface{}{
				"applicableLaws":  "Law",
				"audit":           "Annual",
				"confidentiality": "Standard",
				"damages":         "Limited",
				"forceMajeure":    "Standard",
				"ownership":       "Noi SPA",
				"warranties":      "None",
			},
			"scope": map[string]interface{}{
				"definition":       "Full",
				"exclusive":        false,
				"geographicalArea": []string{"EU"},
				"language":         "en",
				"permanent":        true,
				"restrictions":     "None",
				"rights":           []string{"Reuse"},
			},
			"termination": map[string]interface{}{
				"continuityConditions":  "Standard",
				"terminationConditions": "Breach",
			},
		},
		"pricingPlans": map[string]interface{}{
			"en": []map[string]interface{}{
				{
					"billingDuration":        "Monthly",
					"maxTransactionQuantity": "Unlimited",
					"name":                   "Free",
					"offering":               []string{"Basic access"},
					"price":                  "0",
					"priceCurrency":          "USD",
					"unit":                   "month",
				},
			},
		},
		"recommendedDataProducts": []string{
			"https://example.org/product1",
			"https://example.org/product2",
		},
		"support": map[string]interface{}{
			"documentationURL":  "https://opendatahub.com/",
			"email":             "help@opendatahub.com",
			"emailServiceHours": "9-17",
			"phoneNumber":       "",
			"phoneServiceHours": "9-17",
		},
	}

	details := map[string]interface{}{
		"summary":     "This is an API catalog for ODHS datasets.",
		"description": "Detailed description of the API product, aggregating multiple datasets from Open Data Hub.",
		"language":    "en",
		"metadata": map[string]interface{}{
			"datasetCount": len(datasets),
		},
	}

	return map[string]interface{}{
		"schema":  "source/schema/odps.yaml",
		"version": "3.1",
		"product": product,
		"details": details,
	}
}

// toODPS30 maps datasets to an ODPS 3.0 (dev) structure conforming to the provided interface.
func toODPS30(datasets []Dataset) map[string]interface{} {
	// Use the first dataset as basis (or default if none)
	var first Dataset
	if len(datasets) > 0 {
		first = datasets[0]
	} else {
		first = Dataset{
			ID:        "default",
			Self:      "https://example.org/default",
			Type:      "unknown",
			Shortname: "Default Dataset",
			ApiUrl:    "https://example.org/api/default",
		}
	}

	// Collect unique categories from all datasets.
	categoriesSet := make(map[string]bool)
	for _, ds := range datasets {
		if ds.Type != "" {
			categoriesSet[ds.Type] = true
		}
	}
	var categories []string
	for cat := range categoriesSet {
		categories = append(categories, cat)
	}

	// Build useCases from each dataset.
	var useCases []map[string]interface{}
	for _, ds := range datasets {
		useCases = append(useCases, map[string]interface{}{
			"useCase": map[string]interface{}{
				"useCaseTitle":       ds.Shortname,
				"useCaseDescription": fmt.Sprintf("Dataset type: %s", ds.Type),
				"useCaseURL":         ds.ApiUrl,
			},
		})
	}
	if len(useCases) == 0 {
		useCases = []map[string]interface{}{
			{"useCase": map[string]interface{}{
				"useCaseTitle":       "Default Dataset",
				"useCaseDescription": "Default dataset description",
				"useCaseURL":         "https://example.org/api/default",
			}},
		}
	}

	// Build the product object using the first dataset (or defaults).
	product := map[string]interface{}{
		"en": map[string]interface{}{
			"name":              first.Shortname,
			"productID":         first.ID,
			"valueProposition":  "Aggregated API Catalog",
			"description":       fmt.Sprintf("Catalog for dataset type: %s", first.Type),
			"productSeries":     "v1.0",
			"visibility":        "public",
			"status":            "active",
			"version":           "v1",
			"categories":        categories,
			"standards":         []string{"Default Standard"},
			"tags":              []string{"API", "Dataset"},
			"brandSlogan":       "Open Data for All",
			"type":              first.Type,
			"logoURL":           first.Self,
			"OutputFileFormats": []string{"JSON", "YAML"},
			"useCases":          useCases,
		},
	}

	// Default values for other properties.
	recommendedDataProducts := []string{
		"https://example.org/product1",
		"https://example.org/product2",
	}

	pricingPlans := map[string]interface{}{
		"en": []map[string]interface{}{
			{
				"name":                   "Free",
				"priceCurrency":          "USD",
				"price":                  "0",
				"billingDuration":        "Monthly",
				"unit":                   "month",
				"maxTransactionQuantity": "Unlimited",
				"offering":               []string{"Basic access"},
			},
		},
	}

	dataOps := map[string]interface{}{
		"data": map[string]interface{}{
			"schemaLocationURL": "https://example.org/schema",
		},
		"lineage": map[string]interface{}{
			"dataLineageTool":   "Default Lineage Tool",
			"dataLineageOutput": "Lineage output info",
		},
		"infrastructure": map[string]interface{}{
			"platform":         "Kubernetes",
			"region":           "us-east",
			"storageTechnology": "S3",
			"storageType":      "object",
			"containerTool":    "Docker",
		},
		"build": map[string]interface{}{
			"format":                     "docker",
			"hashType":                   "sha256",
			"checksum":                   "abc123",
			"signatureType":              "PGP",
			"scriptURL":                  "https://example.org/build.sh",
			"deploymentDocumentationURL": "https://example.org/deploy",
		},
	}

	dataAccess := map[string]interface{}{
		"type":                 "REST",
		"authenticationMethod": "None",
		"specification":        "OpenAPI",
		"format":               "JSON",
		"documentationURL":     "https://example.org/docs",
	}

	SLA := []map[string]interface{}{
		{
			"dimension": "Availability",
			"displaytitle": []map[string]interface{}{
				{"en": "Availability SLA"},
			},
			"objective": 99.9,
			"unit":      "%",
			"monitoring": map[string]interface{}{
				"type":      "Service Level",
				"reference": "https://example.org/monitoring",
				"spec":      "SLA Monitoring Spec",
			},
		},
	}

	support := map[string]interface{}{
		"phoneNumber":       "123-456-7890",
		"phoneServiceHours": "9-5",
		"email":             "help@opendatahub.com",
		"emailServiceHours": "9-5",
		"documentationURL":  "https://example.org/support/docs",
	}

	dataQuality := []map[string]interface{}{
		{
			"dimension": "Accuracy",
			"displaytitle": []map[string]interface{}{
				{"en": "Data Accuracy"},
			},
			"objective": 95.0,
			"unit":      "%",
			"monitoring": map[string]interface{}{
				"type":      "Quality",
				"reference": "https://example.org/quality",
				"spec":      "Quality Spec",
			},
		},
	}

	license := map[string]interface{}{
		"scope": map[string]interface{}{
			"definition":       "Full",
			"language":         "en",
			"restrictions":     "None",
			"geographicalArea": []string{"EU"},
			"permanent":        true,
			"exclusive":        false,
			"rights":           []string{"Reuse"},
		},
		"termination": map[string]interface{}{
			"terminationConditions": "Breach",
			"continuityConditions":  "Standard",
		},
		"governance": map[string]interface{}{
			"ownership":       "Noi SPA",
			"damages":         "Limited",
			"confidentiality": "Standard",
			"applicableLaws":  "Law X",
			"warranties":      "None",
			"audit":           "Annual",
			"forceMajeure":    "Standard",
		},
	}

	dataHolder := map[string]interface{}{
		"taxID":              "02595720216",
		"vatID":              "02595720216",
		"businessDomain":     "Tourism",
		"logoURL":            "https://opendatahub.com/img/NOI_OPENDATAHUB_NEW_BK_nospace-01.svg",
		"description":        "Data holder description",
		"URL":                "https://noi.bz.it",
		"telephone":          "+390471066600",
		"streetAddress":      "Via Volta 13/A",
		"postalCode":         "39100",
		"addressRegion":      "Alto Adige",
		"addressLocality":    "Bolzano",
		"addressCountry":     "IT",
		"aggregateRating":    "4.5",
		"ratingCount":        100,
		"slogan":             "Open Data for All",
		"parentOrganization": "Noi SPA",
	}

	result := map[string]interface{}{
		"schema":                  "source/schema/odps-dev-yaml-schema.yaml",
		"version":                 "dev",
		"product":                 product,
		"recommendedDataProducts": recommendedDataProducts,
		"pricingPlans":            pricingPlans,
		"dataOps":                 dataOps,
		"dataAccess":              dataAccess,
		"SLA":                     SLA,
		"support":                 support,
		"dataQuality":             dataQuality,
		"license":                 license,
		"dataHolder":              dataHolder,
	}

	return result
}

// --- Handlers ---

//dcatHandler has json as default and supports ?format=json.
func dcatHandler(w http.ResponseWriter, r *http.Request) {
	page := getPageNumber(r)
	datasets, err := fetchDatasets(page)
	if err != nil || len(datasets) == 0 {
		log.Println("Error fetching datasets or no data; using default data")
		datasets = getDefaultDatasets()
	}

	output := toDCAT(datasets)
	format := r.URL.Query().Get("format") // check query parameter

	if format == "yaml" {
		yamlData, err := yaml.Marshal(output)
		if err != nil {
			http.Error(w, "Error marshaling YAML", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write(yamlData)
	} else {
		// Default to JSON response
		jsonData, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			http.Error(w, "Error marshaling JSON", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	}
}


func odpsHandler(w http.ResponseWriter, r *http.Request) {
	page := getPageNumber(r)
	datasets, err := fetchDatasets(page)
	if err != nil || len(datasets) == 0 {
		log.Println("Error fetching datasets or no data; using default data")
		datasets = getDefaultDatasets()
	}

	output := toODPS(datasets)
	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		http.Error(w, "Error marshaling JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

// odps31Handler has YAML as default and supports ?format=json.
func odps31Handler(w http.ResponseWriter, r *http.Request) {
	page := getPageNumber(r)
	datasets, err := fetchDatasets(page)
	if err != nil || len(datasets) == 0 {
		log.Println("Error fetching datasets or no data; using default data")
		datasets = getDefaultDatasets()
	}

	output := toODPS31(datasets)
	format := r.URL.Query().Get("format")

	if format == "json" {
		// Return JSON if requested
		jsonData, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			http.Error(w, "Error marshaling JSON", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	} else {
		// Otherwise default to YAML
		yamlData, err := yaml.Marshal(output)
		if err != nil {
			http.Error(w, "Error marshaling YAML", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write(yamlData)
	}
}

// odps30Handler has YAML as default and supports ?format=json.
func odps30Handler(w http.ResponseWriter, r *http.Request) {
	page := getPageNumber(r)
	datasets, err := fetchDatasets(page)
	if err != nil || len(datasets) == 0 {
		log.Println("Error fetching datasets or no data; using default data")
		datasets = getDefaultDatasets()
	}

	output := toODPS30(datasets)
	format := r.URL.Query().Get("format")

	if format == "json" {
		// Return JSON if requested
		jsonData, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			http.Error(w, "Error marshaling JSON", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	} else {
		// Otherwise default to YAML
		yamlData, err := yaml.Marshal(output)
		if err != nil {
			http.Error(w, "Error marshaling YAML", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write(yamlData)
	}
}

// getPageNumber extracts the "page" query parameter.
func getPageNumber(r *http.Request) int {
	pageStr := r.URL.Query().Get("page")
	if pageStr == "" {
		return 1
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return 1
	}
	return page
}

func main() {
	http.HandleFunc("/dcat", dcatHandler)
	http.HandleFunc("/odps", odpsHandler)
	http.HandleFunc("/odps31", odps31Handler)
	http.HandleFunc("/odps30", odps30Handler)

	fmt.Println("Server running on http://localhost:8878")
	log.Fatal(http.ListenAndServe(":8878", nil))
}
