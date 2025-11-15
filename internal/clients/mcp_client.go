package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

type McpClient struct {
	Client  *http.Client
	BaseURL string
	APIKey  string
}

func NewMcpClient(baseURL, apiKey string) *McpClient {
	return &McpClient{
		Client: &http.Client{
			Timeout: 15 * time.Second,
		},
		BaseURL: baseURL,
		APIKey:  apiKey,
	}
}

func (mc *McpClient) newRequest(method, path string, query url.Values, body interface{}) (*http.Request, error) {

	fullURL, err := url.Parse(mc.BaseURL + path)
	if err != nil {
		return nil, err
	}

	fullURL.RawQuery = query.Encode()

	log.Printf("MCP REQUEST URL: %s", fullURL.String())

	var bodyReader io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, fullURL.String(), bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-API-KEY", mc.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	return req, nil
}

func (mc *McpClient) do(req *http.Request, v interface{}) error {
	resp, err := mc.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("MCP server returned error: %s", resp.Status)
	}

	if v != nil {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return fmt.Errorf("failed to decode MCP response: %v", err)
		}
	}
	return nil
}

func (mc *McpClient) GetGeocode(city string) (*McpGeocodeResponse, error) {
	log.Printf("MCP Client: Geocoding city: %s", city)

	query := url.Values{}
	query.Set("city", city)

	req, err := mc.newRequest("GET", "/mcp/geo/geocode", query, nil)
	if err != nil {
		return nil, err
	}

	var response McpGeocodeResponse
	if err := mc.do(req, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (mc *McpClient) GetForecast(lat, lon float64, date string) (*McpForecastResponse, error) {
	log.Printf("MCP Client: Getting forecast for %s", date)

	query := url.Values{}
	query.Set("lat", fmt.Sprintf("%f", lat))
	query.Set("lon", fmt.Sprintf("%f", lon))
	query.Set("date", date)

	req, err := mc.newRequest("GET", "/mcp/weather/forecast", query, nil)
	if err != nil {
		return nil, err
	}

	var response McpForecastResponse
	if err := mc.do(req, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (mc *McpClient) GetNearby(lat, lon float64, query string) (*McpNearbyResponse, error) {
	log.Printf("MCP Client: Getting nearby venues for query: %s", query)

	queryVals := url.Values{}
	queryVals.Set("lat", fmt.Sprintf("%f", lat))
	queryVals.Set("lon", fmt.Sprintf("%f", lon))
	queryVals.Set("query", query)
	queryVals.Set("radius_m", "5000")
	queryVals.Set("limit", "15")

	req, err := mc.newRequest("GET", "/mcp/geo/nearby", queryVals, nil)
	if err != nil {
		return nil, err
	}

	var response McpNearbyResponse
	if err := mc.do(req, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (mc *McpClient) GetAQI(lat, lon float64, date string) (*McpAirQualityResponse, error) {
	log.Printf("MCP Client: Getting AQI for %f, %f", lat, lon)

	queryVals := url.Values{}
	queryVals.Set("lat", fmt.Sprintf("%f", lat))
	queryVals.Set("lon", fmt.Sprintf("%f", lon))
	queryVals.Set("date", date)

	req, err := mc.newRequest("GET", "/mcp/air/aqi", queryVals, nil)
	if err != nil {
		return nil, err
	}

	var response McpAirQualityResponse
	if err := mc.do(req, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (mc *McpClient) GetHolidays(countryCode string, year int) (McpHolidayResponse, error) {
	log.Printf("MCP Client: Getting holidays for %s, %d", countryCode, year)

	queryVals := url.Values{}
	queryVals.Set("country_code", countryCode)
	queryVals.Set("year", fmt.Sprintf("%d", year))

	req, err := mc.newRequest("GET", "/mcp/calendar/holidays", queryVals, nil)
	if err != nil {
		return nil, err
	}

	var response McpHolidayResponse
	if err := mc.do(req, &response); err != nil {
		return nil, err
	}
	return response, nil
}

func (mc *McpClient) GetEta(reqBody EtaRequest) (*McpEtaResponse, error) {
	log.Printf("MCP Client: Getting ETA for profile: %s", reqBody.Profile)

	req, err := mc.newRequest("POST", "/mcp/route/eta", nil, reqBody)
	if err != nil {
		return nil, err
	}

	var response McpEtaResponse
	if err := mc.do(req, &response); err != nil {
		return nil, err
	}
	return &response, nil
}
