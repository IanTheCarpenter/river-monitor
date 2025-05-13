package forecaster

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/IanTheCarpenter/river-monitor/db"
	"github.com/IanTheCarpenter/river-monitor/schemas"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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

func Init() {

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

func fetch_river_definitions() ([]schemas.River, error) {

	// var bson_result bson.D
	// var result river.Schema

	bson_response, err := db.RIVER_DEFINITIONS.Find(context.TODO(), bson.D{{}}, options.Find())
	if err != nil {
		return nil, err
	}

	var results []schemas.River

	for bson_response.Next(context.TODO()) {
		var element schemas.River

		err := bson_response.Decode(&element)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		results = append(results, element)

	}

	return results, nil
}

func fetch_data(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return []byte{}, err
	}
	return body, nil
}

func fetch_lcra_data(url string) ([]SiteSample, error) {
	var output []SiteSample
	raw_data, err := fetch_data(url)
	if err != nil {
		return output, err
	}
	// define schema for this API
	type lcra_records struct {
		DateTime string  `json:"dateTime"`
		Stage    float64 `json:"value1"`
		Flow     float64 `json:"value2"`
	}

	type lcra_data struct {
		SiteName string         `json:"siteName"`
		Records  []lcra_records `json:"records"`
	}

	// convert this data into type []SiteSample
	var lcra lcra_data
	json.Unmarshal(raw_data, &lcra)

	for _, sample := range lcra.Records {
		var current_record SiteSample

		timestamp, err := time.Parse(time.RFC3339, sample.DateTime)
		if err != nil {
			fmt.Printf("Could not parse timestamp for data from: %s", url)
			continue
		}

		current_record.TimeStamp = timestamp
		current_record.Flow = sample.Flow
		current_record.Stage = sample.Stage

		output = insertDataInDescendingOrder(output, current_record)
	}

	return output, nil

}

func fetch_usgs_data(url string) ([]SiteSample, error) {
	// this is currently STAGE ONLY
	var output []SiteSample
	raw_data, err := fetch_data(url)
	if err != nil {
		return output, err
	}

	// i do not like this sam i am
	type usgs_data_sample struct {
		Stage    string `json:"value"`
		DateTime string `json:"dateTime"`
	}

	type usgs_values2 struct {
		Value []usgs_data_sample `json:"value"`
	}

	type usgs_timeseries struct {
		Values []usgs_values2 `json:"values"`
	}

	type usgs_value1 struct {
		TimeSeries []usgs_timeseries `json:"timeSeries"`
	}

	type usgs_data struct {
		Value usgs_value1 `json:"value"`
	}
	var usgs usgs_data
	json.Unmarshal(raw_data, &usgs)

	// convert this data into type []SiteSample

	for _, sample := range usgs.Value.TimeSeries[0].Values[0].Value {
		var current_record SiteSample
		timestamp, err := time.Parse(time.RFC3339, sample.DateTime)
		if err != nil {
			fmt.Printf("Could not parse timestamp for data from: %s", url)
			continue
		}
		float, err := strconv.ParseFloat(sample.Stage, 64)
		if err != nil {
			fmt.Printf("Could not parse stage into string for data from: %s", url)
			continue
		}

		current_record.TimeStamp = timestamp
		current_record.Stage = float
		current_record.Flow = 0

		output = insertDataInDescendingOrder(output, current_record)
	}

	return output, nil

}

func insertDataInDescendingOrder(data []SiteSample, sampletoinsert SiteSample) []SiteSample {
	if len(data) == 0 || sampletoinsert.TimeStamp.After(data[0].TimeStamp) {
		data = append(data, sampletoinsert)
	} else {
		data = append([]SiteSample{sampletoinsert}, data...)
	}
	return data
}

// PLACEHOLDER

// func fetch_usgs_data(url string) ([]SiteSample, error) {
// 	var output []SiteSample
// 	raw_data, err := fetch_data(url)
// 	if err != nil {
// 		return output, err
// 	}
// 	json.Unmarshal(raw_data, &output)
// // define schema for this API

// // convert this data into the proper type

// 	return output, nil

// }
