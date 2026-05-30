package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type rawConfig struct {
	Env                  string `json:"Env"`
	DBSource             string `json:"DBSource"`
	HttpServerAddress    string `json:"HttpServerAddress"`
	TokenSecretKey       string `json:"TokenSecretKey"`
	AccessTokenDuration  string `json:"AccessTokenDuration"`
	RefreshTokenDuration string `json:"RefreshTokenDuration"`
}

type Config struct {
	Env                  string
	DBSource             string
	HttpServerAddress    string
	TokenSecretKey       string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

func LoadConfig(path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	var raw rawConfig
	if err := json.NewDecoder(file).Decode(&raw); err != nil {
		return Config{}, err
	}

	accessDuration, err := time.ParseDuration(raw.AccessTokenDuration)
	if err != nil {
		return Config{}, fmt.Errorf("invalid AccessTokenDuration: %w", err)
	}

	refreshDuration, err := time.ParseDuration(raw.RefreshTokenDuration)
	if err != nil {
		return Config{}, fmt.Errorf("invalid RefreshTokenDuration: %w", err)
	}

	return Config{
		Env:                  raw.Env,
		DBSource:             raw.DBSource,
		HttpServerAddress:    raw.HttpServerAddress,
		TokenSecretKey:       raw.TokenSecretKey,
		AccessTokenDuration:  accessDuration,
		RefreshTokenDuration: refreshDuration,
	}, nil
}
