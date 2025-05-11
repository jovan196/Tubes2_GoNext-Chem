package main

import (
	"encoding/json"
	"fmt"
	"os"
)

var Graph map[string][][2]string

func LoadGraph(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&Graph)
	if err != nil {
		return err
	}
	fmt.Printf("Loaded %d elements from %s\n", len(Graph), filename)
	return nil
}

func LoadTier(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(&Tier)
}
