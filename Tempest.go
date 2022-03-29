// Temperature & Relative Humidity from GPS location (https://github.com/EO2/Tempest)
// Data from The Norwegian Meteorological Institute (https://api.met.no/doc/)
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
	lat          float32 = 60.3929 // max 4 decimal coords
	lon          float32 = 5.3241
	alt          int8    = 0 // meters above sea level. Negative number possible.. 5 default
)

var (
	userAgent string = fmt.Sprintf("Tempest/%d.%d github.com/EO2/Tempest", versionUpper, versionLower)
	url       string = fmt.Sprintf("https://api.met.no/weatherapi/locationforecast/2.0/compact?lat=%.4f&lon=%.4f&altitude=%d", lat, lon, alt)
)

var (
	data         Sensors
	lastModified time.Time
	expires      time.Time
)

func init() {
	fmt.Printf("Tempest v.%d.%d\n", versionUpper, versionLower)
	fmt.Println("Get Weather Forecasts from MET.no API")
	content, errJSON := ioutil.ReadFile("Sensors.json")
	if errJSON != nil {
		fmt.Println(errJSON)
	}
	err := json.Unmarshal([]byte(content), &data)
	if err != nil {
		fmt.Println("Unmarshal Sensors: %w", err)
	}
}

func main() {
	// Get Weather from MET.no API
	for {
		if expires.IsZero() || expires.Before(time.Now()) {
			getSensors()
			fmt.Println("MET Forecast:", time.Now().Format(time.RFC1123), "> ", data.Properties.Timeseries[0].Data.Instant.Details.AirTemperature, "℃ ", data.Properties.Timeseries[0].Data.Instant.Details.RelativeHumidity, "% RH")

			//fmt.Println(data.Properties.Timeseries[0].Data.Instant.Details.AirTemperature)
			//fmt.Println(data.Properties.Timeseries[0].Data.Instant.Details.RelativeHumidity)

			/*for _, forecast := range data.Properties.Timeseries {
				ts, _ := time.Parse(time.RFC3339, forecast.Time)
				t := forecast.Data.Instant.Details.AirTemperature
				rh := forecast.Data.Instant.Details.RelativeHumidity
				fmt.Printf("%v\t\t%.1f C\t\t%.1f rH\n", ts.Format("01.02 15:04"), t, rh)
			}*/
		}
		time.Sleep(5 * time.Minute)
	}
}

func getSensors() {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)
	if !lastModified.IsZero() {
		req.Header.Set("If-Modified-Since", lastModified.Format(time.RFC1123))
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		if resp.StatusCode == 203 {
			log.Fatal("MET: 203 - Service Deprecated!")
		}
		if resp.StatusCode == 429 {
			log.Fatal("MET: 429 - Service Throttling!")
		}
		if resp.StatusCode == 304 {
			fmt.Println("MET: 304 - Resource Not Modified")
			return // No Body
		}
		log.Fatalf("MET Response StatusCode: %v", resp.StatusCode)
	}

	expires, err = time.Parse(time.RFC1123, resp.Header.Get("Expires"))
	if err != nil {
		log.Fatal(err)
	}

	lastModified, err = time.Parse(time.RFC1123, resp.Header.Get("Last-Modified"))
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile("Sensors.json", body, 0666); err != nil {
		log.Fatal(err)
	}
}
