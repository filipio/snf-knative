package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
)

var c = cache.New(5*time.Minute, 10*time.Minute)

type Payload struct {
	ID int `json:"id"`
}

type Response struct {
	Value float64 `json:"value"`
}

func handleError(w http.ResponseWriter, successor string, err string) {
	errMsg := "Error sending POST request to " + successor + " : " + err
	fmt.Println(errMsg)
	http.Error(w, "Error sending POST request to "+successor+" : "+err, http.StatusInternalServerError)
}

func handler(w http.ResponseWriter, r *http.Request) {
	functionName := os.Getenv("FUNCTION_NAME")
	if functionName == "" {
		http.Error(w, "functionName variable not set", http.StatusInternalServerError)
		return
	}

	successors := os.Getenv("SUCCESSORS")

	// split it by comma
	successorsNames := strings.Split(successors, ",")
	fmt.Println("functionName: ", functionName, " successors: ", successorsNames)

	requestBody := "ping"

	if len(successorsNames) > 0 {
		for _, successor := range successorsNames {
			if successor == "" {
				continue
			}
			fmt.Printf("sending request to '%s'", successor)
			// send a POST request to each successor
			if statusCode, response, err := Post(successor, requestBody); err != nil {
				handleError(w, successor, err.Error())
			} else {
				if statusCode != 200 {
					handleError(w, successor, response)
					return
				}
				fmt.Println("Response from ", successor, ": ", response)
			}
		}
	}

	// return 200 with "success" response
	w.Write([]byte("success"))

	// var p Payload
	// err := json.NewDecoder(r.Body).Decode(&p)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	// val, found := c.Get(string(p.ID))
	// if !found {
	// 	val = rand.Float64()
	// 	c.Set(string(p.ID), val, cache.DefaultExpiration)
	// }

	// resp := Response{Value: val.(float64)}
	// json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
