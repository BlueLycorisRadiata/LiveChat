package config

import "os"

type Config struct {
	OpenRouterAPIKey  string
	OpenRouterBaseURL string
	DefaultAIModel    string
}

func Load() *Config {
	baseURL := os.Getenv("OPENROUTER_BASE_URL")
	if baseURL == "" {
		baseURL = "https://openrouter.ai/api/v1"
	}

	model := os.Getenv("OPENROUTER_DEFAULT_MODEL")
	if model == "" {
		model = "nvidia/nemotron-3-super-120b-a12b:free"
	}

	return &Config{
		OpenRouterAPIKey:  os.Getenv("OPENROUTER_API_KEY"),
		OpenRouterBaseURL: baseURL,
		DefaultAIModel:    model,
	}
}
