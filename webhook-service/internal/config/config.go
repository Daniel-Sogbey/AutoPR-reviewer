package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type Config struct {
	LLMProvider       string `yaml:"llm_provider"`
	Model             string `yaml:"model"`
	AuthKey           string `yaml:"auth_key"`
	GuidelinesPath    string `yaml:"guidelines"`
	Language          string `yaml:"language,omitempty"`
	GuidelinesContent []byte `yaml:"-"`
}

func LoadConfig(wd string) (*Config, error) {
	configFilePath := filepath.Join(wd, "ai-reviewer.yml")
	file, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err = yaml.Unmarshal(file, &cfg); err != nil {
		return nil, err
	}

	guidelineFilePath := filepath.Join(wd, cfg.GuidelinesPath)
	guidelines, err := os.ReadFile(guidelineFilePath)
	if err != nil {
		return nil, err
	}

	cfg.GuidelinesContent = guidelines

	apiKey := os.Getenv("TOGETHERAI_API_KEY")
	cfg.AuthKey = apiKey

	return &cfg, nil
}
