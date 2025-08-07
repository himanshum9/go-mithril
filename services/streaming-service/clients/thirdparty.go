package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type LocationData struct {
	TenantID  string  `json:"tenant_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp string  `json:"timestamp"`
}

type ThirdPartyClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewThirdPartyClient(baseURL string) *ThirdPartyClient {
	return &ThirdPartyClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *ThirdPartyClient) SendLocationData(data LocationData) error {
	url := c.BaseURL + "/api/location"
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send location data, status code: %d", resp.StatusCode)
	}

	return nil
}
