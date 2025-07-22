package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func router(chunkChan chan Envelope, llmResponseChan chan Envelope, errorChan chan error, engine any) http.Handler {
	r := http.NewServeMux()

	r.HandleFunc("POST /webhook", func(w http.ResponseWriter, r *http.Request) {
		handleWebhook(w, r, chunkChan, llmResponseChan, errorChan, engine)
	})

	r.HandleFunc("POST /init", func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Repo     string `json:"repo"`
			PRNumber string `json:"pr_number"`
			APIKey   string `json:"api_key"`
		}

		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			log.Println("ERROR DECODING INIT CREDENTIALS:", err)
			return
		}

		//credentials := fmt.Sprintf("%s:%s:%s", input.Repo, input.PRNumber, input.APIKey)
		//
		//err = os.Setenv("REQUEST_CREDENTIALS", credentials)
		//if err != nil {
		//	log.Println("ERROR SETTING ENV VALUES: ", err)
		//	return
		//}

		configInstance.AuthKey = input.APIKey

		log.Println("INPUT CREDENTIALS: ", input)

		w.WriteHeader(http.StatusOK)
	})

	return r

}
