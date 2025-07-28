package main

import (
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
