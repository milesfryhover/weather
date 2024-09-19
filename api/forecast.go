// Package api abstracts the logic for making API calls to the Open-Meteo API and the Google Geocode API. It also contains
// the WeeklyForecast struct that holds the daily forecast from the Open-Meteo API.
package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	// forecastPathTemplate defines the URL path template for fetching forecast data from the Open-Meteo API.
	forecastPathTemplate = "/v1/forecast?latitude=%f&longitude=%f&current=temperature_2m&daily=temperature_2m_max,temperature_2m_min&temperature_unit=fahrenheit&wind_speed_unit=mph&precipitation_unit=inch"
)

// forecastResponse holds the forecast response from the API
type forecastResponse struct {
	Current struct {
		// Temperature2M represents the current temp in Fahrenheit
		Temperature2M float64 `json:"temperature_2m"`
	} `json:"current"`
	// WeeklyForecast contains the daily forecast data for a week
	WeeklyForecast `json:"daily"`
}

// WeeklyForecast holds the daily forecast from the Open-Meteo API.
// It includes dates and temperature ranges in slices, where each index corresponds to the same day.
type WeeklyForecast struct {
	// Time is a slice of dates covered by the forecast.
	Time []string `json:"time"`
	// Temperature2MMax holds the max temperatures for each day.
	Temperature2MMax []float64 `json:"temperature_2m_max"`
	// Temperature2MMin holds the min temperatures for each day.
	Temperature2MMin []float64 `json:"temperature_2m_min"`
}

// GetForecast retrieves the current temperature and weekly forecast for the given latitude and longitude.
// It requires the base URL of the API server and returns the current temperature, weekly forecast, and an error if any.
func GetForecast(latitude float64, longitude float64, baseURL string) (float64, WeeklyForecast, error) {
	// Build the full API request URL
	path := fmt.Sprintf(forecastPathTemplate, latitude, longitude)
	url := fmt.Sprintf(baseURL + path)

	// Make the HTTP GET request to the API
	resp, err := http.Get(url)
	if err != nil {
		return 0, WeeklyForecast{}, fmt.Errorf("error making GET request: %v", err)
	}
	defer resp.Body.Close()

	// Check the HTTP status code
	if resp.StatusCode != http.StatusOK {
		return 0, WeeklyForecast{}, fmt.Errorf("received non-OK HTTP status: %s", resp.Status)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, WeeklyForecast{}, fmt.Errorf("error reading response body: %v", err)
	}

	// Unmarshal the JSON data into the forecast struct
	forecast := forecastResponse{}
	err = json.Unmarshal(body, &forecast)
	if err != nil {
		return 0, WeeklyForecast{}, fmt.Errorf("error unmarshalling response body: %v", err)
	}

	// Return the current temperature and weekly forecast
	return forecast.Current.Temperature2M, forecast.WeeklyForecast, nil
}
