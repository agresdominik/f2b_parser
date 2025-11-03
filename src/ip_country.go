package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type OkResponseJson struct {
	Ipaddress		string	`json:"ipaddress"`
	Continent_code 	string  `json:"continent_code"`
	Continent_name 	string	`json:"continent_name"`
	Country_code 	string	`json:"country_code"`
	Country_name 	string	`json:"country_name"`
}

func getIpAddressCountry(ipAddress string) (country string, err error) {

	url := fmt.Sprintf("https://api.ipaddress.com/iptocountry?format=json&ip=%s", ipAddress)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("Failed creating request: %v", err)
	}

	client := &http.Client{}
	response , err := client.Do(request)
	if err != nil {
		return "", fmt.Errorf("Failed making request: %v", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("Failed reading body: %s", err)
	}

	var responseFormat OkResponseJson
	err = json.Unmarshal(body, &responseFormat)
	if err != nil {
		return "", fmt.Errorf("Failed parsing response body with error: %v", err)
	}

	return responseFormat.Country_name, nil

}
