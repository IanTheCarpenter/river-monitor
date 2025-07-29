package external_apis

// import (
// 	"encoding/json"
// 	"fmt"
// 	"time"
// )

// func Fetch_lcra_data(url string) (SiteData, error) {
// 	var output SiteData
// 	raw_data, err := Fetch_data(url)
// 	if err != nil {
// 		return output, err
// 	}
// 	// define schema for this API
// 	type lcra_records struct {
// 		DateTime string  `json:"dateTime"`
// 		Stage    float64 `json:"value1"`
// 		Flow     float64 `json:"value2"`
// 	}

// 	type lcra_data struct {
// 		SiteName string         `json:"siteName"`
// 		Records  []lcra_records `json:"records"`
// 	}

// 	// convert this data into type []SiteSample
// 	var lcra lcra_data
// 	json.Unmarshal(raw_data, &lcra)

// 	for _, sample := range lcra.Records {
// 		var current_record SiteSample

// 		timestamp, err := time.Parse(time.RFC3339, sample.DateTime)
// 		if err != nil {
// 			fmt.Printf("Could not parse timestamp for data from: %s", url)
// 			continue
// 		}

// 		current_record.TimeStamp = timestamp
// 		current_record.Flow = sample.Flow
// 		current_record.Stage = sample.Stage

// 		output = InsertDataInDescendingOrder(output, current_record)
// 	}
// 	return output, nil
// }
