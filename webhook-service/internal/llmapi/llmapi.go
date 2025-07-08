package llmapi

import (
	"review-pr/webhook-service/internal/config"
	"review-pr/webhook-service/internal/github"
)

type LLMQuery[I any, O any] interface {
	QueryLLM(input I) (O, error)
}

type PromptGenerator[I any] interface {
	Generate(chunk github.DiffChunk, config config.Config) I
}

type LLMEngine[I any, O any] struct {
	Prompt PromptGenerator[I]
	Query  LLMQuery[I, O]
}
