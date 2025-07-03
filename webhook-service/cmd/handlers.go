package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"review-pr/webhook-service/internal/github"
	"review-pr/webhook-service/internal/githubapi"
)

type Queue struct {
	chunkChan       chan Envelope
	llmResponseChan chan Envelope
	errorChan       chan error
}

type Envelope struct {
	data any
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {
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

	queue := Queue{
		chunkChan:       make(chan Envelope, 2),
		llmResponseChan: make(chan Envelope, 2),
		errorChan:       make(chan error, 5),
	}

	go func() {
		ExtractDiffChunk(queue.chunkChan, queue.errorChan, *prChangedFilesResponse)
		close(queue.chunkChan)
	}()

	go func() {
		QueryLLMWithChunks(queue.llmResponseChan, queue.chunkChan, queue.errorChan, installationTokenResponse.Token, pullRequestEventModel.Number, prMetaDataResponse.Head.Sha, *prChangedFilesResponse)
		close(queue.llmResponseChan)
	}()

	select {
	case llmResponse := <-queue.llmResponseChan:
		log.Println("LLM RESPONSE", llmResponse)
	case err = <-queue.errorChan:
		log.Println("Error in PR REVIEW SYSTEM", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
