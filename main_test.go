package main

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExtractMeasurements(t *testing.T) {
	f, err := os.Open("Agromet _ Estaciones.html")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	v, err := extractMeasurements(f)
	if err != nil {
		t.Fatal(err)
	}

	assert.Len(t, v, 36)
	// first row
	assert.Equal(t, time.Date(2018, 1, 2, 0, 0, 0, 0, time.UTC), v[0].Date)
	assert.EqualValues(t, 10.5, v[0].Temperature)
	assert.EqualValues(t, 100, v[0].Humidity)
	assert.EqualValues(t, 18.4, v[0].WindVelocity)
	// last row
	assert.Equal(t, time.Date(2018, 1, 3, 11, 0, 0, 0, time.UTC), v[35].Date)
	assert.EqualValues(t, 17.4, v[35].Temperature)
	assert.EqualValues(t, 82, v[35].Humidity)
	assert.EqualValues(t, 10.4, v[35].WindVelocity)
}

func TestMergeMeasurementsNilM1(t *testing.T) {
	now := time.Now()
	m := mergeMeasurements(nil, []Measurements{{
		Date:         now,
		Temperature:  1.0,
		Humidity:     2.0,
		WindVelocity: 3.0,
	}})
	assert.Len(t, m, 1)
	assert.Equal(t, now, m[0].Date)
	assert.Equal(t, 1.0, m[0].Temperature)
	assert.Equal(t, 2.0, m[0].Humidity)
	assert.Equal(t, 3.0, m[0].WindVelocity)
}

func TestMergeMeasurements(t *testing.T) {
	m1 := []Measurements{{
		Date:         time.Date(2018, 1, 4, 21, 0, 0, 0, time.UTC),
		Temperature:  1.0,
		Humidity:     2.0,
		WindVelocity: 3.0,
	}}
	m2 := []Measurements{{
		Date:         time.Date(2018, 1, 4, 23, 0, 0, 0, time.UTC),
		Temperature:  1.0,
		Humidity:     2.0,
		WindVelocity: 3.0,
	}, {
		Date:         time.Date(2018, 1, 4, 21, 0, 0, 0, time.UTC),
		Temperature:  1.0,
		Humidity:     2.0,
		WindVelocity: 3.0,
	}, {
		Date:         time.Date(2018, 1, 4, 22, 0, 0, 0, time.UTC),
		Temperature:  1.0,
		Humidity:     2.0,
		WindVelocity: 3.0,
	}}
	m := mergeMeasurements(m1, m2)
	assert.Len(t, m, 3)
	assert.Equal(t, time.Date(2018, 1, 4, 21, 0, 0, 0, time.UTC), m[0].Date)
	assert.Equal(t, 1.0, m[0].Temperature)
	assert.Equal(t, 2.0, m[0].Humidity)
	assert.Equal(t, 3.0, m[0].WindVelocity)
	assert.Equal(t, time.Date(2018, 1, 4, 22, 0, 0, 0, time.UTC), m[1].Date)
	assert.Equal(t, 1.0, m[1].Temperature)
	assert.Equal(t, 2.0, m[1].Humidity)
	assert.Equal(t, 3.0, m[1].WindVelocity)
	assert.Equal(t, time.Date(2018, 1, 4, 23, 0, 0, 0, time.UTC), m[2].Date)
	assert.Equal(t, 1.0, m[2].Temperature)
	assert.Equal(t, 2.0, m[2].Humidity)
	assert.Equal(t, 3.0, m[2].WindVelocity)
}
