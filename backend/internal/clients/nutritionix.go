package clients

import (
	"backend/internal/config"
	"backend/internal/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type NutritionixClient struct {
	AppID      string
	AppKey     string
	BaseURL    string
	HTTPClient *http.Client
}

func NewNutritionixClient(envs *config.Envs) *NutritionixClient {
	return &NutritionixClient{
		AppID:      envs.NutritionixAppID,
		AppKey:     envs.NutritionixAppKey,
		BaseURL:    "https://trackapi.nutritionix.com/v2",
		HTTPClient: &http.Client{},
	}
}

func (c *NutritionixClient) GetNutritionData(query string) (*models.NutritionixResponse, error) {
	endpoint := c.BaseURL + "/natural/nutrients"

	requestBody := struct {
		Query string `json:"query"`
	}{
		Query: query,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("x-app-id", c.AppID)
	req.Header.Set("x-app-key", c.AppKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s, response %s", resp.Status, string(body))
	}

	var response models.NutritionixResponse

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(response.Foods) == 0 {
		return nil, fmt.Errorf("no food data found for query: %s", query)
	}

	return &response, nil

}
