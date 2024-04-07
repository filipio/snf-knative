package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
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

func handler(w http.ResponseWriter, r *http.Request) {
	var p Payload
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	val, found := c.Get(string(p.ID))
	if !found {
		val = rand.Float64()
		c.Set(string(p.ID), val, cache.DefaultExpiration)
	}

	resp := Response{Value: val.(float64)}
	json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
