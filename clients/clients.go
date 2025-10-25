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
	Population float64    `json:"population"`
	Currencies []Currency `json:"currencies"`
	FlagURL    string     `json:"flag" gorm:"size:512"`
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