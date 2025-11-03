package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type StatsByIp struct {
	IpAddress string	`json:"ipAddress"`
	TotalLogs int `json:"totalLogs"`
	TotalFound int `json:"totalFound"`
	TotalBanned int `json:"totalBanned"`
	TotalUnbanned int `json:"totalUnbanned"`
	Country string `json:"county"`
}

func analyseLogs() {

	data, err := os.ReadFile(ParsedJson)
	if err != nil {
		fmt.Printf("Error opening file: %v", err)
		return
	}

	statsByIp := make(map[string]*StatsByIp)
	var fileData []Logs
	err = json.Unmarshal(data, &fileData)
	if err != nil {
		fmt.Printf("Error unmarshaling data: %v", err)
		return
	}

	for _, entry := range fileData {

		ip := entry.IpAddress
		if ip == "" {
			continue
		}

		action := entry.Message


		if _, exists := statsByIp[ip]; !exists {
			country, err := getIpAddressCountry(ip)
			if err != nil {
				fmt.Printf("Failed getting ip-address with error: %v", err)
			}
			statsByIp[ip] = &StatsByIp{
				IpAddress: ip,
				Country: country,
			}
		}

		statsByIp[ip].TotalLogs += 1

		switch action {
			case "Found":
				statsByIp[ip].TotalFound += 1
			case "Ban", "already banned":
				statsByIp[ip].TotalBanned += 1
			case "Unban":
				statsByIp[ip].TotalUnbanned += 1
		}
	}

	jsonData, err := json.MarshalIndent(statsByIp, "", "   ")
	if err != nil {
		fmt.Println("Error marshalling the stats file.")
		return
	}
	err = os.WriteFile(StatsByIPFile, jsonData, 0644)

}


func analyseExtractedData() {

	totalCounter := make(map[string]float64)
	totalBans := make(map[string]float64)
	totalConnections := make(map[string]float64)

	data, err := os.ReadFile(StatsByIPFile)
	if err != nil {
		fmt.Printf("Error opening file: %v", err)
		return
	}

	var statsByIp map[string]StatsByIp
	err = json.Unmarshal(data, &statsByIp)
	if err != nil {
		fmt.Printf("Error unmarshaling data: %v", err)
		return
	}

	for _, entry := range statsByIp {
		country := entry.Country

		if _, exists := totalCounter[country]; !exists {
			totalCounter[country] = 0
		} else {
			totalCounter[country] += 1
		}
		if entry.TotalBanned > 0 {
			if _, exists := totalBans[country]; !exists {
				totalBans[country] = float64(entry.TotalBanned)
			} else {
				totalBans[country] += float64(entry.TotalBanned)
			}
		}
		if entry.TotalFound > 0 {
			if _, exists := totalConnections[country]; !exists {
				totalConnections[country] = float64(entry.TotalFound)
			} else {
				totalConnections[country] += float64(entry.TotalFound)
			}
		}
	}

	barChart("Individual IPs", "Country", totalCounter)
	barChart("Total Banned Ips", "Country", totalBans)
	barChart("Total found connections", "Country", totalConnections)

}
