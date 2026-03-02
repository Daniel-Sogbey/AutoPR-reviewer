package llmapi

import (
	"review-pr/webhook-service/internal/config"
	"review-pr/webhook-service/internal/github"
)

type Gemini struct {
}

type GeminiRequestModel struct {
}

type GeminiResponseModel struct {
}

func NewGeminiModel() *Gemini {
	return &Gemini{}
}

func (g *Gemini) QueryLLM(in *GeminiRequestModel) (out *GeminiResponseModel, err error) {
	return nil, nil
}

func (g *Gemini) Generate(chunk github.DiffChunk, config config.Config) *GeminiRequestModel {
	return nil
}
