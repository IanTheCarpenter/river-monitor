package river_data

import (
	"github.com/IanTheCarpenter/river-monitor/forecaster/schemas"
)

var sanGabrielRiver = schemas.River{
	RiverName: "san_gabriel",
	PointsOfInterest: []schemas.PointOfInterest{
		{
			Name:       "Highway 29 Crossing",
			Type:       "Public Access",
			Latitude:   30.645658650442574,
			Longitude:  -97.58441396892775,
			Indicators: []string{"usgs-08105300"},
		},
	},
	DataCollectionSites: []schemas.DataCollectionSite{
		{
			Name:      "usgs-08105300",
			Agency:    "usgs",
			URL:       "https://waterservices.usgs.gov/nwis/iv/?sites=08105300&agencyCd=USGS&period=P1D&parameterCd=00065&format=json",
			Latitude:  30.646,
			Longitude: -97.5852,
			Thresholds: schemas.Thresholds{
				High: 12,
				Low:  7.5,
			},
		},
	},
}
