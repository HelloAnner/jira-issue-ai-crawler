package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Jira struct {
		URL      string `yaml:"url"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		JQL      string `yaml:"jql"`
	} `yaml:"jira"`

	AI struct {
		APIKey      string  `yaml:"api_key"`
		Model       string  `yaml:"model"`
		BaseURL     string  `yaml:"base_url"`
		Temperature float64 `yaml:"temperature"`
	} `yaml:"ai"`

	DB struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
	} `yaml:"db"`

	Sync struct {
		Interval int `yaml:"interval"` // in minutes
	} `yaml:"sync"`
}

func LoadConfig(path string) (*Config, error) {
	config := &Config{}

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(file, config)
	if err != nil {
		return nil, err
	}

	return config, nil
} 