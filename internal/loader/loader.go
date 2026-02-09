package loader

import (
	"os"

	"github.com/dawgdevv/probe/pkg/models"
	"gopkg.in/yaml.v3"
)

func LoadSuite(path string) (*models.TestSuite, error) {
	data, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	var suite models.TestSuite
	if err := yaml.Unmarshal(data, &suite); err != nil {
		return nil, err
	}

	return &suite, nil
}

// LoadSuiteFromString parses a YAML test suite from a string
func LoadSuiteFromString(yamlContent string) (*models.TestSuite, error) {
	var suite models.TestSuite
	if err := yaml.Unmarshal([]byte(yamlContent), &suite); err != nil {
		return nil, err
	}

	return &suite, nil
}
