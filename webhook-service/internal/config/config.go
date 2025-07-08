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
	GuidelinesPath    string `json:"guidelines"`
	Language          string `yaml:"language,omitempty"`
	GuidelinesContent []byte
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

	guidelineFilePath := filepath.Join(wd, "CODING_GUIDELINES.md")
	guidelines, err := os.ReadFile(guidelineFilePath)
	if err != nil {
		return nil, err
	}

	cfg.GuidelinesContent = guidelines

	return &cfg, nil
}
