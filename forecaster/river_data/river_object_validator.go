package river_data

import (
	"fmt"

	"github.com/IanTheCarpenter/river-monitor/forecaster/schemas"
)

// A function that reads the river data object and checks that all the points of interest
// only reference datacollection sites within the river object

func validate(r schemas.River) error {

	// generate map containing the names of the datacollectionsites to be validated later
	valid_name_map := make(map[string]bool)
	for _, data_site := range r.DataCollectionSites {
		if data_site.Name == "" {
			return fmt.Errorf("river schema validation error:\n -- datacollectionsite name cannot be empty string")
		}
		if valid_name_map[data_site.Name] {
			return fmt.Errorf("river schema validation error:\n -- duplicate datacollectionsite name: %s", data_site.Name)
		}

		valid_name_map[data_site.Name] = true
	}

	// validate references in each point of interest, checking that they exist
	for _, point := range r.PointsOfInterest {
		for _, indicator := range point.Indicators {
			if !valid_name_map[indicator] {
				return fmt.Errorf("river schema validation error:\n -- pointofinterest %s references non-existent datacollectionSite: %s", point.Name, indicator)
			}
		}
	}
	return nil
}

// func findInListByName([]schemas.DataCollectionSite, searchstring string) bool {
// 	for {}
// }
