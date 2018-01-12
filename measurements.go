package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func loadMeasurements(file string) ([]Measurements, error) {
	var currentValues []Measurements
	j, err := os.Open(file)
	if err != nil {
		if os.IsNotExist(err) {
			return currentValues, nil
		}
		return nil, err
	}
	defer j.Close()
	err = json.NewDecoder(j).Decode(&currentValues)
	if err != nil {
		return nil, err
	}
	return currentValues, nil
}

func saveMeasurements(file string, values []Measurements) error {
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(values)
	if err != nil {
		return err
	}
	return nil
}

func stationData(station string, from, to time.Time) (io.ReadCloser, error) {
	values := make(url.Values)
	values.Set("estaciones[]", station)
	values.Add("variables[]", "2007-variable")
	values.Add("variables[]", "2002-variable")
	values.Add("variables[]", "2003-variable")
	values.Add("intervalos", "hour")
	values.Add("desde", from.Format("02-01-2006"))
	values.Add("hasta", to.Format("02-01-2006"))
	values.Add("desde_meses", "01")
	values.Add("desde_anos", "2017")
	values.Add("hasta_meses", "12")
	values.Add("hasta_anos", "2017")
	values.Add("desde_anos_a", "2017")
	values.Add("hasta_anos_a", "2017")
	values.Add("html", "html")
	resp, err := http.PostForm(agrometURL, values)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

type Measurements struct {
	Date         time.Time `json:"date"`
	Temperature  float64   `json:"temperature"`
	Humidity     float64   `json:"humidity"`
	WindVelocity float64   `json:"wind_velocity"`
}

func extractMeasurements(r io.Reader) ([]Measurements, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, nil
	}

	var values []Measurements
	i := 0
	var currentValue Measurements
	doc.Find(".table.table-striped.table-bordered > tbody > tr > td").Each(func(_ int, node *goquery.Selection) {
		switch i % 5 {
		case 1:
			date, err := time.Parse("02-01-2006 15:04", node.Text())
			if err != nil {
				log.Print(err)
			}
			currentValue.Date = date.UTC()
		case 2:
			fixedSz := strings.Replace(node.Text(), ",", ".", 1)
			temp, err := strconv.ParseFloat(fixedSz, 64)
			if err != nil {
				log.Print(err)
			}
			currentValue.Temperature = temp
		case 3:
			fixedSz := strings.Replace(node.Text(), ",", ".", 1)
			hum, err := strconv.ParseFloat(fixedSz, 64)
			if err != nil {
				log.Print(err)
			}
			currentValue.Humidity = hum
		case 4:
			fixedSz := strings.Replace(node.Text(), ",", ".", 1)
			windVel, err := strconv.ParseFloat(fixedSz, 64)
			if err != nil {
				log.Print(err)
			}
			currentValue.WindVelocity = windVel
			values = append(values, currentValue)
		}
		i++
	})
	return values, nil
}

func mergeMeasurements(m1, m2 []Measurements) []Measurements {
	merged := make([]Measurements, 0, len(m1)+len(m2))
	if len(m1) > 0 {
		merged = append(merged, m1...)
	}

	dates := make(map[time.Time]struct{})
	for _, m := range m1 {
		dates[m.Date] = struct{}{}
	}
	for _, m := range m2 {
		_, ok := dates[m.Date]
		if !ok {
			merged = append(merged, m)
		}
	}

	sort.Slice(merged, func(i, j int) bool {
		return merged[i].Date.Before(merged[j].Date)
	})
	return merged
}
