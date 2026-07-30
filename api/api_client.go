package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/akash329d/storj_exporter/models"
)

type ApiClient struct {
	BaseURL    string
	NodeID     string
	NodeName   string
	Satellites []models.Satellite
	httpClient *http.Client
}

func NewApiClient(baseURL, nodeName string) *ApiClient {
	client := &ApiClient{
		BaseURL:  baseURL,
		NodeName: nodeName,
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}

	nodeData, err := client.Node()
	if err != nil {
		panic(fmt.Sprintf("failed to make initial connection to node %s: %v", baseURL, err))
	}

	client.NodeID = nodeData.NodeID
	client.Satellites = nodeData.Satellites

	return client
}

func (c *ApiClient) get(endpoint string, target interface{}) error {
	resp, err := c.httpClient.Get(c.BaseURL + endpoint)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(target)
}

func (c *ApiClient) Node() (models.NodeData, error) {
	var data models.NodeData
	err := c.get("/api/sno/", &data)
	if err != nil {
		return data, fmt.Errorf("API Request for node data failed: %w", err)
	}
	return data, nil
}

func (c *ApiClient) Payout() (models.PayoutResponse, error) {
	var data models.PayoutResponse
	err := c.get("/api/sno/estimated-payout", &data)
	if err != nil {
		return data, fmt.Errorf("API Request for payout data failed: %w", err)
	}
	return data, nil
}

func (c *ApiClient) Satellite(satelliteId string) (models.SatelliteResponse, error) {
	var data models.SatelliteResponse
	satelliteApiUrl := fmt.Sprintf("/api/sno/satellite/%s", satelliteId)
	err := c.get(satelliteApiUrl, &data)
	if err != nil {
		return data, fmt.Errorf("API Request for sattelite data failed with API URL %s, %w",satelliteApiUrl, err)
	}
	return data, nil
}
