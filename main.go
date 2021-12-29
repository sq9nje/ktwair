package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	u "net/url"
	"strconv"
	"time"
)

type station struct {
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	Address string   `json:"address"`
	Lat     float64  `json:"lat"`
	Lon     float64  `json:"lon"`
	Sensors []sensor `json:"sensors"`
}

type sensor struct {
	Name string        `json:"name"`
	Unit string        `json:"unit"`
	Data []measurement `json:"data"`
}

type measurement struct {
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

	log.Printf("INFO: %s %s", url, resp.Status)

	body, _ := ioutil.ReadAll(resp.Body)
	return body, nil
}

func printLatest(stationData *station) {
	loc, _ := time.LoadLocation("Europe/Warsaw")

	for _, sens := range stationData.Sensors {
		fmt.Printf("\t - %s\n", sens.Name)
		timestamp, _ := time.ParseInLocation("2006-01-02 15:04:05", sens.Data[len(sens.Data)-1].Timestamp, loc)
		value, _ := strconv.ParseFloat(sens.Data[len(sens.Data)-1].Value, 64)
		fmt.Printf("\t\t%v\t%f %s\n", timestamp, value, sens.Unit)
	}
}

func main() {

	stationID := 80
	interval := 60

	ticker := time.NewTicker(time.Duration(interval) * time.Second)

	lastTimestamp := time.Time{}
	loc, _ := time.LoadLocation("Europe/Warsaw")

	go func() {
		for range ticker.C {
			stationJSON, err := getStationData(stationID, lastTimestamp)
			if err != nil {
				log.Printf("ERROR: %v\n", err)
			} else {
				stationData := station{}
				json.Unmarshal(stationJSON, &stationData)
				lastTimestamp, _ = time.ParseInLocation("2006-01-02 15:04:05", stationData.Sensors[0].Data[len(stationData.Sensors[0].Data)-1].Timestamp, loc)

				printLatest(&stationData)
			}
		}
	}()

	for true {
	}

}
