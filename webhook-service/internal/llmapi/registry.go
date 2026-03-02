package llmapi

import (
	"log"
	"review-pr/webhook-service/internal/config"
)

func NewLLMEngineRegistry(cfg *config.Config) any {
	switch cfg.LLMProvider {
	case "openai":
		return LLMEngine[*OpenAiRequestModel, *OpenAiResponseModel]{
			Prompt: NewOpenAi(cfg.AuthKey),
			Query:  NewOpenAi(cfg.AuthKey),
		}
	case "togetherai":
		return LLMEngine[*TogetherAiRequestModel, *TogetherAiResponseModel]{
			Prompt: NewTogetherAI(cfg.AuthKey),
			Query:  NewTogetherAI(cfg.AuthKey),
		}
	case "gemini":
		return LLMEngine[*GeminiRequestModel, *GeminiResponseModel]{
			Prompt: NewGeminiModel(),
			Query:  NewGeminiModel(),
		}
	default:
		log.Printf("Unsupported LLM Provider: %s", cfg.LLMProvider)
		return nil
	}
}
