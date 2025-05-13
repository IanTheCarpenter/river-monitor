package schemas

import "go.mongodb.org/mongo-driver/v2/bson"

type River struct {
	RiverName           string               `bson:"river_name"`
	ObjectID            bson.ObjectID        `bson:"_id,omitempty"`
	DataCollectionSites []DataCollectionSite `bson:"data_collection_sites"`
	PointsOfInterest    []PointOfInterest    `bson:"points_of_interest"`
}

type DataCollectionSite struct {
	Name       string     `bson:"name"`
	Agency     string     `bson:"agency"`
	URL        string     `bson:"url"`
	Latitude   float64    `bson:"latitude"`
	Longitude  float64    `bson:"longitude"`
	Thresholds Thresholds `bson:"thresholds"`
}

type Thresholds struct {
	High float64 `bson:"high"`
	Low  float64 `bson:"low"`
}

type PointOfInterest struct {
	Name       string   `bson:"name"`
	Type       string   `bson:"type"`
	Latitude   float64  `bson:"latitude"`
	Longitude  float64  `bson:"longitude"`
	Indicators []string `bson:"indicators"`
}
