package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type State struct {
	Containers []string
	Freezers   []Freezer
}

type Freezer struct {
	Name     string
	Contents []Item
}

type Item struct {
	Name       string
	Date       string
	Containers []string
}

func readContents(f string) (State, error) {
	jsonBytes, err := os.ReadFile(f)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return State{}, err
	}

	var state State
	err = json.Unmarshal(jsonBytes, &state)
	if err != nil {
		fmt.Println("Error parsing file:", err)
		return State{}, err
	}

	return state, nil
}

func writeContents(f string, state State) error {
	fmt.Println("Writing to file:", f)

	jsonData, err := json.MarshalIndent(state, "", "  ")

	if err != nil {
		fmt.Println("Error writing file:", err)
		return err
	}

	file, err := os.Create(f)
	if err != nil {
		fmt.Println("Error creating file", err)
		return err
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			fmt.Println("Error closing file:", cerr)
		}
	}()

	_, err = file.Write(jsonData)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return err
	}

	return nil
}
