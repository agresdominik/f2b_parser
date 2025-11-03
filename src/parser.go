package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type State struct {
	Offset		int64 	`json:"offset"`
}

func starter() {

	destinationDirectory := flag.String("destDir", "", "Destination Directory")
	flag.StringVar(destinationDirectory, "d", "", "Destination Directory (shorthand)")
	source := flag.String("source", "", "Source Log File")
	flag.StringVar(source, "s", "", "Source Log File (shorthand)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s --destDir <dir> --source <file>\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if *destinationDirectory == "" || *source == "" {
		flag.Usage()
		os.Exit(1)
	}

	checkParameters(*destinationDirectory, *source)

}

/*
 * TODO: This function does not check read/write permissions yet
 */
func checkParameters(destinationDirectory string, source string) {


	if _, err := os.Stat(source); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: source file does not exist: %s\n", source)
		os.Exit(1)
	} else if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed reading source file: %s with error: %s\n", source, err)
		os.Exit(1)
	}

	if _, err := os.Stat(destinationDirectory); os.IsNotExist(err){
		err = os.MkdirAll(destinationDirectory, 0755)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed creating directory file: %s with error: %s\n", source, err)
			os.Exit(1)
		}
	} else if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed reading directory file: %s with error: %s\nDid not try to create.", source, err)
		os.Exit(1)
	}

	stateFile, err := initState(destinationDirectory)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed initialising state file: %s\n", err)
		os.Exit(1)
	}

	parseFile(stateFile, source, destinationDirectory)
}

func initState(destinationDirectory string) (stateFilePath string, err error) {

	stateFile := filepath.Join(destinationDirectory, "state.json")

	if _, err := os.Stat(stateFile); err == nil {
		return stateFile, nil
	} else if os.IsNotExist(err) {
		fmt.Println("No state file found, creating...")
	}

	state := State{
		Offset: 0,
	}
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return "", err
	}

	err = os.WriteFile(stateFile, data, 0644)
	if err != nil {
		return "", err
	}

	fmt.Printf("State file initialised: %s\n", stateFile)
	return stateFile, nil
}

func checkState(stateFile string)  State {

	data, err := os.ReadFile(stateFile)
	if err != nil {
		fmt.Printf("Error opening file: %v", err)
		return State{}
	}

	var state State

	err = json.Unmarshal(data, &state)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading state file: %s", err)
		return State{}
	}
	return state
}

func updateState(stateFile string, newState State) error {

	data, err := json.MarshalIndent(newState, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(stateFile, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
