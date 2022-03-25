// MET API: Temperature & Relative Humidity from GPS location
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	versionUpper uint8   = 0
	versionLower uint8   = 1
	lat          float32 = 60.3929
	lon          float32 = 5.3241
	alt          int8    = 0 // alt altitude / msl meters above sea level / negative number? / 5 default
)

var (
	userAgent string = fmt.Sprintf("Tempest/%v.%v github.com/EO2", versionUpper, versionLower)
	url       string = fmt.Sprintf("https://api.met.no/weatherapi/locationforecast/2.0/compact?lat=%.4f&lon=%.4f", lat, lon)
)

var (
	data         Sensors
	dataJSON     []byte
	expires      time.Time
	lastModified time.Time
)

func init() {
	content, errJSON := ioutil.ReadFile("Sensors.json")
	if errJSON != nil {
		fmt.Println(errJSON)
	}
	err := json.Unmarshal([]byte(content), &data)
	if err != nil {
		fmt.Println("Unmarshal Sensors: %w", err)
	}
	lastModified, err = time.Parse(time.RFC3339Nano, data.Properties.Meta.UpdatedAt)
	expires = lastModified // Todo: Check expires with Head
}

func main() {
	// Get again if over an hour since last modified
	if expires.Before(time.Now().Add(-time.Hour * 1)) {
		fmt.Println("Data expired - Get again..")
		// then look at lastModified, then download.
		getSensors()
	} else {
		fmt.Println("Sensors up to date")
	}

	for _, forecast := range data.Properties.Timeseries {
		ts, _ := time.Parse(time.RFC3339, forecast.Time)
		t := forecast.Data.Instant.Details.AirTemperature
		rh := forecast.Data.Instant.Details.RelativeHumidity
		fmt.Printf("%v\t\t%.1f C\t\t%.1f rH\n", ts.Format("01.02 15:04"), t, rh)
	}
}

func getSensors() {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	// If-Modified-Since (lastModified) 304 not modified or 200
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalf("MET Response StatusCode: %v", resp.StatusCode)
	} // Todo: 203 (deprecated product) and 429 (throttling)

	expires, err = time.Parse(time.RFC1123, resp.Header.Get("Expires"))
	lastModified, err = time.Parse(time.RFC1123, resp.Header.Get("Last-Modified"))

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	dataJSON = body
	data = Sensors{}
	err = json.Unmarshal(body, &data)

	if err := os.WriteFile("Sensors.json", body, 0666); err != nil {
		log.Fatal(err)
	}
}
