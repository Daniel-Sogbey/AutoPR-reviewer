package llmapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"review-pr/webhook-service/internal/config"
	"review-pr/webhook-service/internal/github"
	"review-pr/webhook-service/internal/requester"
	"strings"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type TogetherAiRequestModel struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type TogetherAiResponseModel struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Prompt  []any  `json:"prompt"`
	Choices []struct {
		FinishReason string      `json:"finish_reason"`
		Seed         json.Number `json:"seed"`
		Logprobs     any         `json:"logprobs"`
		Index        int         `json:"index"`
		Message      struct {
			Role      string `json:"role"`
			Content   string `json:"content"`
			ToolCalls []any  `json:"tool_calls"`
		} `json:"message"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
		CachedTokens     int `json:"cached_tokens"`
	} `json:"usage"`
}

type TogetherAI struct {
	AuthKey string
}

func NewTogetherAI(authKey string) *TogetherAI {
	return &TogetherAI{AuthKey: authKey}
}

func (t TogetherAI) QueryLLM(query *TogetherAiRequestModel) (*TogetherAiResponseModel, error) {
	//authorization := os.Getenv("TOGETHER_AI_LLM_AUTH_KEY")
	headers := http.Header{}
	headers.Add("Content-Type", "application/json")

	bytes, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	response, err := requester.Requester[TogetherAiResponseModel](http.MethodPost, "https://api.together.xyz/v1/chat/completions", t.AuthKey, headers, strings.NewReader(string(bytes)))
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (t TogetherAI) Generate(chunk github.DiffChunk, config config.Config) *TogetherAiRequestModel {

	msg := fmt.Sprintf(`
You are a code reviewer.

# Coding Guidelines:
%s

# Diff Hunk:
%s

## Your task:
- ONLY list issues that violate the guidelines.
- Keep it short.
- If nothing is wrong, say: "No guideline violations found."

Format:
- [ ] Violation: <description>
  - Line: <relevant line>
  - Suggestion: <fix if applicable>
`, config.GuidelinesContent, chunk.CleanedCode)

	return &TogetherAiRequestModel{
		Model: config.Model,
		Messages: []Message{
			{
				Role:    "user",
				Content: msg,
			},
		},
	}
}
