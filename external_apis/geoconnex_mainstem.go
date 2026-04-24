package external_apis

type GeoConnexMainstem struct {
	Properties GeoConnexProperties `json:"properties"`
}

type GeoConnexProperties struct {
	Name     string              `json:"name_at_outlet"`
	Datasets []GeoConnexLocation `json:"datasets"`
	OutletID int64               `json:"outlet_nhdpv1_comid"`
}

type GeoConnexLocation struct {
	Name            string `json:"siteName"`
	Type            string `json:"type"`
	MeasurementType string `json:"variableMeasured"`
	Coords          string `json:"wkt"`
}

type GeoConnexRawNavigationData struct {
	Features []GeoConnexRawNavigationDataSite `json:"features"`
}

type GeoConnexRawNavigationDataSite struct {
	Geometry   GeoConnexRawNavigationDataCoordinates `json:"geometry"`
	Properties GeoConnexRawProperties                `json:"properties"`
}

type GeoConnexRawNavigationDataCoordinates struct {
	Coordinates string `json:"coordinates"`
}

type GeoConnexRawProperties struct {
	Name string `json:"name"`
	Id   string `json:"identifier"`
}
