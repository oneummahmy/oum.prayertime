package oum_prayertime

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"net/http"
	"strings"
)

type Location struct {
	Code        string
	Description string
	State       string
}

type Prayer struct {
	PrayerTime []struct {
		Hijri   string `json:"hijri"`
		Date    string `json:"date"`
		Day     string `json:"day"`
		Imsak   string `json:"imsak"`
		Fajr    string `json:"fajr"`
		Syuruk  string `json:"syuruk"`
		Dhuhr   string `json:"dhuhr"`
		Asr     string `json:"asr"`
		Maghrib string `json:"maghrib"`
		Isha    string `json:"isha"`
	} `json:"prayerTime"`
	Status     string `json:"status"`
	ServerTime string `json:"serverTime"`
	PeriodType string `json:"periodType"`
	Lang       string `json:"lang"`
	Zone       string `json:"zone"`
	Bearing    string `json:"bearing"`
}

func GetLocations() ([]Location, error) {
	c := colly.NewCollector()

	var locations []Location

	c.OnHTML("#inputzone", func(e *colly.HTMLElement) {
		e.ForEach("optgroup", func(i int, element *colly.HTMLElement) {
			state := element.Attr("label")

			element.ForEach("option", func(i int, element2 *colly.HTMLElement) {
				descriptionSplit := strings.Split(element2.Text, "-")

				description := strings.TrimSpace(descriptionSplit[1])
				location := Location{
					Code:        element2.Attr("value"),
					Description: description,
					State:       state,
				}

				locations = append(locations, location)
			})
		})
	})

	err := c.Visit("https://www.e-solat.gov.my/index.php")
	if err != nil {
		return nil, err
	}

	return locations, nil
}

func GetPrayerTime(location Location) (*Prayer, error) {
	requestURL := fmt.Sprintf("https://www.e-solat.gov.my/index.php?r=esolatApi/TakwimSolat&period=today&zone=%s", location.Code)
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(res.Body)

	prayer := new(Prayer)
	err = decoder.Decode(prayer)
	if err != nil {
		return nil, err
	}

	return prayer, nil
}
