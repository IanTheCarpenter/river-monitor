package external_apis

type Mainstem struct {
	Properties Properties `json:"properties"`
}

type Properties struct {
	Name     string     `json:"name_at_outlet"`
	Datasets []Location `json:"datasets"`
	OutletID int64      `json:"outlet_nhdpv1_comid"`
}

type Location struct {
	Name            string `json:"siteName"`
	Type            string `json:"type"`
	MeasurementType string `json:"variableMeasured"`
	Coords          string `json:"wkt"`
}

type RawNavigationData struct {
	Features []RawNavigationDataSite `json:"features"`
}

type RawNavigationDataSite struct {
	Geometry   RawNavigationDataCoordinates `json:"geometry"`
	Properties RawProperties                `json:"properties"`
}

type RawNavigationDataCoordinates struct {
	Coordinates string `json:"coordinates"`
}

type RawProperties struct {
	Name string `json:"name"`
	Id   string `json:"identifier"`
}
