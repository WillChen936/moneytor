package utils

import (
	"encoding/json"
	"os"
)

type Config struct {
	DBSource string `json:"db_source"`
}

func LoadConfig(path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}

	defer file.Close()

	var config Config
	err = json.NewDecoder(file).Decode(&config)

	return config, err
}
