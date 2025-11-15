package clients

// McpGeocodeResponse is the response from /mcp/geo/geocode
type McpGeocodeResponse struct {
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	DisplayName string  `json:"display_name"`
}

type McpNearbyPlace struct {
	Name string                 `json:"name"`
	Lat  float64                `json:"lat"`
	Lon  float64                `json:"lon"`
	Tags map[string]interface{} `json:"tags"`
}

// McpNearbyResponse is the response from /mcp/geo/nearby
type McpNearbyResponse struct {
	Places []McpNearbyPlace `json:"places"`
}

// McpForecastResponse is the response from /mcp/weather/forecast
type McpForecastResponse struct {
	TempC      float64 `json:"temp_c"`
	PrecipProb float64 `json:"precip_prob"`
	WindKph    float64 `json:"wind_kph"`
	Summary    string  `json:"summary"`
}

// McpAirQualityResponse is the response from /mcp/air/aqi
type McpAirQualityResponse struct {
	PM25     float64 `json:"pm25"`
	PM10     float64 `json:"pm10"`
	NO2      float64 `json:"no2"`
	O3       float64 `json:"o3"`
	Category string  `json:"category"`
}

// McpEtaResponse is the response from /mcp/route/eta
type McpEtaResponse struct {
	DistanceKm  float64 `json:"distance_km"`
	DurationMin float64 `json:"duration_min"`
	Polyline    string  `json:"polyline"`
}
