package main

import "net/http"

func router() http.Handler {
	r := http.NewServeMux()

	r.HandleFunc("POST /", handleWebhook)

	return r
}
