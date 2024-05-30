package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type JsonMap map[string]any

func Post(url string, body string) (int, JsonMap, error) {
	// encode into json
	mapBody := JsonMap{"body": body}
	jsonBody, err := json.Marshal(mapBody)
	if err != nil {
		log.Fatal(err)
		return 0, JsonMap{}, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Fatal(err)
		return 0, JsonMap{}, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return 0, JsonMap{}, err
	}

	defer resp.Body.Close()

	responseBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return 0, JsonMap{}, err
	}

	// unmarshal
	var responseJson JsonMap
	err = json.Unmarshal(responseBodyBytes, &responseJson)
	if err != nil {
		log.Fatal(err)
		return 0, JsonMap{}, err
	}

	return resp.StatusCode, responseJson, nil
}
