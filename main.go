package main

import (
	"context"
	"fmt"
	"log"

	"github.com/IanTheCarpenter/river-monitor/db"
	"github.com/IanTheCarpenter/river-monitor/riverdata"
	"github.com/IanTheCarpenter/river-monitor/schemas"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func main() {
	db.Init()
	rebuild_all_rivers()

	// NS_TO_MINUTES := 60000000000

	// // begin loop
	// for {
	// 	fmt.Println("Regenerating Forecasts...")
	// 	rivers, err := external_apis.Fetch_river_definitions()
	// 	if err != nil {
	// 		fmt.Println("UNABLE TO FETCH RIVERS")
	// 		fmt.Println(err.Error())
	// 	}

	// 	for _, current_river := range rivers {
	// 		fmt.Printf("...Generating forecast for: %s river\n", current_river.RiverName)
	// 		var current_forecast = schemas.Forecast{
	// 			River:                     current_river.RiverName,
	// 			RiverObjectID:             current_river.ObjectID,
	// 			LastUpdated:               time.Now().GoString(),
	// 			PointsOfInterestForecasts: []schemas.PointOfInterestForecast{},
	// 		}

	// 		for _, data_site := range current_river.DataCollectionSites {
	// 			var telemetry = external_apis.SiteData{
	// 				SiteName: data_site.Name,
	// 			}
	// 			switch data_site.Agency {
	// 			case "lcra":
	// 				telemetry.Records, err = external_apis.Fetch_lcra_data(data_site.URL)
	// 			case "usgs":
	// 				telemetry.Records, err = external_apis.Fetch_usgs_data(data_site.URL)
	// 			}
	// 			if err != nil {
	// 				fmt.Printf("Error fetching data for site: %s", data_site.Name)
	// 			}

	// 			analyze(telemetry, data_site, &current_forecast)

	// 		}
	// 		insert_forecast(current_forecast, current_river.ObjectID)
	// 	}
	// 	fmt.Println("...Done!")
	// 	time.Sleep(time.Duration(5 * NS_TO_MINUTES))
	// }
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
		fmt.Printf("Inserting new river forecast for: %s\n", forecast.River)

		db.RIVER_REPORTS.InsertOne(context.TODO(), forecast_bson)
	} else {
		fmt.Printf("Overwrote river forecast for: %s\n", forecast.River)

	}
}

func rebuild_all_rivers() {

	deleted, err := db.RIVER_DEFINITIONS.DeleteMany(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal("FATAL: UNABLE TO CLEAR RIVER SCHEMA DB")
	}

	fmt.Printf("Deleted %d river schemas!", deleted.DeletedCount)

	rivers := []string{
		"2639515",
	}

	for _, riverID := range rivers {
		current_river, err := riverdata.Build_river(riverID)
		if err != nil {
			fmt.Sprintf("ERROR: Build failed on river ID: %s\n%s", riverID, err)
			continue
		}
		fmt.Println(current_river.RiverName)
		insert_river(current_river)

	}

}

//TODO:: troubleshoot why the schemas are being built incorrectly

func insert_river(input schemas.River) {
	fmt.Printf("inserting river: %s", input.RiverName)
	river_bson, err := bson.Marshal(input)
	if err != nil {
		fmt.Printf("Unable to marshal river schema: %s\n", input.RiverName)
		return
	}

	db.RIVER_DEFINITIONS.InsertOne(context.TODO(), river_bson)

}
