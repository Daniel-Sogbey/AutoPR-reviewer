package requester

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func Requester[T any](method, url, authorization string, headers http.Header, body io.Reader) (*T, error) {
	headers.Add("Authorization", fmt.Sprintf("Bearer %s", authorization))

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	client := http.Client{}

	req.Header = headers

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	log.Println("Raw response:", string(responseBytes))

	var t T
	err = json.Unmarshal(responseBytes, &t)
	if err != nil {
		return nil, err
	}

	return &t, nil
}
