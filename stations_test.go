package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractStations(t *testing.T) {
	f, err := os.Open("testdata/estaciones_php")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	stations, err := extractStations(f)
	require.Nil(t, err)

	assert.Len(t, stations, 165)
	// first station
	assert.Equal(t, "81-inia", stations[0].ID)
	assert.Equal(t, "Azapa Alto", stations[0].Name)
	assert.EqualValues(t, -18.577454, stations[0].Latitude)
	assert.EqualValues(t, -69.9472636, stations[0].Longitude)
	// last station
	assert.Equal(t, "13-ceaza", stations[164].ID)
	assert.Equal(t, "Vicuña [INIA]", stations[164].Name)
	assert.EqualValues(t, -30.03832, stations[164].Latitude)
	assert.EqualValues(t, -70.69655, stations[164].Longitude)
}
