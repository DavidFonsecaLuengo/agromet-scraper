package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

type Station struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

func loadStations(file string) ([]Station, error) {
	var currentStations []Station
	j, err := os.Open(file)
	if err != nil {
		if os.IsNotExist(err) {
			return currentStations, nil
		}
		return nil, err
	}
	defer j.Close()
	err = json.NewDecoder(j).Decode(&currentStations)
	if err != nil {
		return nil, err
	}
	return currentStations, nil

}

func updateStations(file string) error {
	resp, err := http.Get(agrometURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	stations, err := extractStations(resp.Body)
	if err != nil {
		return err
	}

	// save to file
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(stations)
	if err != nil {
		return err
	}
	return nil
}

func extractStations(r io.Reader) ([]Station, error) {
	var buf bytes.Buffer
	tee := io.TeeReader(r, &buf)

	// get stations (id and name)
	doc, err := goquery.NewDocumentFromReader(tee)
	if err != nil {
		return nil, nil
	}
	var stations []Station
	stationIdx := make(map[string]int)
	doc.Find("#estaciones > optgroup > option").Each(func(_ int, node *goquery.Selection) {
		value, ok := node.Attr("value")
		if !ok {
			log.Printf("value not found in station option %s", node.Text())
			return
		}
		label, ok := node.Attr("label")
		if !ok {
			log.Printf("label not found in station option %s", node.Text())
			return
		}
		stations = append(stations, Station{
			ID:   value,
			Name: label,
		})
		stationIdx[value] = len(stations) - 1
	})

	// get coordinates of each station
	pointRegexp, err := regexp.Compile("var point = new google\\.maps\\.LatLng\\((?P<lng>.+), (?P<lat>.+)\\);")
	if err != nil {
		panic(err)
	}
	onclickRegexp, err := regexp.Compile("onclick=seleccionarEma\\('(?P<stationId>.+)'\\)")
	if err != nil {
		panic(err)
	}
	var lng, lat float64
	scanner := bufio.NewScanner(&buf)
	for scanner.Scan() {
		pointMatches := pointRegexp.FindStringSubmatch(scanner.Text())
		if pointMatches != nil {
			lat, err = strconv.ParseFloat(pointMatches[1], 64)
			if err != nil {
				return nil, err
			}
			lng, err = strconv.ParseFloat(pointMatches[2], 64)
			if err != nil {
				return nil, err
			}
		}
		onclickMatches := onclickRegexp.FindStringSubmatch(scanner.Text())
		if onclickMatches != nil {
			idx := stationIdx[onclickMatches[1]]
			stations[idx].Longitude = lng
			stations[idx].Latitude = lat
		}
	}
	if scanner.Err() != nil {
		return nil, scanner.Err()
	}

	return stations, nil
}
