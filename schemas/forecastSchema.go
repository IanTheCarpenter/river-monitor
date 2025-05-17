package schemas

import "go.mongodb.org/mongo-driver/v2/bson"

type Forecast struct {
	River                     string                    `bson:"river_name"`
	RiverObjectID             bson.ObjectID             `bson:"river_object_id"`
	LastUpdated               string                    `bson:"last_updated"`
	PointsOfInterestForecasts []PointOfInterestForecast `bson:"points_of_interest"`
}

type PointOfInterestForecast struct {
	Name     string `bson:"name" `
	Runnable bool   `bson:"runnable"`
}
