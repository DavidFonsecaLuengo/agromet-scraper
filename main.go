package main

import (
	"flag"
	"log"
	"time"
)

const (
	agrometURL   = "http://agromet.inia.cl/estaciones.php"
	stationsJSON = "stations.json"
)

func main() {
	var updStations, all bool
	flag.BoolVar(&updStations, "upd", false, "update stations.json")
	flag.BoolVar(&all, "all", false, "fetch data of all stations")

	var stationID, from, to string
	flag.StringVar(&stationID, "station", "", "station code name")
	flag.StringVar(&from, "from", "", "initial date (dd-MM-yyyy)")
	flag.StringVar(&to, "to", "", "final date (dd-MM-yyyy)")

	var waitTime time.Duration
	flag.DurationVar(&waitTime, "wt", time.Minute, "wait time between requests (to avoid flooding)")

	flag.Parse()

	if updStations {
		log.Printf("updating station list...")
		err := updateStations(stationsJSON)
		if err != nil {
			log.Fatal(err)
		}
	}

	// no measurement fetching, finishing early
	if !all && stationID == "" && from == "" && to == "" {
		return
	}

	fromDate, err := time.Parse("02-01-2006", from)
	if err != nil {
		log.Fatal(err)
	}
	toDate, err := time.Parse("02-01-2006", to)
	if err != nil {
		log.Fatal(err)
	}

	stations, err := loadStations(stationsJSON)
	if err != nil {
		log.Fatal(err)
	}

	switch {
	case stationID != "":
		for _, station := range stations {
			if stationID == station.ID {
				err := fetchStationMeasurements(stationID, fromDate, toDate)
				if err != nil {
					log.Fatal(err)
				}
				return
			}
		}
		log.Printf("station %s not found in %s", stationID, stationsJSON)
	case all:
		for _, station := range stations {
			log.Printf("fetching measurements of %s (%s)", station.Name, station.ID)
			err := fetchStationMeasurements(station.ID, fromDate, toDate)
			if err != nil {
				log.Printf("error fetching measurements of %s: %s", station.ID, err)
			}
			if waitTime > 0 {
				log.Printf("waiting %v...", waitTime)
				time.Sleep(waitTime)
			}
		}
	}
}

func fetchStationMeasurements(stationID string, fromDate, toDate time.Time) error {
	data, err := stationData(stationID, fromDate, toDate)
	if err != nil {
		return err
	}
	defer data.Close()
	values, err := extractMeasurements(data)
	if err != nil {
		return err
	}

	jsonName := stationID + ".json"
	currentMsr, err := loadMeasurements(jsonName)
	if err != nil {
		return err
	}
	mergedValues := mergeMeasurements(currentMsr, values)
	return saveMeasurements(jsonName, mergedValues)
}
