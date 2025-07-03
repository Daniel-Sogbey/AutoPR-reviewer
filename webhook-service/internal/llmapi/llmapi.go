package llmapi

import (
	"encoding/json"
	"net/http"
	"os"
	"review-pr/webhook-service/internal/llm"
	"review-pr/webhook-service/internal/requester"
	"strings"
)

func QueryLMM(query llm.TogetherAiRequestModel) (*llm.TogetherAiResponseModel, error) {
	authorization := os.Getenv("TOGETHER_AI_LLM_AUTH_KEY")
	headers := http.Header{}
	headers.Add("Content-Type", "application/json")

	bytes, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	response, err := requester.Requester[llm.TogetherAiResponseModel](http.MethodPost, "https://api.together.xyz/v1/chat/completions", authorization, headers, strings.NewReader(string(bytes)))
	if err != nil {
		return nil, err
	}

	return response, nil
}
