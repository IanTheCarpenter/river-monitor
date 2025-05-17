package forecaster_logic

import (
	"github.com/IanTheCarpenter/river-monitor/schemas"
)

func analyze(telemetry SiteData, data_site schemas.DataCollectionSite, current_forecast *schemas.Forecast) {

	successes := 0

	for i := range 10 {
		if data_site.Thresholds.Low < telemetry.Records[i].Stage && telemetry.Records[i].Stage < data_site.Thresholds.High {
			successes = successes + 1
		}
	}

	current_forecast.PointsOfInterestForecasts = append(current_forecast.PointsOfInterestForecasts, schemas.PointOfInterestForecast{
		Name:     data_site.Name,
		Runnable: successes > 6,
	})

}
