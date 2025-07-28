package schemas

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type River struct {
	RiverName           string               `bson:"river_name"`
	ObjectID            bson.ObjectID        `bson:"_id,omitempty"`
	DataCollectionSites []DataCollectionSite `bson:"data_collection_sites"`
	PointsOfInterest    []PointOfInterest    `bson:"points_of_interest"`
}

type DataCollectionSite struct {
	Name      string          `bson:"name"`
	SiteID    string          `bson:"site_id"`
	Agency    string          `bson:"agency"`
	Latitude  float64         `bson:"latitude"`
	Longitude float64         `bson:"longitude"`
	HasFlow   bool            `bson:"hasflow"`
	Stage     StageThresholds `bson:"stage"`
	// URL        string     `bson:"url"`
}

type StageThresholds struct {
	Limit            float64   `bson:"high"`
	Baseline         float64   `bson:"baseline"`
	Samples          int       `bson:"samples"`
	MostRecentSample time.Time `bson:"most_recent_sample_averaged"`
}

type PointOfInterest struct {
	Name       string   `bson:"name"`
	Type       string   `bson:"type"`
	Latitude   float64  `bson:"latitude"`
	Longitude  float64  `bson:"longitude"`
	Indicators []string `bson:"indicators"`
}
