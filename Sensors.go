// Sensor weather data - Auto JSON => Struct
package main

// AutoStruct: http://json2struct.mervine.net/
// Fixed some ints that where floats + using float32 intead of 64.

type Sensors struct {
	Geometry struct {
		Coordinates []float32 `json:"coordinates"`
		Type        string    `json:"type"`
	} `json:"geometry"`
	Properties struct {
		Meta struct {
			Units struct {
				AirPressureAtSeaLevel string `json:"air_pressure_at_sea_level"`
				AirTemperature        string `json:"air_temperature"`
				CloudAreaFraction     string `json:"cloud_area_fraction"`
				PrecipitationAmount   string `json:"precipitation_amount"`
				RelativeHumidity      string `json:"relative_humidity"`
				WindFromDirection     string `json:"wind_from_direction"`
				WindSpeed             string `json:"wind_speed"`
			} `json:"units"`
			UpdatedAt string `json:"updated_at"`
		} `json:"meta"`
		Timeseries []struct {
			Data struct {
				Instant struct {
					Details struct {
						AirPressureAtSeaLevel float32 `json:"air_pressure_at_sea_level"`
						AirTemperature        float32 `json:"air_temperature"`
						CloudAreaFraction     float32 `json:"cloud_area_fraction"`
						RelativeHumidity      float32 `json:"relative_humidity"`
						WindFromDirection     float32 `json:"wind_from_direction"`
						WindSpeed             float32 `json:"wind_speed"`
					} `json:"details"`
				} `json:"instant"`
				Next12Hours struct {
					Summary struct {
						SymbolCode string `json:"symbol_code"`
					} `json:"summary"`
				} `json:"next_12_hours"`
				Next1Hours struct {
					Details struct {
						PrecipitationAmount float32 `json:"precipitation_amount"`
					} `json:"details"`
					Summary struct {
						SymbolCode string `json:"symbol_code"`
					} `json:"summary"`
				} `json:"next_1_hours"`
				Next6Hours struct {
					Details struct {
						PrecipitationAmount float32 `json:"precipitation_amount"`
					} `json:"details"`
					Summary struct {
						SymbolCode string `json:"symbol_code"`
					} `json:"summary"`
				} `json:"next_6_hours"`
			} `json:"data"`
			Time string `json:"time"`
		} `json:"timeseries"`
	} `json:"properties"`
	Type string `json:"type"`
}
