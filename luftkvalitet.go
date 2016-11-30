package luftkvalitet

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var endpoint = "https://api.nilu.no/"

type Area struct {
	Zone         string `json:"zone"`
	Municipality string `json:"municipality"`
	Area         string `json:"area"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Station struct {
	Area
	Location
	Station string `json:"station"`
}

type Measurement struct {
	Station
	Eoi       string    `json:"eoi"`
	Component string    `json:"component"`
	FromTime  time.Time `json:"fromTime"`
	ToTime    time.Time `json:"toTime"`
	Value     float64   `json:"value"`
	Unit      string    `json:"unit"`
	Index     int       `json:"index"`
	Color     string    `json:"color"`
}

type Point struct {
	Location
	Radius float64
}

type Filter struct {
	Areas      []string
	Stations   []string
	Components []string
	Within     Point
	Nearest    Point
}

type AqiResult struct {
	Component string `json:"component"`
	Unit      string `json:"unit"`
	Aqis      []Aqi  `json:"aqis"`
}

type Aqi struct {
	Index            int     `json:"index"`
	FromValue        float64 `json:"fromValue"`
	ToValue          float64 `json:"toValue"`
	Color            string  `json:"color"`
	Text             string  `json:"text"`
	ShortDescription string  `json:"shortDescription"`
	Description      string  `json:"description"`
}

func GetMeasurements(f Filter) ([]Measurement, error) {
	u := endpoint + url.QueryEscape("aq/utd.json?")

	if len(f.Areas) > 0 {
		query := url.QueryEscape("areas=" + strings.Join(f.Areas, ";"))
		u = u + query
	}

	if len(f.Stations) > 0 {
		query := url.QueryEscape("&stations=" + strings.Join(f.Stations, ";"))
		u = u + query
	}

	if len(f.Components) > 0 {
		query := url.QueryEscape("&components=" + strings.Join(f.Components,
			";"))
		u = u + query
	}

	if f.Within.Latitude != 0 && f.Within.Longitude != 0 && f.Within.Radius != 0 {
		lat := strconv.FormatFloat(f.Within.Latitude, 'E', -1, 64)
		long := strconv.FormatFloat(f.Within.Longitude, 'E', -1, 64)
		radius := strconv.FormatFloat(f.Within.Radius, 'E', -1, 64)

		query := url.QueryEscape("&within=" + strings.Join([]string{lat, long, radius}, ";"))
		u = u + query
	}

	if f.Nearest.Latitude != 0 && f.Nearest.Longitude != 0 && f.Nearest.Radius != 0 {
		lat := strconv.FormatFloat(f.Nearest.Latitude, 'E', -1, 64)
		long := strconv.FormatFloat(f.Nearest.Longitude, 'E', -1, 64)
		radius := strconv.FormatFloat(f.Nearest.Radius, 'E', -1, 64)

		query := url.QueryEscape("&nearest=" + strings.Join([]string{lat, long, radius}, ";"))
		u = u + query
	}

	resp, err := http.Get(u)

	if err != nil {
		return []Measurement{}, err
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var measurements []Measurement
	err = json.Unmarshal(body, &measurements)
	if err != nil {
		return []Measurement{}, err
	}

	return measurements, nil

}

func GetAreas() ([]Area, error) {
	u := endpoint + "lookup/areas"
	resp, err := http.Get(u)

	if err != nil {
		return []Area{}, err
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var areas []Area
	err = json.Unmarshal(body, &areas)
	if err != nil {
		return []Area{}, err
	}

	return areas, nil

}

func GetStations() ([]Station, error) {
	u := endpoint + "lookup/stations"
	resp, err := http.Get(u)

	if err != nil {
		return []Station{}, err
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var stations []Station
	err = json.Unmarshal(body, &stations)
	if err != nil {
		return []Station{}, err
	}

	return stations, nil
}

func GetAqis() ([]AqiResult, error) {
	u := endpoint + "lookup/aqis"
	resp, err := http.Get(u)

	if err != nil {
		return []AqiResult{}, err
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var results []AqiResult
	err = json.Unmarshal(body, &results)
	if err != nil {
		return []AqiResult{}, err
	}

	return results, nil

}
