package main

import (
	"flag"
	"log"
	"time"

	"github.com/davecgh/go-spew/spew"
)

const agrometURL = "http://agromet.inia.cl/estaciones.php"

func main() {
	var stationName, from, to string
	flag.StringVar(&stationName, "station", "", "station code name")
	flag.StringVar(&from, "from", "", "initial date (dd-MM-yyyy)")
	flag.StringVar(&to, "to", "", "final date (dd-MM-yyyy)")
	flag.Parse()

	fromDate, err := time.Parse("02-01-2006", from)
	if err != nil {
		log.Fatal(err)
	}
	toDate, err := time.Parse("02-01-2006", to)
	if err != nil {
		log.Fatal(err)
	}
	data, err := stationData(stationName, fromDate, toDate)
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	values, err := extractMeasurements(data)
	if err != nil {
		log.Fatal(err)
	}

	jsonName := stationName + ".json"
	currentMsr, err := loadMeasurements(jsonName)
	if err != nil {
		log.Fatal(err)
	}
	mergedValues := mergeMeasurements(currentMsr, values)
	spew.Dump(mergedValues)
	err = saveMeasurements(jsonName, mergedValues)
	if err != nil {
		log.Fatal(err)
	}
}
