package main

import (
	"encoding/json"
	"fmt"

	"github.com/IanTheCarpenter/river-monitor/external_apis"
	"github.com/IanTheCarpenter/river-monitor/schemas"
)

func analyze(telemetry external_apis.SiteData, current_forecast *schemas.Forecast) {

	successes := 0

	for i := range 10 {
		if telemetry.StageThreshold < telemetry.Stage[i].Value && telemetry.Stage[i].Value < data_site.Thresholds.High {
			successes = successes + 1
		}
	}
	current_forecast.PointsOfInterestForecasts = append(current_forecast.PointsOfInterestForecasts, schemas.PointOfInterestForecast{
		Name:     data_site.Name,
		Runnable: successes > 6,
	})

}

func getData() {
	obj := make(map[string]any)

	myJsonData := []byte(`{"name": "John", "age": 30, "region": {"name": "Texas", "state": "TX"}}`)

	err := json.Unmarshal(myJsonData, &obj)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println(obj["name"])
	fmt.Println(obj["region"])

	myStruct := struct {
		Name   string `json:"name"`
		Age    int    `json:"age"`
		Region struct {
			Name  string `json:"name"`
			State string `json:"state"`
		} `json:"region"`
	}{}

	// check if region is present and if it is, do a type assertion to get the region
	region := obj["region"]

	var exists bool
	region, exists = obj["region"]
	if !exists {
		return
	}
	fmt.Println(region)

	// Do a type assertion to get the region
	if region, ok := obj["region"].(map[string]any); ok {
		name, ok := region["name"].(string)
		if ok {
			fmt.Println(name)
		}
		myStruct.Region.State = region["state"].(string)
	}

	// another way to do it:

	// check if region is present
	if region2, exists := obj["region"]; exists {
		if region, ok := region2.(map[string]any); ok {
			name, ok := region["name"].(string)
			if ok {
				fmt.Println(name)
			}
			myStruct.Region.State = region["state"].(string)
		}
	}

	// and a risky way to do it:
	myStruct.Region.Name = obj["region"].(map[string]any)["name"].(string)
	myStruct.Region.State = obj["region"].(map[string]any)["state"].(string)

}
