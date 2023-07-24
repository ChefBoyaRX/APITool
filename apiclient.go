package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type APIClient struct {
	Status       string   `json:"status"`
	ResponseTime int      `json:"responseTime"`
	Message      []string `json:"message"`
	Results      Results  `json:"Results"`
}

type Results struct {
	Series []Series `json:"series"`
}

type Series struct {
	SeriesID string `json:"seriesID"`
	Data     []Data `json:"data"`
}

type Data struct {
	Year       string      `json:"year"`
	Period     string      `json:"period"`
	PeriodName string      `json:"periodName"`
	Value      string      `json:"value"`
	Footnotes  []Footnotes `json:"footnotes"`
}

type Footnotes struct {
	Code string `json:"code"`
	Text string `json:"text"`
}


type PayloadData struct {
	SeriesID  []string `json:"seriesid"`
	StartYear string   `json:"startyear"`
	EndYear   string   `json:"endyear"`
	APIKey    string   `json:"registrationkey"`
}

func (c *APIClient) FetchData() ([]Series, []byte, error) {
	apiURL := "https://api.bls.gov/publicAPI/v2/timeseries/data/"
	
	apiKey, err := getAPIKeyFromFile()
    if err != nil {
        return nil, nil, err
    }

	// Modify the series and year data here if needed
	series1 := "CUUR0000SA0"
	series2 := "CES0500000003"
	startYear := "2012"
	endYear := "2022"

	// Create the payload data
	payloadData := PayloadData{
		SeriesID:  []string{series1, series2},
		StartYear: startYear,
		EndYear:   endYear,
		APIKey:    apiKey,
	}

	// Convert payload data to JSON
	payload, err := json.Marshal(payloadData)
	if err != nil {
		return nil, nil, err
	}

	// Create a new HTTP POST request
	req, err := http.NewRequest("POST", apiURL, bytes.NewReader(payload))
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	
	fmt.Println("Raw JSON Response:")
	fmt.Println(string(body))

	// Parse the JSON response into the APIClient struct
	var responseData APIClient
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return nil, nil, err
	}

	
	return responseData.Results.Series, body, nil
}

func getAPIKeyFromFile() (string, error) {
	apiKey, err := os.ReadFile("api_key.txt")
    if err != nil {
        return "", err
    }
    return strings.TrimSpace(string(apiKey)), nil
}