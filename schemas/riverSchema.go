package river

type Schema struct {
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
