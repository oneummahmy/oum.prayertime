package oum_prayertime

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gocolly/colly"
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

type Reversegeocode struct {
	XMLName     xml.Name `xml:"reversegeocode"`
	Text        string   `xml:",chardata"`
	Timestamp   string   `xml:"timestamp,attr"`
	Attribution string   `xml:"attribution,attr"`
	Querystring string   `xml:"querystring,attr"`
	Result      struct {
		Text        string `xml:",chardata"`
		PlaceID     string `xml:"place_id,attr"`
		OsmType     string `xml:"osm_type,attr"`
		OsmID       string `xml:"osm_id,attr"`
		Lat         string `xml:"lat,attr"`
		Lon         string `xml:"lon,attr"`
		Boundingbox string `xml:"boundingbox,attr"`
		PlaceRank   string `xml:"place_rank,attr"`
		AddressRank string `xml:"address_rank,attr"`
	} `xml:"result"`
	Addressparts struct {
		Text         string `xml:",chardata"`
		Village      string `xml:"village"`
		District     string `xml:"district"`
		State        string `xml:"state"`
		ISO31662Lvl4 string `xml:"ISO3166-2-lvl4"`
		Postcode     string `xml:"postcode"`
		Country      string `xml:"country"`
		CountryCode  string `xml:"country_code"`
	} `xml:"addressparts"`
}

type JakimLocation struct {
	District string `json:"district"`
	Zone     string `json:"zone"`
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

func GetPrayerTimeByLatLng(lat float64, lng float64) {
	res, err := ioutil.ReadFile("jakim_location.json")
	if err != nil {
		panic(err)
	}

	var jakim []JakimLocation

	err = json.Unmarshal(res, &jakim)
	if err != nil {
		panic(err)
	}

	rev, err := GetInfoByCoordinate(lat, lng)
	if err != nil {
		panic(err)
	}

	var location Location
	for _, v := range jakim {
		if v.District == rev.Addressparts.District {
			location = Location{
				Code: v.Zone,
			}
		}
	}

	prayer, err := GetPrayerTime(location)
	if err != nil {
		panic(err)
	}

	fmt.Println(prayer.PrayerTime)
}

func GetInfoByCoordinate(lat float64, lng float64) (*Reversegeocode, error) {
	requestURL := fmt.Sprintf("https://nominatim.openstreetmap.org/reverse?lat=%f&lon=%f", lat, lng)
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	byteValue, _ := ioutil.ReadAll(res.Body)

	var reverseGeo Reversegeocode

	xml.Unmarshal(byteValue, &reverseGeo)

	return &reverseGeo, nil
}
