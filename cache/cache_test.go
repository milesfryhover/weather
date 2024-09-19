package cache

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/mfryhover/weather/api"
)

var c = GetCacheInstance()

func TestCache_Add(t *testing.T) {
	key := "TestCache_Add"
	currentTemp := 75.5
	weeklyForecast := api.WeeklyForecast{
		Time:             []string{"2024-09-25"},
		Temperature2MMax: []float64{75.5},
		Temperature2MMin: []float64{75.2},
	}

	// Add entry to cacheInstance
	c.Add(key, currentTemp, weeklyForecast)

	// Get entry from cacheInstance
	temp, forecast, ok := c.Get(key)
	if !ok {
		t.Errorf("Expected key %s to be found in cacheInstance", key)
	}
	if temp != currentTemp {
		t.Errorf("Expected temp %f, got %f", currentTemp, temp)
	}
	if !strings.EqualFold(forecast.Time[0], weeklyForecast.Time[0]) {
		t.Errorf("Expected forecast time %+v, got %+v", weeklyForecast.Time[0], forecast.Time[0])
	}
	if forecast.Temperature2MMax[0] != weeklyForecast.Temperature2MMax[0] {
		t.Errorf("Expected forecast max %f, got %f", weeklyForecast.Temperature2MMax[0], forecast.Temperature2MMax[0])
	}
	if forecast.Temperature2MMin[0] != weeklyForecast.Temperature2MMin[0] {
		t.Errorf("Expected forecast min %f, got %f", weeklyForecast.Temperature2MMin[0], forecast.Temperature2MMin[0])
	}
}

func TestCache_Get(t *testing.T) {
	c.SetEntryTTL(1 * time.Second)
	key := "TestCache_Get"
	currentTemp := 75.5
	weeklyForecast := api.WeeklyForecast{
		Time:             []string{"2024-09-25"},
		Temperature2MMax: []float64{75.5},
		Temperature2MMin: []float64{75.2},
	}

	// Add entry to cacheInstance
	c.Add(key, currentTemp, weeklyForecast)

	// Sleep so it will expire
	time.Sleep(2 * time.Second)

	// Try to get expired entry
	temp, forecast, ok := c.Get(key)
	if ok {
		t.Errorf("Expected key %s to be removed from the cacheInstance", key)
	}
	if temp != 0 {
		t.Error("Expected temp 0 for expired entry")
	}
	if !reflect.DeepEqual(forecast, api.WeeklyForecast{}) {
		t.Error("Expected forecast to be empty for expired entry")
	}

	// Try to get entry that doesn't exist
	temp, forecast, ok = cacheInstance.Get(key)
	if ok {
		t.Errorf("Expected key %s to not exist", key)
	}
	if temp != 0 {
		t.Error("Expected temp 0 for entry that doesn't exist")
	}
	if !reflect.DeepEqual(forecast, api.WeeklyForecast{}) {
		t.Error("Expected forecast to be empty for entry that doesn't exist")
	}
}

func TestCache_Delete(t *testing.T) {
	key := "TestCache_Delete"
	currentTemp := 75.5
	weeklyForecast := api.WeeklyForecast{
		Time:             []string{"2024-09-25"},
		Temperature2MMax: []float64{75.5},
		Temperature2MMin: []float64{75.2},
	}

	c.Add(key, currentTemp, weeklyForecast)

	c.Delete(key)

	temp, forecast, ok := c.Get(key)
	if ok {
		t.Errorf("Expected key %s to be deleted from cacheInstance", key)
	}
	if temp != 0 {
		t.Errorf("Expected temp 0 for deleted entry, got %f", temp)
	}
	if !reflect.DeepEqual(forecast, api.WeeklyForecast{}) {
		t.Error("Expected forecast to be empty for deleted entry")
	}
}

func TestCache_PurgeCache(t *testing.T) {
	c.SetEntryTTL(1 * time.Second)
	key := "TestCache_PurgeCache"
	key2 := "TestCache_PurgeCache2"
	currentTemp := 75.5
	weeklyForecast := api.WeeklyForecast{
		Time:             []string{"2024-09-25"},
		Temperature2MMax: []float64{75.5},
		Temperature2MMin: []float64{75.2},
	}

	c.Add(key, currentTemp, weeklyForecast)
	time.Sleep(2 * time.Second)
	c.Add(key2, currentTemp, weeklyForecast)
	c.PurgeCache()

	if _, _, ok := c.Get(key); ok {
		t.Errorf("Expected key %s to be purged from cacheInstance", key2)
	}
	if _, _, ok := c.Get(key2); !ok {
		t.Errorf("Expected key %s to remain in cacheInstance", key2)
	}
}
