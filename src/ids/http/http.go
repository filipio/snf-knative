package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
)

// type JsonMap map[string]any

func Post(url string, body string) (int, string, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(body)))
	if err != nil {
		log.Fatal(err)
		return 0, "", err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return 0, "", err
	}

	defer resp.Body.Close()

	responseBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return 0, "", err
	}

	return resp.StatusCode, string(responseBodyBytes), nil
}
