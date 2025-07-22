package main

import (
	"net/http"
)

func router(chunkChan chan Envelope, llmResponseChan chan Envelope, errorChan chan error, engine any) http.Handler {
	r := http.NewServeMux()

	r.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
		handleWebhook(w, r, chunkChan, llmResponseChan, errorChan, engine)
	})

	return r

}
