package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/sq9nje/ktwair/internal/ktwair"
	"github.com/sq9nje/ktwair/internal/shared"
	"github.com/sq9nje/ktwair/pkg/logging"
)

func main() {

	// Get config from file
	_, err := shared.ReadConfig("/home/ux86dp/ktwair/config.json")
	if err != nil {
		log.Fatalf("Could not read config: %v\n", err)
	}

	logging.SetLevelFromString(shared.GlobalConfig.LogLevel)

	lastTimestamp := time.Time{}
	loc, _ := time.LoadLocation("Europe/Warsaw")

	// Infinite loop
	for {
		stationJSON, err := ktwair.GetStationData(shared.GlobalConfig.KTWAir.StationID, lastTimestamp)
		if err != nil {
			logging.Logf(logging.ERROR, "Failed to fetch station data: %v", err)
			continue
		}

		stationData := ktwair.Station{}
		err = json.Unmarshal(stationJSON, &stationData)
		if err != nil {
			logging.Logf(logging.ERROR, "Could not unmarshall station data: %v", err)
			continue
		}

		// Number of returned datapoints can be 0 depending on the startTime
		if len(stationData.Sensors) > 0 {
			lastTimestamp, err = time.ParseInLocation("2006-01-02 15:04:05", stationData.Sensors[0].Data[len(stationData.Sensors[0].Data)-1].Timestamp, loc)
			if err != nil {
				logging.Logf(logging.ERROR, "Parsing last timestamp failed: %v", err)
			} else {
				logging.Logf(logging.INFO, "Number of records %d Last timestamp: %v", len(stationData.Sensors[0].Data), lastTimestamp)
			}

			ktwair.LogLatest(&stationData)

		} else {
			logging.Logf(logging.DEBUG, "No sensor data returned")
		}

		time.Sleep(time.Duration(shared.GlobalConfig.KTWAir.QueryInterval) * time.Second)
	}
}
