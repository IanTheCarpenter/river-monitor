package forecaster

import (
	"context"
	"fmt"
	"time"

	"github.com/IanTheCarpenter/river-monitor/db"
	"github.com/IanTheCarpenter/river-monitor/forecaster/schemas"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type SiteData struct {
	SiteName string
	// HasFlow bool
	// HasStage bool
	Records []SiteSample
}
type SiteSample struct {
	TimeStamp time.Time
	Stage     float64
	Flow      float64
}

func main() {
	NS_TO_MINUTES := 60000000000

	// begin loop
	for {
		fmt.Println("Regenerating Forecasts...")
		rivers, err := fetch_river_definitions()
		if err != nil {
			fmt.Println("UNABLE TO FETCH RIVERS")
			fmt.Println(err.Error())
		}

		for _, current_river := range rivers {
			fmt.Printf("...Generating forecast for: %s river\n", current_river.RiverName)
			var current_forecast = schemas.Forecast{
				River:                     current_river.RiverName,
				RiverObjectID:             current_river.ObjectID,
				LastUpdated:               time.Now().GoString(),
				PointsOfInterestForecasts: []schemas.PointOfInterestForecast{},
			}

			for _, data_site := range current_river.DataCollectionSites {
				var telemetry SiteData
				switch data_site.Agency {
				case "lcra":
					telemetry.Records, err = fetch_lcra_data(data_site.URL)
				case "usgs":
					telemetry.Records, err = fetch_usgs_data(data_site.URL)
				}
				if err != nil {
					fmt.Printf("Error fetching data for site: %s", data_site.Name)
				}

				analyze(telemetry, data_site, &current_forecast)

			}

			insert_forecast(current_forecast, current_river.ObjectID)
		}
		fmt.Println("...Done!")
		time.Sleep(time.Duration(5 * NS_TO_MINUTES))
	}
}

func insert_forecast(forecast schemas.Forecast, river_objectID bson.ObjectID) {
	filter := bson.D{{Key: "river_object_id", Value: river_objectID}}

	forecast_bson, err := bson.Marshal(forecast)
	if err != nil {
		fmt.Printf("Unable to marshal forecast object: %s\n", forecast.River)
		return
	}

	replace_result := db.RIVER_REPORTS.FindOneAndReplace(context.TODO(), filter, forecast_bson)
	if replace_result.Err() == mongo.ErrNoDocuments {
		// insert new document
		db.RIVER_REPORTS.InsertOne(context.TODO(), forecast_bson)
	}
}
