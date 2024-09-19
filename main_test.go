package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mfryhover/weather/api"

	"github.com/mfryhover/weather/cache"
)

func TestMain_displayPrompt(t *testing.T) {
	displayPrompt()
}

func TestMain_displayCurrentForecast(t *testing.T) {
	displayCurrentForecast("3001 Esperanza Crossing, Austin, TX 78758, USA", 78.6, 97.6, 75.8, false)
}

func TestMain_displayExtendedForecast(t *testing.T) {
	weeklyForecast := api.WeeklyForecast{
		Time:             []string{"2024-09-19"},
		Temperature2MMax: []float64{97.6},
		Temperature2MMin: []float64{75.8},
	}
	displayExtendedForecast(weeklyForecast)
}

func TestMain_getPostalCode(t *testing.T) {
	address := "3001 Esperanza Crossing, Austin, TX 78758, USA"
	expectedPostalCode := "78758"
	postalCode := getPostalCode(address)
	if postalCode != expectedPostalCode {
		t.Errorf("Expected postal code %s, got %s", expectedPostalCode, postalCode)
	}
}

func TestMain_getForecast(t *testing.T) {
	c := cache.GetCacheInstance()
	testcases := []struct {
		name                 string
		geocodeStatus        int
		forecastStatus       int
		err                  string
		address              string
		mockGeocodeResponse  string
		mockForecastResponse string
		isFromCache          bool
	}{
		{
			name:           "Success Case",
			geocodeStatus:  http.StatusOK,
			forecastStatus: http.StatusOK,
			mockForecastResponse: `{
							  "current": {
								"temperature_2m": 78.6
							  },
							  "daily": {
								"time": [
								  "2024-09-19"
								],
								"temperature_2m_max": [
								  97.6
								],
								"temperature_2m_min": [
								  75.8
								]
							  }
							}`,
			mockGeocodeResponse: `{
   								"results" : [{
									 "address_components" : [],
									 "formatted_address" : "3001 Esperanza Crossing, Austin, TX 78758, USA",
									 "geometry" : {
										"bounds" : {},
										"location" :
										{
										   "lat" : 30.3985991,
											"lng" : 30.3985991
										},
										"location_type" : "ROOFTOP",
										"viewport" :
										{ }
									 },
									 "place_id" : "ChIJAZqSbXPMRIYRNYouzErXl_4",
									 "types" : []
								}],
								"status" : "OK"
							}`,
			address: "3001 Esperanza Crossing, Austin, TX 78758, USA",
		},
		{
			name:           "Success Case - Retrieve from Cache",
			geocodeStatus:  http.StatusOK,
			forecastStatus: http.StatusOK,
			mockForecastResponse: `{
							  "current": {
								"temperature_2m": 78.6
							  },
							  "daily": {
								"time": [
								  "2024-09-19"
								],
								"temperature_2m_max": [
								  97.6
								],
								"temperature_2m_min": [
								  75.8
								]
							  }
							}`,
			mockGeocodeResponse: `{
   								"results" : [{
									 "address_components" : [],
									 "formatted_address" : "3001 Esperanza Crossing, Austin, TX 78758, USA",
									 "geometry" : {
										"bounds" : {},
										"location" :
										{
										   "lat" : 30.3985991,
											"lng" : 30.3985991
										},
										"location_type" : "ROOFTOP",
										"viewport" :
										{ }
									 },
									 "place_id" : "ChIJAZqSbXPMRIYRNYouzErXl_4",
									 "types" : []
								}],
								"status" : "OK"
							}`,
			address:     "3001 Esperanza Crossing, Austin, TX 78758, USA",
			isFromCache: true,
		},
		{
			name:           "Error - Geocode API Failed",
			geocodeStatus:  http.StatusNotFound,
			forecastStatus: http.StatusOK,
			mockForecastResponse: `{
							  "current": {
								"temperature_2m": 78.6
							  },
							  "daily": {
								"time": [
								  "2024-09-19"
								],
								"temperature_2m_max": [
								  97.6
								],
								"temperature_2m_min": [
								  75.8
								]
							  }
							}`,
			mockGeocodeResponse: `{
   								"results" : [{
									 "address_components" : [],
									 "formatted_address" : "3001 Esperanza Crossing, Austin, TX 78758, USA",
									 "geometry" : {
										"bounds" : {},
										"location" :
										{
										   "lat" : 30.3985991,
											"lng" : 30.3985991
										},
										"location_type" : "ROOFTOP",
										"viewport" :
										{ }
									 },
									 "place_id" : "ChIJAZqSbXPMRIYRNYouzErXl_4",
									 "types" : []
								}],
								"status" : "OK"
							}`,
			address: "3001 Esperanza Crossing, Austin, TX 78758, USA",
			err:     "error retrieving coordinates: received non-OK HTTP status: 404 Not Found",
		},
		{
			name:           "Error - Forecast API Failed",
			geocodeStatus:  http.StatusOK,
			forecastStatus: http.StatusNotFound,
			mockForecastResponse: `{
							  "current": {
								"temperature_2m": 78.6
							  },
							  "daily": {
								"time": [
								  "2024-09-19"
								],
								"temperature_2m_max": [
								  97.6
								],
								"temperature_2m_min": [
								  75.8
								]
							  }
							}`,
			mockGeocodeResponse: `{
   								"results" : [{
									 "address_components" : [],
									 "formatted_address" : "3001 Esperanza Crossing, Austin, TX 78752, USA",
									 "geometry" : {
										"bounds" : {},
										"location" :
										{
										   "lat" : 30.3985991,
											"lng" : 30.3985991
										},
										"location_type" : "ROOFTOP",
										"viewport" :
										{ }
									 },
									 "place_id" : "ChIJAZqSbXPMRIYRNYouzErXl_4",
									 "types" : []
								}],
								"status" : "OK"
							}`,
			address: "3001 Esperanza Crossing, Austin, TX 78758, USA",
			err:     "error retrieving forecast: received non-OK HTTP status: 404 Not Found",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/maps/api/geocode/json" {
					w.WriteHeader(tc.geocodeStatus)
					w.Write([]byte(tc.mockGeocodeResponse))
				}
				if r.URL.Path == "/v1/forecast" {
					w.WriteHeader(tc.forecastStatus)
					w.Write([]byte(tc.mockForecastResponse))
				}
			}))
			defer server.Close()

			fullAddress, currentTemp, weeklyForecast, isFromCache, err := getForecast(tc.address, c, server.URL, server.URL, "testApiKey")
			// Check for error cases
			if tc.err != "" {
				if err != nil {
					if err.Error() != tc.err {
						t.Errorf("Expected '%s', got %s", tc.err, err.Error())
					}
				} else {
					t.Errorf("Expected an error, got nil")
				}
				return
			}

			// Check for success cases
			if fullAddress != tc.address {
				t.Errorf("Expected '%s', got %s", tc.address, fullAddress)
			}
			if currentTemp == 0 {
				t.Errorf("Expected currentTemp to be greater than 0")
			}
			if len(weeklyForecast.Time) == 0 {
				t.Errorf("Expected weeklyForecast to have at least 1 day")
			}
			if isFromCache != tc.isFromCache {
				t.Errorf("Expected isFromCache to be %t, got %t", tc.isFromCache, isFromCache)
			}
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
		})
	}
}
