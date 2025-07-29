package external_apis

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/IanTheCarpenter/river-monitor/db"
	"github.com/IanTheCarpenter/river-monitor/schemas"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type SiteData struct {
	SiteName       string
	StageThreshold float64
	Flow           []SiteSample
	Stage          []SiteSample
}
type SiteSample struct {
	TimeStamp time.Time
	Value     float64
}

func Fetch_river_definitions() ([]schemas.River, error) {
	// var bson_result bson.D
	// var result river.Schema

	bson_response, err := db.RIVER_DEFINITIONS.Find(context.TODO(), bson.D{{}}, options.Find())
	if err != nil {
		fmt.Println("ERROR FETCHING RIVER DEFINITIONS")
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

func Fetch_data(url string) ([]byte, error) {
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

func InsertDataInDescendingOrder(data []SiteSample, sampletoinsert SiteSample) []SiteSample {
	if len(data) == 0 || sampletoinsert.TimeStamp.Before(data[0].TimeStamp) {
		data = append(data, sampletoinsert)
	} else {
		data = append([]SiteSample{sampletoinsert}, data...)
	}
	return data
}
