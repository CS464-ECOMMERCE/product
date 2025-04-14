package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"product/configs"
)

type GoogleMapsService struct {
	apiKey string
	client *http.Client
}

func NewGoogleMapsService() *GoogleMapsService {
	fmt.Println("Google Maps API Key:", configs.GOOGLE_MAPS_API_KEY)
	return &GoogleMapsService{
		apiKey: configs.GOOGLE_MAPS_API_KEY,
		client: &http.Client{},
	}
}

func (g *GoogleMapsService) ValidateAddress(address, country string) error {
	requestBody := map[string]interface{}{
		"address": map[string]interface{}{
			"addressLines": []string{address},
			"regionCode":   country,
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}
	// Create HTTP request
	url := fmt.Sprintf("https://addressvalidation.googleapis.com/v1:validateAddress?key=%s", g.apiKey)
	resp, err := g.client.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("API request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status: %s", resp.Status)
	}

	// Parse response
	var result struct {
		Result struct {
			Verdict struct {
				ValidationGranularity string `json:"validationGranularity"`
				AddressComplete       bool   `json:"addressComplete"`
			} `json:"verdict"`
		} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to parse response: %v", err)
	}
	fmt.Println(result)

	// Check if the address is INVALID
	if !result.Result.Verdict.AddressComplete ||
		result.Result.Verdict.ValidationGranularity == "OTHER" ||
		result.Result.Verdict.ValidationGranularity == "GRANULARITY_UNSPECIFIED" {
		return fmt.Errorf("address is not valid (granularity: %s), completeness: %v",
			result.Result.Verdict.ValidationGranularity,
			result.Result.Verdict.AddressComplete)
	}

	return nil
}
