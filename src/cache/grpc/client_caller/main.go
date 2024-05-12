package main

import (
	"fmt"

	"github.com/filipio/snf-knative/src/cache/grpc/client"
)

const functionType = "f1"

func main() {
	serverAddr := functionType + ".default.10.99.236.6.sslip.io:80"

	message, err := client.Call(serverAddr, "hello")
	if err != nil {
		panic(err)
	}

	fmt.Println(message)
}
