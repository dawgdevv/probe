package loader

import (
	"os"

	"github.com/dawgdevv/apitestercli/pkg/models"
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
