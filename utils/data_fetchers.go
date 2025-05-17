package utils

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
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func Fetch_river_definitions() ([]schemas.River, error) {

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

func insertDataInDescendingOrder(data []SiteSample, sampletoinsert SiteSample) []SiteSample {
	if len(data) == 0 || sampletoinsert.TimeStamp.After(data[0].TimeStamp) {
		data = append(data, sampletoinsert)
	} else {
		data = append([]SiteSample{sampletoinsert}, data...)
	}
	return data
}
