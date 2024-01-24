package main

import (
	"fmt"
	"os"
	"serifhealth-takehome/config"
	"serifhealth-takehome/parser"
	"time"

	"gopkg.in/yaml.v2"
)

func main() {
	defer timer()()

	// Retrieving configuration
	config, err := loadConfig()
	if err != nil {
		fmt.Printf("Failed to retrieve config: %v\n", err)
		return
	}

	if len(config.FilePath) == 0 {
		fmt.Printf("Please set the file path to your json file in the config.yml file!\n")
		return
	}

	// Parsing the file
	fileParser := parser.NewParser(config)
	urls, err := fileParser.ParseFile()
	if err != nil {
		fmt.Printf("Failed to parse file: %v\n", err)
		return
	}

	fmt.Printf("Num urls: %v\n", len(urls))
	err = writeOutput(urls)
	if err != nil {
		fmt.Printf("Failed to write output: %v\n", err)
		return
	}
}

func timer() func() {
	start := time.Now()
	return  func() {
		fmt.Printf("Application took %v\n", time.Since(start))
	}
}

func loadConfig() (*config.Config, error) {
	configFile, err := os.Open("./config.yml")
	if err != nil {
		return nil, err
	}
	defer configFile.Close()
	
	decoder := yaml.NewDecoder(configFile)

	var config config.Config
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func writeOutput(urls []string) error {
	output, err := os.Create("output.txt")
	if err != nil {
		return err
	}
	defer output.Close()

	for _, url := range urls {
		_, err = output.WriteString(fmt.Sprintf("%v\n\n", url))
		if err != nil {
			return err
		}
	}

	return nil
}
