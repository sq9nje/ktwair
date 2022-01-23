package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	u "net/url"
	"strconv"
	"time"

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

func getStationData(stationID int, startTime time.Time) ([]byte, error) {

	var baseURL string = "https://powietrze.katowice.eu/data/station/"

	url := baseURL + strconv.FormatInt(int64(stationID), 10)
	if !startTime.IsZero() {
		url += "?from=" + u.QueryEscape(startTime.Format("2006-01-02 15:04:05"))
	}

	httpClient := http.Client{Timeout: time.Duration(30) * time.Second}
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	logging.Logf(logging.INFO, "%s %s", url, resp.Status)

	body, _ := ioutil.ReadAll(resp.Body)
	return body, nil
}

func logLatest(stationData *Station) {
	loc, _ := time.LoadLocation("Europe/Warsaw")
	for _, s := range stationData.Sensors {
		timestamp, _ := time.ParseInLocation("2006-01-02 15:04:05", s.Data[len(s.Data)-1].Timestamp, loc)
		value, _ := strconv.ParseFloat(s.Data[len(s.Data)-1].Value, 64)
		logging.Logf(logging.DEBUG, "Last measurement %s %v %f %s", s.Name, timestamp, value, s.Unit)
	}
}

func main() {

	logging.SetLevelFromString("DEBUG")

	stationID := 80
	interval := 60

	lastTimestamp := time.Time{}
	loc, _ := time.LoadLocation("Europe/Warsaw")

	for {
		stationJSON, err := getStationData(stationID, lastTimestamp)
		if err != nil {
			logging.Logf(logging.ERROR, "%v", err)
		} else {
			stationData := Station{}
			json.Unmarshal(stationJSON, &stationData)
			// Sort sensor data by timestamp
			// for _, s := range stationData.Sensors {
			// 	sort.Slice(s.Data[:], func(i, j int) bool {
			// 		return s.Data[i].Timestamp < s.Data[j].Timestamp
			// 	})

			// }
			lastTimestamp, err = time.ParseInLocation("2006-01-02 15:04:05", stationData.Sensors[0].Data[len(stationData.Sensors[0].Data)-1].Timestamp, loc)
			if err != nil {
				logging.Logf(logging.ERROR, "Parsing last timestamp failed: %v", err)
			} else {
				logging.Logf(logging.INFO, "Last timestamp: %v", lastTimestamp)
			}

			logLatest(&stationData)
		}

		time.Sleep(time.Duration(interval) * time.Second)
	}
}
