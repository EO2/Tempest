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
	alt          int8    = 0    // alt altitude / msl meters above sea level / negative number? / 5 default
	debug        bool    = true // false
)

var (
	userAgent string = fmt.Sprintf("Tempest/%v.%v github.com/EO2", versionUpper, versionLower)
	url       string = fmt.Sprintf("https://api.met.no/weatherapi/locationforecast/2.0/compact?lat=%.4f&lon=%.4f&altitude=%d", lat, lon, alt)
)

var (
	data         Sensors
	updatedAt    time.Time
	lastModified time.Time
	expires      time.Time
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
	updatedAt, err = time.Parse(time.RFC3339, data.Properties.Meta.UpdatedAt)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Get if over an hour since last modified
	// hourAgo := time.Now().Add(-time.Hour * 1)
	/* Todo:
	Compression (Accept-Encoding: gzip, deflate) as described in RFC 2616.
	Redirects yes
	cache headers, see RFC 2616. For example, use If-Modified-Since requests if the Last-Modified header exists.
	Note that the If-Modified-Since header should be identical to the previous Last-Modified,
	not any random timestamp (and definitely not in the future).
	If 429 - Throttling: Limit traffic / request frequency
	Randomize interval (30 - 60 min?)
	*/

	if debug {
		fmt.Println("updatedAt", updatedAt.Format(time.RFC1123))
		fmt.Println("Last-Modified:", lastModified.Format(time.RFC1123))
		fmt.Println("Expires:", expires.Format(time.RFC1123))
	}

	if expires.Before(time.Now().Add(-time.Minute * 1)) {
		fmt.Println("Expired - ReGet..")
		getSensors()
		time.Sleep(5 * time.Second)
		getSensors()
	}

	//fmt.Println(data.Properties.Timeseries[0].Data.Instant.Details.AirTemperature)

	/*for _, forecast := range data.Properties.Timeseries {
		ts, _ := time.Parse(time.RFC3339, forecast.Time)
		t := forecast.Data.Instant.Details.AirTemperature
		rh := forecast.Data.Instant.Details.RelativeHumidity
		fmt.Printf("%v\t\t%.1f C\t\t%.1f rH\n", ts.Format("01.02 15:04"), t, rh)
	}*/
}

func getSensors() {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("If-Modified-Since", lastModified.Format(time.RFC1123))

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
			return // Still need to save Expires and Last-Modified? Test..
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

	if debug {
		fmt.Println("updatedAt:", updatedAt.Format(time.RFC1123))
		fmt.Println("Last-Modified:", lastModified.Format(time.RFC1123))
		fmt.Println("Expires:", expires.Format(time.RFC1123))
		fmt.Println("Req Header:", req.Header)
		fmt.Println("Resp Header:", resp.Header)
		fmt.Println("Gzipped?:", resp.Uncompressed) // yes
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

func checkModified() {
	// Head request to check if file is changed, without downloading json body.
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Todo: any readon to use app/json here? should not.
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("If-Modified-Since", lastModified.Format(time.RFC1123))

	if debug {
		fmt.Println(req.Header)
	}

	resp, err := new(http.Client).Do(req) // Todo: reuse client, need to close?
	if err != nil {
		log.Fatal(err)
	}

	if debug {
		fmt.Println(resp.StatusCode)
	}

	// if 200, return ok. Then get request for data
	//fmt.Println(resp.StatusCode == 200)
	//fmt.Println(resp)
}
