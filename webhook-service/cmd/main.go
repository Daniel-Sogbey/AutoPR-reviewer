package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {

	fmt.Println(os.Getenv("GITHUB_APP_PRIVATE_KEY_PATH"))

	srv := http.Server{
		Addr:    ":3000",
		Handler: router(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
