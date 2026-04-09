package config

import "os"

var SupportedAIModels = []string{
	"nvidia/nemotron-nano-9b-v2:free",
	"minimax/minimax-m2.5:free",
	"qwen/qwen3.6-plus:free",
	"nvidia/nemotron-3-super-120b-a12b:free",
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

	corsOrigins := []string{"http://localhost:5173"}
	if originsEnv := os.Getenv("CORS_ORIGINS"); originsEnv != "" {
		corsOrigins = []string{originsEnv}
	}

	return &Config{
		OpenRouterAPIKey:  os.Getenv("OPENROUTER_API_KEY"),
		OpenRouterBaseURL: baseURL,
		DefaultAIModel:    DefaultAIModel,
		CORSOrigins:       corsOrigins,
	}
}
