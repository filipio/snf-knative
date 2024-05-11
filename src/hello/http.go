package function

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type JsonMap map[string]any

func Post(url string, jsonMap JsonMap) (JsonMap, error) {

	jsonBytes, err := json.Marshal(jsonMap)
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer resp.Body.Close()

	responseBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var responseBodyJson JsonMap
	err = json.Unmarshal(responseBodyBytes, &responseBodyJson)
	if err != nil {
		log.Fatalf("Error unmarshaling response body: %v", err)
		return nil, err
	}

	return responseBodyJson, nil
}
