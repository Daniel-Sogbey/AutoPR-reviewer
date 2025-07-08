package main

import (
	"log"
	"net/http"
	"os"
	"review-pr/webhook-service/internal/config"
	"review-pr/webhook-service/internal/llmapi"
	"sync"
)

var (
	once           sync.Once
	configInstance *config.Config
)

func getConfigInstance() *config.Config {
	once.Do(func() {

		wd, err := os.Getwd()
		if err != nil {
			log.Println("ERROR GETTING WORKING DIR: ", err)
			os.Exit(1)
		}

		cfg, err := config.LoadConfig(wd)
		if err != nil {
			log.Println("ERROR LOADING CONFIG: ", err)
			os.Exit(1)
		}

		configInstance = cfg
	})

	log.Println("CONFIG INFO: ", configInstance)
	return configInstance
}

func main() {
	cfg := getConfigInstance()

	engine := llmapi.NewLLMEngineRegistry(cfg)

	if engine == nil {
		log.Printf("Unsupported LLM Provider: %s", cfg.LLMProvider)
		return
	}

	chunkChan := make(chan Envelope, 1)
	llmResponseChan := make(chan Envelope, 1)
	errorChan := make(chan error, 5)

	go func() {
		for err := range errorChan {
			log.Println("Received an error: ", err)
		}
	}()

	srv := http.Server{
		Addr:    ":3000",
		Handler: router(chunkChan, llmResponseChan, errorChan, engine),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
