package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"time"
)

const successorTemplate = "http://%s.default.10.99.236.6.sslip.io"

type JsonMap map[string]interface{}

func run(command string, args []string) (string, error) {
	fmt.Println("Running command:", command, args)

	output, err := exec.Command(command, args...).Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func successorUrl(successor string) string {
	return fmt.Sprintf(successorTemplate, successor)
}

func autoscalerType(metric string) string {
	if metric == "rps" || metric == "concurrency" {
		return "kpa"
	} else {
		return "hpa"
	}
}

func createK8sDefinition(data JsonMap) (string, error) {
	templ, err := template.ParseFiles("k8s_resource_template.yaml")
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	err = templ.Execute(&buffer, data)

	if err != nil {
		return "", err
	}

	return buffer.String(), nil
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
			successorsStrings[i] = successorUrl(successor.(string))
		}

		fmt.Println("Function type:", config["function_type"].(string))

		config["successors"] = strings.Join(successorsStrings, ",")
		config["function_name"] = functionName
		config["autoscaler_type"] = autoscalerType(config["metric"].(string))

		resourceDefinition, err := createK8sDefinition(config)
		if err != nil {
			panic(err)
		}

		writeToFile("k8s_resource.yaml", resourceDefinition)
		result, err := run("kubectl", []string{"apply", "-f", "k8s_resource.yaml"})
		fmt.Println(result)
		time.Sleep(1000 * time.Millisecond)

		if err != nil {
			panic(err)
		}

	}

	err = os.Remove("k8s_resource.yaml")
	if err != nil {
		panic(err)
	}
}
