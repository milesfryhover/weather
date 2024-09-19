// Package api abstracts the logic for making API calls to the Open-Meteo API and the Google Geocode API.
package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	// geocodePathTemplate defines the URL path template for the Google Geocode API request.
	geocodePathTemplate = "/maps/api/geocode/json?address=%s&key=%s"
)

// geocodeResponse holds the response from the Google Geocode API.
type geocodeResponse struct {
	// Results is a list of geocoding results.
	Results []geocodeResult `json:"results"`
}

// geocodeResult represents a single result from the geocode API response.
type geocodeResult struct {
	// FormattedAddress is the full address returned by the API.
	FormattedAddress string `json:"formatted_address"`
	// Geometry contains the location data.
	Geometry geocodeGeometry `json:"geometry"`
}

// geocodeGeometry contains the location details.
type geocodeGeometry struct {
	// Location specifies the latitude and longitude.
	Location geocodeLocation `json:"location"`
}

// geocodeLocation represents the latitude and longitude coordinates.
type geocodeLocation struct {
	// Lat is the latitude of the location.
	Lat float64 `json:"lat"`
	// Lng is the longitude of the location.
	Lng float64 `json:"lng"`
}

// AddressToCoordinates converts an address into geographical coordinates.
// It returns the full formatted address, latitude, longitude, and an error if any.
// If the address is not found or an error occurs, it returns zero values and the error.
func AddressToCoordinates(address string, baseURL string, apiKey string) (fullAddress string, latitude, longitude float64, err error) {
	// Build the full API request URL
	path := fmt.Sprintf(geocodePathTemplate, url.QueryEscape(address), apiKey)
	fullURL := fmt.Sprintf(baseURL + path)

	// Make the HTTP GET request to the API
	resp, err := http.Get(fullURL)
	if err != nil {
		return "", 0.0, 0.0, fmt.Errorf("error making GET request: %v", err)
	}
	defer resp.Body.Close()

	// Check the HTTP status code
	if resp.StatusCode != http.StatusOK {
		return "", 0.0, 0.0, fmt.Errorf("received non-OK HTTP status: %s", resp.Status)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0.0, 0.0, fmt.Errorf("error reading response body: %v", err)
	}

	// Unmarshal the JSON data into the geocodeResponse struct
	var googleRes geocodeResponse
	err = json.Unmarshal(body, &googleRes)
	if err != nil {
		return "", 0.0, 0.0, fmt.Errorf("error unmarshalling response body: %v", err)
	}

	// Check if any results were returned
	if len(googleRes.Results) == 0 {
		return "", 0.0, 0.0, fmt.Errorf("no results found for address: %s", address)
	}

	// Return the first result's formatted address and coordinates - assuming the first result is the most relevant
	result := googleRes.Results[0]
	return result.FormattedAddress, result.Geometry.Location.Lat, result.Geometry.Location.Lng, nil
}
