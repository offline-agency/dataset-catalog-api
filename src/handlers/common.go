// © 2024 NOI Techpark <digital@noi.bz.it>
// SPDX-License-Identifier: AGPL-3.0-or-later

package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"opendatahub.com/dataset-catalog-api/transformers"
)

// ConvertDatasets maps a slice of handlers.Dataset (with all properties)
// into a slice of transformers.Dataset.
func ConvertDatasets(d []transformers.Dataset) []transformers.Dataset {
	var out []transformers.Dataset
	for _, ds := range d {
		out = append(out, transformers.Dataset{
			ID:             ds.ID,
			Self:           ds.Self,
			Type:           ds.Type,
			Meta:           transformers.MetaData(ds.Meta),
			ApiUrl:         ds.ApiUrl,
			Output:         ds.Output,
			ApiType:        ds.ApiType,
			BaseUrl:        ds.BaseUrl,
			ODHTags:        ds.ODHTags,
			OdhType:        ds.OdhType,
			Sources:        ds.Sources,
			Category:       ds.Category,
			ApiAccess:      ds.ApiAccess,
			ApiFilter:      ds.ApiFilter,
			Dataspace:      ds.Dataspace,
			OdhTagIds:      ds.OdhTagIds,
			PathParam:      ds.PathParam,
			Shortname:      ds.Shortname,
			Deprecated:     ds.Deprecated,
			LastChange:     ds.LastChange,
			SwaggerUrl:     ds.SwaggerUrl,
			FirstImport:    ds.FirstImport,
			LicenseInfo:    transformers.LicenseInfo(ds.LicenseInfo),
			PublishedOn:    ds.PublishedOn,
			RecordCount:    ds.RecordCount,
			DataProvider:   ds.DataProvider,
			ImageGallery:   convertImageGallery(ds.ImageGallery),
			ApiDescription: ds.ApiDescription,
		})
	}
	return out
}

// convertImageGallery maps a slice of handlers.ImageGalleryItem to a slice of transformers.ImageGalleryItem.
func convertImageGallery(src []transformers.ImageGalleryItem) []transformers.ImageGalleryItem {
	var out []transformers.ImageGalleryItem
	for _, item := range src {
		out = append(out, transformers.ImageGalleryItem{
			Width:         item.Width,
			Height:        item.Height,
			License:       item.License,
			ValidTo:       item.ValidTo,
			ImageUrl:      item.ImageUrl,
			CopyRight:     item.CopyRight,
			ImageDesc:     item.ImageDesc,
			ImageName:     item.ImageName,
			ImageTags:     item.ImageTags,
			ValidFrom:     item.ValidFrom,
			ImageTitle:    item.ImageTitle,
			ImageSource:   item.ImageSource,
			IsInGallery:   item.IsInGallery,
			ImageAltText:  item.ImageAltText,
			ListPosition:  item.ListPosition,
			LicenseHolder: item.LicenseHolder,
		})
	}
	return out
}

const pageSize = 10

type cacheItem struct {
	data       []transformers.Dataset
	expiration time.Time
}

var (
	datasetCache = make(map[int]cacheItem)
	cacheMutex   sync.RWMutex
)

// fetchDatasets retrieves datasets for a given page from the external API,
// caching the result for 5 minutes.
func fetchDatasets(page int) ([]transformers.Dataset, error) {
	cacheMutex.RLock()
	if item, found := datasetCache[page]; found {
		if time.Now().Before(item.expiration) {
			cacheMutex.RUnlock()
			return item.data, nil
		}
	}
	cacheMutex.RUnlock()

	url := fmt.Sprintf("https://tourism.api.opendatahub.com/v1/MetaData?pagenumber=%d&limit=%d", page, pageSize)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data struct {
		TotalResults int       `json:"TotalResults"`
		TotalPages   int       `json:"TotalPages"`
		CurrentPage  int       `json:"CurrentPage"`
		NextPage     string    `json:"NextPage"`
		Items        []transformers.Dataset `json:"Items"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Printf("Error decoding JSON on page %d: %v", page, err)
		return nil, err
	}
	if len(data.Items) == 0 {
		log.Printf("No datasets found on page %d", page)
		return nil, nil
	}
	cacheMutex.Lock()
	datasetCache[page] = cacheItem{
		data:       data.Items,
		expiration: time.Now().Add(5 * time.Minute),
	}
	cacheMutex.Unlock()
	return data.Items, nil
}

// fetchDatasetsResponse retrieves the complete API response for a given page.
func fetchDatasetsResponse(page int) (*struct {
	TotalResults int       `json:"TotalResults"`
	TotalPages   int       `json:"TotalPages"`
	CurrentPage  int       `json:"CurrentPage"`
	NextPage     string    `json:"NextPage"`
	Items        []transformers.Dataset `json:"Items"`
}, error) {
	url := fmt.Sprintf("https://tourism.api.opendatahub.com/v1/MetaData?pagenumber=%d&limit=%d", page, pageSize)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data struct {
		TotalResults int       `json:"TotalResults"`
		TotalPages   int       `json:"TotalPages"`
		CurrentPage  int       `json:"CurrentPage"`
		NextPage     string    `json:"NextPage"`
		Items        []transformers.Dataset `json:"Items"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Printf("Error decoding JSON on page %d: %v", page, err)
		return nil, err
	}
	if len(data.Items) == 0 {
		return nil, nil
	}
	return &data, nil
}

// getDefaultDatasets returns a default dataset (not used if real data is available).
func getDefaultDatasets() []transformers.Dataset {
	return []transformers.Dataset{
		{
			ID:        "default",
			Self:      "https://example.org/default",
			Type:      "unknown",
			Shortname: "Default Dataset",
			ApiUrl:    "https://example.org/api/default",
		},
	}
}

// getPageNumber extracts the "page" query parameter from the request (default=1).
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

// slugify converts a string into a slug.
func slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, " ", "-")
	re := regexp.MustCompile(`[^a-z0-9\-]`)
	return re.ReplaceAllString(s, "")
}

// searchDatasetByID fetches the dataset details directly from the external API using the given ID.
func searchDatasetByID(id string) *transformers.Dataset {
	log.Printf("Directly fetching dataset detail for ID: %s", id)
	url := fmt.Sprintf("https://tourism.api.opendatahub.com/v1/MetaData/%s", id)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching detail for ID %s: %v", id, err)
		return nil
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		log.Printf("Dataset with ID %s not found (404)", id)
		return nil
	}
	var ds transformers.Dataset
	if err := json.NewDecoder(resp.Body).Decode(&ds); err != nil {
		log.Printf("Error decoding dataset detail for ID %s: %v", id, err)
		return nil
	}
	log.Printf("Dataset found: ID: %s, Shortname: %s", ds.ID, ds.Shortname)
	return &ds
}
