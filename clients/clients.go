package clients

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Currency struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

type Country struct {
	Name       string     `json:"name"`
	Capital    string     `json:"capital"`
	Region     string     `json:"region"`
	Population int64      `json:"population"`
	Currencies []Currency `json:"currencies"`
	FlagURL    string     `json:"flag" gorm:"size:512"`
}

type ExchangeRates struct {
	Rates map[string]float64 `json:"rates"`
}

// Fetches a list countries from the RestCountries REST API
func GetCountries() (*[]Country, error) {
	var countries []Country

	client := http.Client{}
	url := "https://restcountries.com/v2/all?fields=name,capital,region,population,flag,currencies"
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Failed to make GET request because ", err.Error())
		return nil, err
	}

	response, err := client.Do(request)
	if err != nil {
		log.Println("Failed to perform GET request because", err.Error())
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		err := fmt.Errorf("unexpected status code: %d", response.StatusCode)
		log.Println("Failed to make request:", err)
		return nil, err
	}

	// Decode the response body into Go struct
	if err := json.NewDecoder(response.Body).Decode(&countries); err != nil {
		log.Println("Failed to decode response body", err.Error())
		return nil, err
	}

	return &countries, nil
}

func GetExchangeRates() (*ExchangeRates, error) {
	var rates ExchangeRates

	// Make HTTP request
	client := http.Client{}

	url := "https://open.er-api.com/v6/latest/USD"

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Failed to make GET request because", err.Error())
		return nil, err
	}

	response, err := client.Do(request)
	if err != nil {
		log.Println("Failed to make HTTP GET request because", err.Error())
		return nil, err
	}

	if response.StatusCode != 200 {
		err := fmt.Errorf("unexpected status code: %d", response.StatusCode)
		log.Println("Failed to make request:", err)
		return nil, err
	}

	// Decode the response body into Go struct
	if err := json.NewDecoder(response.Body).Decode(&rates); err != nil {
		log.Println("Failed to decode response body", err.Error())
		return nil, err
	}
	return &rates, nil
}
