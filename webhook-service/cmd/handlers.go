package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"review-pr/webhook-service/internal/github"
	"review-pr/webhook-service/internal/githubapi"
	"review-pr/webhook-service/internal/llmapi"
	"time"
)

type Queue struct {
	chunkChan       chan Envelope
	llmResponseChan chan Envelope
	errorChan       chan error
}

type Envelope struct {
	data any
}

func handleWebhook(w http.ResponseWriter, r *http.Request, chunkChan chan Envelope, llmResponseChan chan Envelope, errorChan chan error, engine any) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("error reading request body: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	signature := r.Header.Get("X-Hub-Signature-256")
	if !VerifySignature(bodyBytes, "danielsogbeykeys", signature) {
		log.Println("invalid signature")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var pullRequestEventModel github.PullRequestEventModel
	if err = json.Unmarshal(bodyBytes, &pullRequestEventModel); err != nil {
		log.Println("JSON parse err:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if pullRequestEventModel.Action != "opened" && pullRequestEventModel.Action != "synchronize" && pullRequestEventModel.Action != "reopened" {
		log.Println("Ignoring action:", pullRequestEventModel.Action)
		return
	}

	installationTokenResponse, err := githubapi.Auth()
	if err != nil {
		log.Println("1:err", err)
		return
	}

	prChangedFilesResponse, err := githubapi.GetPRChangedFiles(installationTokenResponse.Token, pullRequestEventModel.Number)
	if err != nil {
		log.Println("2:err", err)
		return
	}

	prMetaDataResponse, err := githubapi.GetPRMetadata(installationTokenResponse.Token, pullRequestEventModel.Number)
	if err != nil {
		log.Println("3:err", err)
		return
	}

	fmt.Println("event result:", pullRequestEventModel)
	fmt.Println("pr changed file response:", prChangedFilesResponse)
	fmt.Println("pr meta data response:", prMetaDataResponse)

	go ExtractDiffChunk(chunkChan, errorChan, *prChangedFilesResponse)

	switch llm := engine.(type) {
	case *llmapi.LLMEngine[*llmapi.TogetherAiRequestModel, *llmapi.TogetherAiResponseModel]:
		go QueryLLMWithChunks(llm, llmResponseChan, chunkChan, errorChan, installationTokenResponse.Token, pullRequestEventModel.Number, prMetaDataResponse.Head.Sha, *prChangedFilesResponse)
	case *llmapi.LLMEngine[*llmapi.OpenAiRequestModel, *llmapi.OpenAiResponseModel]:
		go QueryLLMWithChunks(llm, llmResponseChan, chunkChan, errorChan, installationTokenResponse.Token, pullRequestEventModel.Number, prMetaDataResponse.Head.Sha, *prChangedFilesResponse)
	default:
		errorChan <- fmt.Errorf("unsupported LLM Engine")
		return
	}
	
	select {
	case llmResponse := <-llmResponseChan:
		log.Println("LLM RESPONSE RECEIVED", llmResponse)
	case <-time.After(30 * time.Second):
		log.Println("Timeout llm taking too much time to return response")
		w.WriteHeader(http.StatusGatewayTimeout)
	}

}
