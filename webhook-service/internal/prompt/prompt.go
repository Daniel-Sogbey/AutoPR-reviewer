package prompt

import (
	"errors"
	"fmt"
	"os"
	"review-pr/webhook-service/internal/github"
	"review-pr/webhook-service/internal/llm"
)

var Prompt = func(chunk github.DiffChunk) (llm.TogetherAiRequestModel, error) {

	guidelines, err := os.ReadFile("/Users/danielsogbey/code/personal/go/review-pr/CODING_GUIDELINES.md")
	if err != nil {
		return llm.TogetherAiRequestModel{}, errors.New("failed to load guideline files")
	}

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
`, guidelines, chunk.CleanedCode)

	return llm.TogetherAiRequestModel{
		Model: "meta-llama/Meta-Llama-3.1-8B-Instruct-Turbo",
		Messages: []llm.Message{
			{
				Role:    "user",
				Content: msg,
			},
		},
	}, nil
}
