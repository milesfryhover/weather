package api

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func Test_AddressToCoordinates(t *testing.T) {
	os.Setenv("GEOCODE_API_KEY", "test")
	tc := []struct {
		name         string
		mockResponse string
		status       int
		address      string
		lat          float64
		long         float64
		err          string
	}{
		{
			name:   "Success Case",
			status: http.StatusOK,
			mockResponse: `{
   								"results" : [{
									 "address_components" : [],
									 "formatted_address" : "3001 Esperanza Crossing, Austin, TX 78758, USA",
									 "geometry" : {
										"bounds" : {},
										"location" :
										{
										   "lat" : 30.3985991,
										   "lng" : -97.72206659999999
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
			lat:     30.3985991,
			long:    -97.72206659999999,
		},
		{
			name:         "Error unmarshalling",
			status:       http.StatusOK,
			mockResponse: `}`,
			err:          "error unmarshalling response body: invalid character '}' looking for beginning of value",
		},
		{
			name:         "Status Not OK",
			mockResponse: `}`,
			status:       http.StatusNotFound,
			err:          "received non-OK HTTP status: 404 Not Found",
		},
		{
			name:   "No Results",
			status: http.StatusOK,
			mockResponse: `{
   								"results" : [],
								"status" : "OK"
							}`,
			err: "no results found for address: test",
		},
		{
			name:   "Missing FormattedAddress",
			status: http.StatusOK,
			mockResponse: `{
   								"results" : [{
									 "address_components" : [],
									 "geometry" : {
										"bounds" : {},
										"location" :
										{
										   "lat" : 30.3985991,
										   "lng" : -97.72206659999999
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
			address: "",
			lat:     30.3985991,
			long:    -97.72206659999999,
		},
		{
			name:   "Missing Lat",
			status: http.StatusOK,
			mockResponse: `{
   								"results" : [{
									 "address_components" : [],
									 "formatted_address" : "3001 Esperanza Crossing, Austin, TX 78758, USA",
									 "geometry" : {
										"bounds" : {},
										"location" :
										{
										   "lng" : -97.72206659999999
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
			lat:     0,
			long:    -97.72206659999999,
		},
		{
			name:   "Missing Long",
			status: http.StatusOK,
			mockResponse: `{
   								"results" : [{
									 "address_components" : [],
									 "formatted_address" : "3001 Esperanza Crossing, Austin, TX 78758, USA",
									 "geometry" : {
										"bounds" : {},
										"location" :
										{
										   "lat" : 30.3985991
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
			lat:     30.3985991,
			long:    0,
		},
	}

	for _, tc := range tc {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/maps/api/geocode/json" {
					t.Errorf("Expected to request '/maps/api/geocode/json', got: %s", r.URL.Path)
				}
				w.WriteHeader(tc.status)
				w.Write([]byte(tc.mockResponse))
			}))
			defer server.Close()

			address, lat, long, err := AddressToCoordinates("test", server.URL, "TestAPIKey")
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
			if address != tc.address {
				t.Errorf("Expected '%s', got %s", tc.address, address)
			}
			if lat != tc.lat {
				t.Errorf("Expected '%f', got %f", tc.lat, lat)
			}
			if long != tc.long {
				t.Errorf("Expected '%f', got %f", tc.long, long)
			}
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
		})
	}
}
