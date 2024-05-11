package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type JsonMap map[string]interface{}

func run(command string, args []string) (string, error) {
	fmt.Println("Running command:", command, args)

	output, err := exec.Command(command, args...).Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func cacheK8sDefinition(imageType string, functionName string, successors []string) string {
	return fmt.Sprintf(`
apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  name: %s
  namespace: default
spec:
  template:
    spec:
      containers:
        - image: docker.io/notnew77/%s:latest
          ports:
          - containerPort: 8080
          env:
          - name: FUNCTION_NAME
            value: "%s"
          - name: SUCCESSORS
            value: "%s"`, functionName, imageType, functionName, strings.Join(successors, ","))
}

func writeToFile(filename string, data string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(data)
	if err != nil {
		return err
	}

	return nil
}

func readJsonConfig(path string) (JsonMap, error) {
	// read config.json file
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var jsonMap JsonMap

	err = json.Unmarshal(fileBytes, &jsonMap)
	if err != nil {
		return nil, err
	}

	return jsonMap, nil
}

func main() {
	fmt.Println("deleting all ksvc resources in default namespace...")
	result, _ := run("kubectl", []string{"delete", "ksvc", "--all", "-n", "default"})
	fmt.Println(result)

	config, err := readJsonConfig("config.json")
	if err != nil {
		panic(err)
	}

	fmt.Println(config)

	graphNumber := int(config["graph_number"].(float64))

	graphConfig, graphConfigError := readJsonConfig(fmt.Sprintf("config/graphs/%d.json", graphNumber))
	if graphConfigError != nil {
		panic(graphConfigError)
	}

	graph := graphConfig["graph"].(map[string]interface{})

	for key, value := range graph {
		functionName := key
		successors := value.([]interface{})

		successorsStrings := make([]string, len(successors))
		for i, successor := range successors {
			successorsStrings[i] = successor.(string)
		}

		functionType := config["function_type"].(string)
		fmt.Println("Function type:", functionType)

		// fmt.Println(cacheK8sDefinition(functionName, successorsStrings))

		writeToFile("k8s_resource.yaml", cacheK8sDefinition(functionType, functionName, successorsStrings))
		result, err := run("kubectl", []string{"apply", "-f", "k8s_resource.yaml"})
		fmt.Println(result)

		// sleep for 1 seconds
		time.Sleep(1000 * time.Millisecond)

		if err != nil {
			panic(err)
		}

	}

	// delete file k8s_resource.yaml
	err = os.Remove("k8s_resource.yaml")
	if err != nil {
		panic(err)
	}
}
