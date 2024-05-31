package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
)

// var solutions = []string{"knative"}
// var functionTypes = []string{"ids"}
// var protocols = []string{"http"}
// var graphNumbers = []int{1}
// var workloads = []string{"weibull"}

var protocols = []string{"http"}
var functionTypes = []string{"cache"}
var graphNumbers = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
var workloads = []string{"weibull", "poisson", "pareto"}
var solutions = []string{"knative", "deployment"}
var politics = []string{"reactive", "resource-optimized"}

const eventsCount = 100
const workers = 10

// weibull
const weibullShape = 1.5
const weibullScale = 50.0

// poisson
const poissonLambda = 50.0

// pareto
const paretoXmin = 10.0
const paretoAlpha = 2.0

type JsonMap map[string]interface{}

var defaultConfig = JsonMap{
	"protocol":         "grpc",
	"function_type":    "cache",
	"graph_number":     10,
	"metric":           "rps",
	"min_scale":        "1",
	"max_scale":        "5",
	"scale_down_delay": "15m",
	"target":           5,
	"window":           "40s",
	"min_cpu":          "100m",
	"max_cpu":          "300m",
	"min_memory":       "40Mi",
	"max_memory":       "400Mi",
}

func readJson(fileName string) JsonMap {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var config JsonMap
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalf("Failed to decode file: %s", err)
	}

	return config
}

func saveJson(fileName string, content JsonMap) {
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("Failed to create file: %s", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(content)
	if err != nil {
		log.Fatalf("Failed to encode file: %s", err)
	}
}

func outputFormat(functionType string, workload string, solution string, protocol string, graphNumber int, policyName string) string {
	return fmt.Sprintf("results_csv/%s/%s/%s/%s/%s/%d", policyName, functionType, workload, solution, protocol, graphNumber)
}

func main() {
	var host = flag.String("host", "localhost", "destination host address")
	flag.Parse()

	policies := readJson("policies.json")["policies"].([]interface{})

	for _, policyInterface := range policies {
		policy := policyInterface.(map[string]interface{})
		policyName := policy["name"].(string)
		foundPolicy := false
		for _, validPolicies := range politics {
			if policyName == validPolicies {
				foundPolicy = true
				break
			}
		}
		if !foundPolicy {
			log.Println("Policy not found in politics, skipping...")
			continue
		}
		config := defaultConfig
		for key, value := range policy {
			config[key] = value
		}

		for _, solution := range solutions {
			for _, funcType := range functionTypes {
				for _, protocol := range protocols {
					for _, graphNumber := range graphNumbers {
						for _, workload := range workloads {
							config["function_type"] = funcType
							config["protocol"] = protocol
							config["graph_number"] = graphNumber

							saveJson("config.json", config)
							setupInfra()
							output := outputFormat(funcType, workload, solution, protocol, graphNumber, policyName)
							fmt.Println("output is", output)
							GenerateTraffic(workers, *host, workload, eventsCount, output,
								protocol,
								weibullShape,
								weibullScale,
								poissonLambda,
								paretoXmin,
								paretoAlpha)
						}
					}
				}

			}
		}

	}

}

// what do i want to change:
// metric, min_scale, max_scale, target, window, ...
// too many things which change
// hmm

// not changing (or all values need to be tested):
// workloads (3 types) and params
// traffic (2 types)
// function (2 types)
//

// is it worth to explore the parameter space?
// then it would be like a training, and this is related to ML
// so I think I can just define several policies which will be suited for some kind of an app
// what policies?
// longer window + higher target
// smaller window + lower target
// cpu
// number of requests
// so there could be 4 policies
