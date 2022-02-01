package main

import (
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

	// Infinite loop
	for {
		stationData := ktwair.Station{}
		err = stationData.Fetch(shared.GlobalConfig.KTWAir.StationID, lastTimestamp)
		if err != nil {
			logging.Logf(logging.ERROR, "Could not fetch station data: %v", err)
			continue
		}

		// Number of returned datapoints can be 0 depending on the startTime
		if stationData.NumSensors() > 0 {
			lastTimestamp, err = stationData.Sensors[0].LastTimestamp()
			if err != nil {
				logging.Logf(logging.ERROR, "Parsing last timestamp failed: %v", err)
			}

			logging.Logf(logging.INFO, "Number of records %d Last timestamp: %v", stationData.Sensors[0].NumData(), lastTimestamp)
			stationData.LogLatest()

		} else {
			logging.Logf(logging.DEBUG, "No sensor data returned")
		}

		time.Sleep(time.Duration(shared.GlobalConfig.KTWAir.QueryInterval) * time.Second)
	}
}
