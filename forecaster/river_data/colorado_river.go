package river_data

import (
	"github.com/IanTheCarpenter/river-monitor/forecaster/schemas"
)

var coloradoRiver = schemas.River{
	RiverName: "colorado",
	PointsOfInterest: []schemas.PointOfInterest{
		{
			Name:       "Little Webberville Park",
			Type:       "Public Access",
			Latitude:   30.22946835858506,
			Longitude:  -97.51890756915127,
			Indicators: []string{"lcra_5423"},
		},
	},
	DataCollectionSites: []schemas.DataCollectionSite{
		{
			Name:      "lcra_5423",
			Agency:    "lcra",
			URL:       "https://hydromet.lcra.org/api/GetDataBySite/5423/flow",
			Latitude:  30.22946835858506,
			Longitude: -97.51890756915127,
			Thresholds: schemas.Thresholds{
				High: 12,
				Low:  7.25,
			},
		},
	},
}
