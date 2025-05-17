package db

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var DB *mongo.Client
var RIVER_DEFINITIONS *mongo.Collection
var RIVER_REPORTS *mongo.Collection
var USER_DATA *mongo.Collection

func Init() {
	connection_string := os.Getenv("CONNECTION_STRING")
	if connection_string == "" {
		fmt.Println("Loading env_vars from file...")
		err := godotenv.Load()
		if err != nil {
			fmt.Println("UNABLE TO LOAD ENV FILE")
			panic(err)
		}
		connection_string = os.Getenv("CONNECTION_STRING")
		fmt.Println("...Done!")

	} else {
		fmt.Println("Env vars successfully loaded environment variables")
	}
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	connection_options := options.Client().ApplyURI(connection_string).SetServerAPIOptions(serverAPI)

	DB, err := mongo.Connect(connection_options)
	if err != nil {
		panic(err)
	}
	RIVER_DEFINITIONS = DB.Database("RIVER_MONITOR").Collection("RIVER_DEFINITIONS")
	RIVER_REPORTS = DB.Database("RIVER_MONITOR").Collection("RIVER_REPORTS")

}
