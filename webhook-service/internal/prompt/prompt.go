package prompt

import (
	"fmt"
	"review-pr/webhook-service/internal/github"
	"review-pr/webhook-service/internal/llmapi"
)

var Prompt = func(llm llmapi.LMMQuery, cfgGuidelines []byte, chunk github.DiffChunk) (*llmapi.LLMRequest, error) {

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
`, cfgGuidelines, chunk.CleanedCode)

	return nil, nil
}

//llm.TogetherAiRequestModel{
//Model: "meta-llama/Meta-Llama-3.1-8B-Instruct-Turbo",
//Messages: []llm.Message{
//{
//Role:    "user",
//Content: msg,
//},
//},
//}
