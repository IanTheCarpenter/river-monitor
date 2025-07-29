package external_apis

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// i do not like this sam i am
type usgs_data_sample struct {
	Value    string `json:"value"`
	DateTime string `json:"dateTime"`
}

type usgs_values2 struct {
	Value []usgs_data_sample `json:"value"`
}

type usgs_source_info struct {
	Name string `json:"siteName"`
}

type usgs_variable struct {
	VariableName string `json:"variableName"`
}

type usgs_data_streams struct {
	SourceInfo usgs_source_info `json:"sourceInfo"`
	Variable   usgs_variable    `json:"variable"`
	Values     []usgs_values2   `json:"values"`
}

// type usgs

type usgs_value1 struct {
	TimeSeries []usgs_data_streams `json:"timeSeries"`
}

type usgs_data struct {
	Data usgs_value1 `json:"value"`
}

func usgs_convert_to_siteSamples(input []usgs_data_sample, url string) []SiteSample {
	var output []SiteSample

	for _, sample := range input {
		var current_record SiteSample
		timestamp, err := time.Parse(time.RFC3339, sample.DateTime)
		if err != nil {
			fmt.Printf("Could not parse timestamp for data from: %s", url)
			continue
		}
		float, err := strconv.ParseFloat(sample.Value, 64)
		if err != nil {
			fmt.Printf("Could not parse stage into string for data from: %s", url)
			continue
		}

		current_record.TimeStamp = timestamp
		current_record.Value = float

		output = InsertDataInDescendingOrder(output, current_record)
	}
	return output
}

func USGS_FetchSite(site string, days int) (SiteData, error) {

	// format URL
	url := fmt.Sprintf("https://waterservices.usgs.gov/nwis/iv/?sites=%s&agencyCd=USGS&period=P%dD&format=json", site, days)

	raw_data, err := Fetch_data(url)
	if err != nil {
		return SiteData{}, err
	}

	var usgs usgs_data
	json.Unmarshal(raw_data, &usgs)

	// convert this data into type []external_apis.SiteSample
	if len(usgs.Data.TimeSeries) < 1 {
		fmt.Printf("No data returned for site: %s\n", site)
		return SiteData{}, nil
	}
	fmt.Println()

	var output = SiteData{
		SiteName: usgs.Data.TimeSeries[0].SourceInfo.Name,
	}

	//initialize values for site
	for _, dataset := range usgs.Data.TimeSeries {
		if strings.HasPrefix(dataset.Variable.VariableName, "Streamflow") {
			output.Flow = usgs_convert_to_siteSamples(dataset.Values[0].Value, url)
		} else if strings.HasPrefix(dataset.Variable.VariableName, "Gage height") {
			output.Stage = usgs_convert_to_siteSamples(dataset.Values[0].Value, url)

		} else {
			continue
		}
	}
	return output, nil

}

func USGS_FetchLimit(site string) (float64, error) {
	type usgs_stages struct {
		Action float64 `json:"action"`
	}

	type usgs_stages_response struct {
		Stages usgs_stages `json:"stages"`
	}

	var output float64

	// format URL
	url := fmt.Sprintf("https://waterdata.usgs.gov/flood-stage/%s", site)

	raw_data, err := Fetch_data(url)
	if err != nil {
		return output, err
	}
	var empty usgs_stages_response
	var stages usgs_stages_response
	json.Unmarshal(raw_data, &stages)

	if stages == empty {
		return 0, errors.New("ERROR: No limit returned by USGS API")
	}

	fmt.Printf("Limit: %f\n", stages.Stages.Action)

	return stages.Stages.Action, err

}
