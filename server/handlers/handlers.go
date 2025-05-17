package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/IanTheCarpenter/river-monitor/db"
	"github.com/IanTheCarpenter/river-monitor/schemas"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func GetRoot(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("Test"))
}

func GetFullRiverReport(writer http.ResponseWriter, request *http.Request) {
	fmt.Println()
	filter := bson.D{{Key: "river_name", Value: request.PathValue("river_name")}}

	var result schemas.Forecast
	err := db.RIVER_REPORTS.FindOne(context.TODO(), filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		fmt.Println("404")
		http.NotFound(writer, request)
		return
	}

	fmt.Printf("get river: %s\n", result.River)

	result_json, err := json.Marshal(result)
	if err != nil {
		fmt.Println("marshal error")
		fmt.Println(err)
		writer.WriteHeader(501)
		json.NewEncoder(writer).Encode(result)
		return
	}

	fmt.Printf("json: %s", result_json)

	writer.Header().Set("Content-Type", "application/json")
	writer.Write(result_json)

}
