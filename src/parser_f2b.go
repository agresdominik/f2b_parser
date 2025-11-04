package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Logs struct {
	Timestamp 	string	`json:"timestamp"`
	Handler		string  `json:"handler"`
	Level 		string  `json:"level"`
	Source 		string  `json:"source"`
	IpAddress	string  `json:"ipAddress"`
	Message 	string  `json:"message"`
}

var LogFile string = "./data/fail2ban.log"
var ParsedJson string = "./data/json_output.json"
var StatsByIPFile string = "./data/stats_by_ip.json"


func parseFile(stateFilePath string, logFilePath string, destinationDirectory string) {

	// Init Parsed Log File Name
	destinationFilePath := filepath.Join(destinationDirectory, "parsed.json")

	// Load metadata
	metadata := checkState(stateFilePath)
	offset := metadata.Offset

	// Check if the log file has rolled over
	file, err := os.Open(logFilePath)
	if err != nil {
		fmt.Printf("Error opening file: $%v", err)
		return
	}
	defer file.Close()

	stat, _ := file.Stat()
	if stat.Size() < offset {
		offset = 0
	}

	// Read out existing parsed log file
	var logs []Logs
    if data, err := os.ReadFile(destinationFilePath); err == nil && len(data) > 0 {
        _ = json.Unmarshal(data, &logs)
    }

    // Define regex logic
	const lenTimestamp = 23
	dateRegex := regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)
	handlerRegex := regexp.MustCompile(`fail2ban\.\w+`)
	ipRegex := regexp.MustCompile(`(\b25[0-5]|\b2[0-4][0-9]|\b[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)
	levelRegex := regexp.MustCompile(`\s*(?:[A-Z]+)\s+`)
	serviceRegex := regexp.MustCompile(`\s*(?:\[[a-z]+\])\s+`)
	actionRegex := regexp.MustCompile(`(Found|already banned|Ban|Unban)`)

	logEntry := Logs{}


	_, err = file.Seek(offset, io.SeekStart)
	if err != nil {
		log.Fatalf("Error going to offset: %s\n", err)
	}

	// Parse the file and append to existing log files
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		line := scanner.Text()

		if len(line) < lenTimestamp {
			continue
		} else if !dateRegex.MatchString(line[:lenTimestamp]) {
			continue
		}

		timestamp := line[:lenTimestamp]; timestamp = strings.TrimSpace(timestamp)
		logString := line[lenTimestamp:]

		ipAddress := strings.TrimSpace(ipRegex.FindString(logString))
		handler := strings.TrimSpace(handlerRegex.FindString(logString))
		level := strings.TrimSpace(levelRegex.FindString(logString))
		service := strings.TrimSpace(serviceRegex.FindString(logString))
		action := strings.TrimSpace(actionRegex.FindString(logString))

		logEntry.IpAddress = ipAddress
		logEntry.Timestamp = timestamp
		logEntry.Handler = handler
		logEntry.Level = level
		logEntry.Source = service
		logEntry.Message = action

		logs = append(logs, logEntry)
	}

	// Write parsed content and update metadata
	jsonData, err := json.MarshalIndent(logs, "", "   ")
	_ = os.WriteFile(destinationFilePath, jsonData, 0644)

	// Get parsed log file size
	logFile, err := os.Open(destinationFilePath)
	if err != nil {
		fmt.Printf("Error opening file: $%v", err)
		return
	}
	defer logFile.Close()
	stat, _ = logFile.Stat()
	parsedLogfileSize := stat.Size()

	newOffset, _ := file.Seek(0, io.SeekCurrent)
	newState := State{
		Offset: newOffset,
		ParsedFileSize: parsedLogfileSize,
	}
	updateState(stateFilePath, newState)
}
