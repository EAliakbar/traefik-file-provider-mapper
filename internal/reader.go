package internal

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

func ReadConfigs(dirPath string) (Configurations, error) {
	configs := make(Configurations)

	// Read YAML files from the specified directory.
	err := filepath.Walk(dirPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(filePath) == ".yaml" {
			fmt.Printf("Found Config: %s\n", filePath)
			data, err := os.ReadFile(filePath)
			if err != nil {
				return err
			}

			var yamlData Configuration
			if err := yaml.Unmarshal(data, &yamlData); err != nil {
				return err
			}

			configs[filePath] = &yamlData
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return configs, nil
}
