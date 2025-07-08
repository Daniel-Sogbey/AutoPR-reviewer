package llmapi

import (
	"encoding/json"
	"net/http"
	"review-pr/webhook-service/internal/config"
	"review-pr/webhook-service/internal/github"
	"review-pr/webhook-service/internal/requester"
	"strings"
)

type OpenAI struct {
	AuthKey string
}

func NewOpenAi(authKey string) *OpenAI {
	return &OpenAI{AuthKey: authKey}
}

type OpenAiRequestModel struct {
}

type OpenAiResponseModel struct {
}

func (o OpenAI) QueryLLM(query *OpenAiRequestModel) (*OpenAiResponseModel, error) {
	headers := http.Header{}
	headers.Add("Content-Type", "application/json")

	bytes, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	response, err := requester.Requester[OpenAiResponseModel](http.MethodPost, "https://api.together.xyz/v1/chat/completions", o.AuthKey, headers, strings.NewReader(string(bytes)))
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (t OpenAI) Generate(chunk github.DiffChunk, config config.Config) *OpenAiRequestModel {

	//	msg := fmt.Sprintf(`
	//You are a code reviewer.
	//
	//# Coding Guidelines:
	//%s
	//
	//# Diff Hunk:
	//%s
	//
	//## Your task:
	//- ONLY list issues that violate the guidelines.
	//- Keep it short.
	//- If nothing is wrong, say: "No guideline violations found."
	//
	//Format:
	//- [ ] Violation: <description>
	//  - Line: <relevant line>
	//  - Suggestion: <fix if applicable>
	//`, config.GuidelinesContent, chunk.CleanedCode)

	return &OpenAiRequestModel{}
}
