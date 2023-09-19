package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/EAliakbar/traefik-file-provider-mapper/internal"
)

func main() {
	inputPath := "sample-inputs/"
	mapperConfig := internal.ReadMapperConfigFromENV()

	// Set up an HTTP endpoint to serve the merged YAML data.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		jsonData := ProcessConfigs(inputPath, mapperConfig)
		w.Write(jsonData)
	})

	// Start the HTTP server.
	fmt.Println("Listening on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func ProcessConfigs(inputPath string, mapperConfig internal.MapperConfig) []byte {
	inputConfigs, err := internal.ReadConfigs(inputPath)
	if err != nil {
		fmt.Println("Couldn't Read Configs")
	}
	mergedConfig := internal.MergeConfigurations(inputConfigs)
	mergedConfig.Mapper(mapperConfig)

	jsonData, err := json.Marshal(mergedConfig)
	if err != nil {
		fmt.Println("Error Convert configs to json", err)
		return nil
	}
	return jsonData
}
