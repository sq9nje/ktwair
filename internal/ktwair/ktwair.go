package ktwair

import (
	"io/ioutil"
	"net/http"
	u "net/url"
	"strconv"
	"time"

	"github.com/sq9nje/ktwair/internal/shared"
	"github.com/sq9nje/ktwair/pkg/logging"
)

type Station struct {
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	Address string   `json:"address"`
	Lat     float64  `json:"lat"`
	Lon     float64  `json:"lon"`
	Sensors []Sensor `json:"sensors"`
}

type Sensor struct {
	Name string        `json:"name"`
	Unit string        `json:"unit"`
	Data []Measurement `json:"data"`
}

type Measurement struct {
	Timestamp  string `json:"timestamp"`
	Value      string `json:"value"`
	StatusCode int    `json:"status_code"`
}

// Get the measurement data from the specified station as a JSON
func GetStationData(stationID int, startTime time.Time) ([]byte, error) {

	url := shared.GlobalConfig.KTWAir.BaseURL + strconv.FormatInt(int64(stationID), 10)
	if !startTime.IsZero() {
		url += "?from=" + u.QueryEscape(startTime.Add(1*time.Second).Format("2006-01-02 15:04:05"))
	}

	httpClient := http.Client{Timeout: time.Duration(shared.GlobalConfig.KTWAir.Timeout) * time.Second}
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	logging.Logf(logging.INFO, "%s %s %s", resp.Request.Method, url, resp.Status)

	body, _ := ioutil.ReadAll(resp.Body)
	return body, nil
}

// Translate sensor names to English
func SensorNameToEN(name string) string {
	switch name {
	case "temperatura":
		return "temperature"
	case "ciśnienie":
		return "pressure"
	case "wilgotność":
		return "humidity"
	}
	return name
}

// Print latest measurement values in the log
func LogLatest(stationData *Station) {
	loc, _ := time.LoadLocation("Europe/Warsaw")
	for _, s := range stationData.Sensors {
		timestamp, err := time.ParseInLocation("2006-01-02 15:04:05", s.Data[len(s.Data)-1].Timestamp, loc)
		if err != nil {
			logging.Logf(logging.ERROR, "%s: Parsing timestamp failed: %v", s.Name, err)
			return
		}
		value, err := strconv.ParseFloat(s.Data[len(s.Data)-1].Value, 64)
		if err != nil {
			logging.Logf(logging.ERROR, "%s: Parsing value failed: %v", s.Name, err)
			return
		}
		logging.Logf(logging.DEBUG, "%s %v %f %s", SensorNameToEN(s.Name), timestamp, value, s.Unit)
	}
}
