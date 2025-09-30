package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	configFileName = "/.gatorconfig.json"
)

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	current_directory, _ := os.UserHomeDir() // replace with os.Getwd() to use .gatorconfig.json in Gator directory
	filepath := current_directory + configFileName
	configFile, err := os.Open(filepath)
	if err != nil {
		return Config{}, fmt.Errorf("error opening config file: %w", err)
	}
	defer configFile.Close()
	config := Config{}
	decoder := json.NewDecoder(configFile)
	if err = decoder.Decode(&config); err != nil {
		return Config{}, fmt.Errorf("error parsing config file: %w", err)
	}

	return config, nil
}

func (c *Config) SetUser(username string) error {
	c.CurrentUserName = username
	if err := write(*c); err != nil {
		return fmt.Errorf("error setting user: %w", err)
	}
	return nil
}

func write(config Config) error {
	current_directory, _ := os.UserHomeDir() // replace with os.Getwd() to use .gatorconfig.json in Gator directory
	filepath := current_directory + configFileName
	configFile, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("error opening config file: %w", err)
	}
	defer configFile.Close()
	encoder := json.NewEncoder(configFile)
	if err = encoder.Encode(&config); err != nil {
		return fmt.Errorf("error writing config to file: %w", err)
	}
	return nil
}
