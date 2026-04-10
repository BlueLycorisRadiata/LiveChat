package config

import (
	"os"
	"strings"
)

var SupportedAIModels = []string{
	"nvidia/nemotron-nano-9b-v2:free",
	"nvidia/nemotron-3-super-120b-a12b:free",
	"google/gemma-4-26b-a4b-it:free",
	"arcee-ai/trinity-large-preview:free",
	"openai/gpt-oss-120b:free",
}

var DefaultAIModel = "nvidia/nemotron-3-super-120b-a12b:free"

type Config struct {
	OpenRouterAPIKey  string
	OpenRouterBaseURL string
	DefaultAIModel    string
	CORSOrigins       []string
}

func IsModelSupported(model string) bool {
	for _, m := range SupportedAIModels {
		if m == model {
			return true
		}
	}
	return false
}

func GetDefaultAIModel() string {
	return DefaultAIModel
}

func Load() *Config {
	baseURL := os.Getenv("OPENROUTER_BASE_URL")
	if baseURL == "" {
		baseURL = "https://openrouter.ai/api/v1"
	}

	corsOrigins := []string{
		"http://localhost:5173",
		"https://live-chat-frontend-sage.vercel.app",
	}

	if originsEnv := os.Getenv("CORS_ORIGINS"); originsEnv != "" {
		parts := strings.Split(originsEnv, ",")
		corsOrigins = make([]string, 0, len(parts))
		for _, origin := range parts {
			origin = strings.TrimSpace(origin)
			if origin != "" {
				corsOrigins = append(corsOrigins, origin)
			}
		}
	}

	return &Config{
		OpenRouterAPIKey:  os.Getenv("OPENROUTER_API_KEY"),
		OpenRouterBaseURL: baseURL,
		DefaultAIModel:    DefaultAIModel,
		CORSOrigins:       corsOrigins,
	}
}
