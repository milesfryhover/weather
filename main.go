// Package main is the entry point for the World's Best Weather App.
// It interacts with the user to provide current and extended weather forecasts.
package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/mfryhover/weather/api"
	"github.com/mfryhover/weather/cache"
)

// displayPrompt displays the user prompt instructions.
func displayPrompt() {
	fmt.Println("To exit please enter q")
	fmt.Println("Otherwise, please enter your address")
	fmt.Print("-> ")
}

// displayCurrentForecast displays the current weather forecast for the given address.
// It shows the current temperature, today's high and low, and indicates if the data was retrieved from the cache.
func displayCurrentForecast(address string, currentTemp, maxTemp, minTemp float64, isFromCache bool) {
	fmt.Println()
	if isFromCache {
		fmt.Println("***Retrieved forecast from cache***")
	}
	fmt.Printf("Here is the weather for address: %s\n", address)
	fmt.Println("---------------------------")
	fmt.Printf("The current temperature is %.1f F\n", currentTemp)
	fmt.Printf("The high for today is %.1f F\n", maxTemp)
	fmt.Printf("The low for today is %.1f F\n", minTemp)
	fmt.Println()
}

// displayExtendedForecast displays the extended weather forecast for the week.
func displayExtendedForecast(weeklyForecast api.WeeklyForecast) {
	fmt.Println("Extended Forecast: ")
	fmt.Println("---------------------------")
	for dayIndex := range weeklyForecast.Time {
		fmt.Println(weeklyForecast.Time[dayIndex])
		fmt.Printf("Max Temp: %.1f F\n", weeklyForecast.Temperature2MMax[dayIndex])
		fmt.Printf("Min Temp: %.1f F\n", weeklyForecast.Temperature2MMin[dayIndex])
		fmt.Println("--------------------")
	}
	fmt.Println()
}

// getPostalCode extracts and returns the postal code from a full address string.
// For example, given an address like "3001 Esperanza Crossing, Austin, TX 78758, USA",
// it returns "78758".
func getPostalCode(address string) string {
	statePostalCode := regexp.MustCompile(`\b[A-Z]{2}\s\d{5}\b`)
	sp := statePostalCode.FindString(address)

	if len(sp) == 8 {
		return sp[len(sp)-5:]
	}

	return ""
}

// getForecast retrieves the current temperature and weekly forecast for the given address.
// It returns the full formatted address, current temperature, weekly forecast, and a boolean indicating if the data was retrieved from the cache.
func getForecast(address string, c *cache.Cache, geocodeURL string, forecastURL string, apiKey string) (string, float64, api.WeeklyForecast, bool, error) {
	// Get the latitude and longitude of the address
	addressFull, lat, lng, err := api.AddressToCoordinates(address, geocodeURL, apiKey)
	if err != nil || addressFull == "" || lat == 0 || lng == 0 {
		return "", 0, api.WeeklyForecast{}, false, fmt.Errorf("error retrieving coordinates: %v", err)
	}

	// Get the postal code from the address
	pc := getPostalCode(addressFull)

	// Get the current temperature and weekly forecast
	isFromCache := true
	currentTemp, weeklyForecast, ok := c.Get(pc)
	if !ok {
		currentTemp, weeklyForecast, err = api.GetForecast(lat, lng, forecastURL)
		if err != nil {
			return "", 0, api.WeeklyForecast{}, false, fmt.Errorf("error retrieving forecast: %v", err)
		}
		c.Add(pc, currentTemp, weeklyForecast)
		isFromCache = false
		return addressFull, currentTemp, weeklyForecast, isFromCache, nil
	}

	return addressFull, currentTemp, weeklyForecast, isFromCache, nil
}

func main() {
	c := cache.GetCacheInstance()
	c.StartAutoPurge(1 * time.Hour)

	// Retrieve the API key once
	apiKey := os.Getenv("GEOCODE_API_KEY")
	if apiKey == "" {
		fmt.Println("GEOCODE_API_KEY environment variable is not set.")
		os.Exit(1)
	}

	fmt.Println("World's Best Weather App")
	fmt.Println("---------------------------")

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		displayPrompt()
		address := scanner.Text()

		if strings.EqualFold(address, "q") {
			fmt.Println("Thanks for using the World's Best Weather App!")
			break
		}

		addressFull, currentTemp, weeklyForecast, isFromCache, err := getForecast(address, c, "https://maps.googleapis.com", "https://api.open-meteo.com", apiKey)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if len(weeklyForecast.Temperature2MMax) > 0 && len(weeklyForecast.Temperature2MMin) > 0 && len(weeklyForecast.Time) > 0 {
			displayCurrentForecast(addressFull, currentTemp, weeklyForecast.Temperature2MMax[0], weeklyForecast.Temperature2MMin[0], isFromCache)
			displayExtendedForecast(weeklyForecast)
		} else {
			fmt.Println("Forecast data is unavailable. Please try again!")
		}

	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
	}

}
