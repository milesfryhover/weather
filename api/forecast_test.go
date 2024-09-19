package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_GetForecast(t *testing.T) {
	tc := []struct {
		name           string
		mockResponse   string
		currentTemp    float64
		status         int
		weeklyForecast WeeklyForecast
		error          string
	}{
		{
			name:   "Success Case",
			status: http.StatusOK,
			mockResponse: `{
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
			currentTemp: 78.6,
			weeklyForecast: WeeklyForecast{
				Time:             []string{"2024-09-19"},
				Temperature2MMax: []float64{97.6},
				Temperature2MMin: []float64{75.8},
			},
		},
		{
			name:         "Error Unmarshalling",
			status:       http.StatusOK,
			mockResponse: `}`,
			currentTemp:  78.6,
			weeklyForecast: WeeklyForecast{
				Time:             []string{"2024-09-19"},
				Temperature2MMax: []float64{97.6},
				Temperature2MMin: []float64{75.8},
			},
			error: "error unmarshalling response body: invalid character '}' looking for beginning of value",
		},
		{
			name:         "Status Not OK",
			status:       http.StatusNotFound,
			mockResponse: `{}`,
			currentTemp:  78.6,
			weeklyForecast: WeeklyForecast{
				Time:             []string{"2024-09-19"},
				Temperature2MMax: []float64{97.6},
				Temperature2MMin: []float64{75.8},
			},
			error: "received non-OK HTTP status: 404 Not Found",
		},
	}

	for _, tc := range tc {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/v1/forecast" {
					t.Errorf("Expected to request '/v1/forecast', got: %s", r.URL.Path)
				}
				w.WriteHeader(tc.status)
				w.Write([]byte(tc.mockResponse))
			}))
			defer server.Close()

			currentTemp, weeklyForecast, err := GetForecast(0, 0, server.URL)
			// Check for error cases
			if tc.error != "" {
				if err != nil {
					if err.Error() != tc.error {
						t.Errorf("Expected '%s', got %s", tc.error, err.Error())
					}
				} else {
					t.Errorf("Expected an error, got nil")
				}
				return
			}

			// Check for success cases
			if currentTemp != tc.currentTemp {
				t.Errorf("Expected '%f', got %f", tc.currentTemp, currentTemp)
			}
			if len(weeklyForecast.Time) > 0 {
				if weeklyForecast.Time[0] != tc.weeklyForecast.Time[0] {
					t.Errorf("Expected '%s', got %s", tc.weeklyForecast.Time[0], weeklyForecast.Time[0])
				}
				if weeklyForecast.Temperature2MMin[0] != tc.weeklyForecast.Temperature2MMin[0] {
					t.Errorf("Expected '%f', got %f", tc.weeklyForecast.Temperature2MMin[0], weeklyForecast.Temperature2MMin[0])
				}
				if weeklyForecast.Temperature2MMax[0] != tc.weeklyForecast.Temperature2MMax[0] {
					t.Errorf("Expected '%f', got %f", tc.weeklyForecast.Temperature2MMax[0], weeklyForecast.Temperature2MMax[0])
				}
			}
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
		})
	}
}
