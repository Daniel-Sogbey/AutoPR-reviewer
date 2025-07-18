package githubapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"review-pr/webhook-service/internal/github"
	"review-pr/webhook-service/internal/requester"
)

func GetPRMetadata(authorization string, prNumber int) (*github.PRMetadataResponseModel, error) {
	headers := http.Header{}
	headers.Add("Accept", "application/vnd.github+json")
	headers.Add("X-GitHub-Api-Version", "2022-11-28")
	prMetadataResponse, err := requester.Requester[github.PRMetadataResponseModel](http.MethodGet, fmt.Sprintf("https://api.github.com/repos/Daniel-Sogbey/code-reviewer/pulls/%d", prNumber), authorization, headers, nil)
	if err != nil {
		return nil, err
	}
	return prMetadataResponse, nil
}

func GetPRChangedFiles(authorization string, prNumber int) (*github.PRFileChangesResponseModel, error) {
	headers := http.Header{}
	headers.Add("Accept", "application/vnd.github+json")
	headers.Add("X-GitHub-Api-Version", "2022-11-28")
	prFileChangesResponse, err := requester.Requester[github.PRFileChangesResponseModel](http.MethodGet, fmt.Sprintf("https://api.github.com/repos/Daniel-Sogbey/code-reviewer/pulls/%d/files", prNumber), authorization, headers, nil)
	if err != nil {
		return nil, err
	}
	return prFileChangesResponse, nil
}

func CreateReviewCommentOnPR(authorization string, prNumber int, requestModel github.PRReviewRequestModel) (*github.PRReviewResponseModel, error) {
	headers := http.Header{}
	headers.Add("Accept", "application/vnd.github+json")
	headers.Add("X-GitHub-Api-Version", "2022-11-28")

	reqBytes, err := json.Marshal(requestModel)
	if err != nil {
		return nil, err
	}

	body := bytes.NewReader(reqBytes)

	prReviewCommentResponse, err := requester.Requester[github.PRReviewResponseModel](http.MethodPost, fmt.Sprintf("https://api.github.com/repos/Daniel-Sogbey/code-reviewer/pulls/%d/reviews", prNumber), authorization, headers, body)
	if err != nil {
		return nil, err
	}

	return prReviewCommentResponse, nil
}
