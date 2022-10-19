package oum_prayertime_test

import (
	oumprayertime "github.com/oneummahmy/oum.prayertime"
	"testing"
)

func TestGetLocations(t *testing.T) {
	locations, err := oumprayertime.GetLocations()
	if err != nil {
		t.Error(err)
	}

	if len(locations) == 0 {
		t.Error("it is supposed to be more than 0")
	}
}

func TestGetPrayerTime(t *testing.T) {
	location := oumprayertime.Location{
		Code:        "MLK01",
		Description: "Melaka Test",
		State:       "Melaka",
	}

	prayer, err := oumprayertime.GetPrayerTime(location)
	if err != nil {
		t.Error(err)
	}

	if prayer.Zone != "MLK01" {
		t.Error("it is supposed to be MLK01")
	}
}

func TestGetCoordinate(t *testing.T) {
	oumprayertime.GetPrayerTimeByLatLng(2.453873065668171, 102.18036740789528)
}
