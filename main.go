package main

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("UNABLE TO LOAD ENV FILE")
		return
	}

	connection_string := os.Getenv("CONNECTION_STRING")

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	connection_options := options.Client().ApplyURI(connection_string).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(connection_options)
	if err != nil {
		panic(err)
	}

	collection := client.Database("RIVER_MONITOR").Collection("RIVER_DEFINITIONS")

	filter := bson.D{{"river_name", "colorado"}}

	// var bson_result bson.D
	var result RiverSchema

	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		panic(err)
	}

	fmt.Println(result.RiverName)
	fmt.Println(result.DataCollectionSites[0].Latitude)

}

type RiverSchema struct {
	RiverName           string               `bson:"river_name"`
	DataCollectionSites []DataCollectionSite `bson:"data_collection_sites"`
}

type DataCollectionSite struct {
	URL        string     `bson:"url"`
	Latitude   float64    `bson:"latitude"`
	Longitude  float64    `bson:"longitude"`
	Thresholds Thresholds `bson:"thresholds"`
}

type Thresholds struct {
	High int64 `bson:"high"`
	Low  int64 `bson:"low"`
}
