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

func (mc *McpClient) newRequest(method, path string, body interface{}) (*http.Request, error) {
	fullURL, err := url.JoinPath(mc.BaseURL, path)
	if err != nil {
		return nil, err
	}

	var bodyReader io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, fullURL, bodyReader)
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
	path := fmt.Sprintf("/mcp/geo/geocode?city=%s", url.QueryEscape(city))

	req, err := mc.newRequest("GET", path, nil)
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
	path := fmt.Sprintf("/mcp/weather/forecast?lat=%f&lon=%f&date=%s", lat, lon, date)

	req, err := mc.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var response McpForecastResponse
	if err := mc.do(req, &response); err != nil {
		return nil, err
	}
	return &response, nil
}
