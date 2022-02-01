package ktwair

import (
	"encoding/json"
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
	Lat     string   `json:"lat"`
	Lon     string   `json:"lon"`
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

// Return the number of datapoints for a Sensor
func (s *Sensor) NumData() int {
	return len(s.Data)
}

// Return the time of the last timestamp in the Sensor Data
func (s *Sensor) LastTimestamp() (time.Time, error) {
	loc, _ := time.LoadLocation("Europe/Warsaw")
	return time.ParseInLocation("2006-01-02 15:04:05", s.Data[len(s.Data)-1].Timestamp, loc)
}

// Translate sensor names to English
func (sens *Sensor) ToEN() string {
	switch sens.Name {
	case "temperatura":
		return "temperature"
	case "ciśnienie":
		return "pressure"
	case "wilgotność":
		return "humidity"
	default:
		return sens.Name
	}
}

// Get the measurement data from the specified station as a JSON
func getRawJSON(stationID int, startTime time.Time) ([]byte, error) {

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

// Fetch data for a monitoring station and unmarshal to a Station struct
func (s *Station) Fetch(stationID int, startTime time.Time) error {
	rawJSON, err := getRawJSON(stationID, startTime)
	if err != nil {
		return err
	}

	err = json.Unmarshal(rawJSON, s)
	return err
}

// Return the number of sensors for a Station
func (s *Station) NumSensors() int {
	return len(s.Sensors)
}

// Print latest measurement values in the log
func (s *Station) LogLatest() {
	for _, sens := range s.Sensors {
		timestamp, err := sens.LastTimestamp()
		if err != nil {
			logging.Logf(logging.ERROR, "%s: Parsing timestamp failed: %v", s.Name, err)
			continue
		}
		value, err := strconv.ParseFloat(sens.Data[len(sens.Data)-1].Value, 64)
		if err != nil {
			logging.Logf(logging.ERROR, "%s: Parsing value failed: %v", s.Name, err)
			continue
		}
		logging.Logf(logging.DEBUG, "%s %v %f %s", sens.ToEN(), timestamp, value, sens.Unit)
	}
}
