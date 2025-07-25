package main

import (
	"fmt"
	"github.com/sourcegraph/go-diff/diff"
	"log"
	"review-pr/webhook-service/internal/github"
	"review-pr/webhook-service/internal/githubapi"
	"review-pr/webhook-service/internal/llmapi"
	"strings"
)

func ExtractDiffChunk(chunkChan chan<- Envelope, errorChan chan error, prFiles github.PRFileChangesResponseModel) {
	var chunks []github.DiffChunk

	for _, file := range prFiles {
		if file.Patch == "" {
			continue
		}

		fullPatch := fmt.Sprintf("--- a/%s\n+++ b/%s\n%s", file.Filename, file.Filename, file.Patch)
		fileDiff, err := diff.ParseFileDiff([]byte(fullPatch))
		if err != nil {
			errorChan <- fmt.Errorf("PARSING FILE DIFF ERROR %v", err)
		}

		log.Println("fileDiff:", fileDiff)

		if fileDiff != nil {
			for _, hunk := range fileDiff.Hunks {
				rawChunk := string(hunk.Body)
				log.Println("rawChunk:", rawChunk)
				cleaned := cleanChunk(rawChunk)
				log.Println("cleaned:", cleaned)

				chunks = append(chunks, github.DiffChunk{
					FilePath:     file.Filename,
					CleanedCode:  cleaned,
					OriginalDiff: rawChunk,
					HunkHeader: fmt.Sprintf(
						"@@ -%d,%d +%d,%d @@ %s",
						hunk.OrigStartLine,
						hunk.OrigLines,
						hunk.NewStartLine,
						hunk.NewLines,
						hunk.Section,
					),
					HunkStartLine: int(hunk.NewStartLine),
				})
			}
		}
	}

	chunkChan <- Envelope{
		data: chunks,
	}
}

func cleanChunk(raw string) string {
	var lines []string

	for _, line := range strings.Split(raw, "\n") {
		if strings.HasPrefix(line, "+") {
			lines = append(lines, strings.TrimPrefix(line, "+"))
		} else if strings.HasPrefix(line, "-") {
			lines = append(lines, "// removed: "+strings.TrimPrefix(line, "-"))
		} else if !strings.HasPrefix(line, "@@") && !strings.HasPrefix(line, "\\") {
			lines = append(lines, line)
		}

	}

	return strings.Join(lines, "\n")
}

func getPosition(patch string) int {
	position := 0
	lines := strings.Split(patch, "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "@@") {
			continue
		}
		position++

		if strings.HasPrefix(line, "+") {
			break
		}
	}

	return position
}

func QueryLLMWithChunks[I, O any](engine *llmapi.LLMEngine[I, O], llmResponseChan chan<- Envelope, chunkChan <-chan Envelope, errorChan chan error, authorization string, prNumber int, commitSHA string, prFiles github.PRFileChangesResponseModel) {
	for msg := range chunkChan {
		var llmResponse O
		var err error
		for _, chunk := range msg.data.([]github.DiffChunk) {
			llmRequest := engine.Prompt.Generate(chunk, *configInstance)

			llmResponse, err = engine.Query.QueryLLM(llmRequest)
			if err != nil {
				errorChan <- fmt.Errorf("LLM RESPONSE ERROR: %v", err)
				return
			}

			content := extractContent(llmResponse)
			if strings.Contains(content, "No guideline violations found.") {
				continue
			}

			log.Println("llm response:", llmResponse)

			prCommentRequestModel := github.PRReviewRequestModel{
				CommitID: commitSHA,
				Body:     "Automated LLM review based on coding guidelines.",
				Event:    "COMMENT",
				Comments: []github.Comment{
					{
						Position: int64(getPosition(chunk.OriginalDiff)),
						Body:     content,
						Path:     chunk.FilePath,
					},
				},
			}

			//push to pr code review
			var reviewCommentOnPR *github.PRReviewResponseModel
			reviewCommentOnPR, err = githubapi.CreateReviewCommentOnPR(authorization, prNumber, prCommentRequestModel)
			if err != nil {
				errorChan <- fmt.Errorf("REVIEW COMMENT ON PR ERROR: %v", err)
			}

			log.Println("REVIEW COMMENT ON PR:", reviewCommentOnPR)

		}

		llmResponseChan <- Envelope{data: llmResponse}
	}
}

func extractContent(response any) string {
	switch r := response.(type) {
	case *llmapi.TogetherAiResponseModel:
		if len(r.Choices) > 0 {
			return r.Choices[0].Message.Content
		}
	case *llmapi.OpenAiResponseModel:

	}

	return ""
}
