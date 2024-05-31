package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

	pb "github.com/filipio/snf-knative/experiments/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var httpClient = &http.Client{}
var message = "ping"
var grpcOpts = []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
var ctx = context.Background()

func Post(url string, body string) (int, string, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(body)))
	if err != nil {
		log.Fatal(err)
		return 0, "", err
	}

	resp, err := httpClient.Do(req)
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

// TODO: implement for gRPC

// generateWeibull generates a Weibull distributed random variable.
func generateWeibull(shape, scale float64) float64 {
	u := rand.Float64() // Generate a U(0,1) random variable
	return scale * math.Pow(-math.Log(1-u), 1/shape)
}

// generateWeibullEvents generates multiple Weibull distributed event times.
func generateWeibullEvents(shape float64, scale float64, size int) []float64 {
	events := make([]float64, size)
	for i := range events {
		events[i] = generateWeibull(shape, scale)
	}
	return events
}

func generatePoissonProcess(lambda float64, numEvents int) []float64 {
	events := make([]float64, numEvents)
	for i := range events {
		events[i] = rand.ExpFloat64() / lambda * 1000 // multiply by 1000 to convert to milliseconds, so lambda = 50 means in avg 50 events per second
	}
	return events
}

func generatePareto(xmin, alpha float64) float64 {
	u := rand.Float64() // Generate a U(0,1) random variable
	return xmin / math.Pow(u, 1/alpha)
}

func generateParetoEvents(xmin, alpha float64, size int) []float64 {
	events := make([]float64, size)
	for i := range events {
		events[i] = generatePareto(xmin, alpha)
	}
	return events
}

func responseTimeWorker(workerID int, outputDir string, workersChannel chan time.Duration, wg *sync.WaitGroup) {
	file, err := os.Create(fmt.Sprintf("%s/results_%d.csv", outputDir, workerID))
	if err != nil {
		log.Fatalf("Failed creating file: %s", err)
	}

	writer := csv.NewWriter(file)
	csvHeader := []string{"response_time"}
	writer.Write(csvHeader)

	for responseTime := range workersChannel {
		// print in ms
		fmt.Println("Response time: ", responseTime.Milliseconds())
		writer.Write([]string{fmt.Sprintf("%d", responseTime.Milliseconds())})
	}

	writer.Flush()
	file.Close()
	wg.Done()
}

func httpWorker(host string, sleepTimeInMilliseconds []float64, workersChannel chan time.Duration) {
	for _, sleepTime := range sleepTimeInMilliseconds {
		fmt.Printf("sleeping for %f ms\n", sleepTime)

		time.Sleep(time.Duration(sleepTime) * time.Millisecond)

		httpHost := fmt.Sprintf("http://%s", host)
		start := time.Now()
		statusCode, _, err := Post(httpHost, message)
		elapsed := time.Since(start)
		if err != nil {
			log.Fatal(err)
		}
		if statusCode != 200 {
			log.Fatalf("Received status code %d", statusCode)
		}

		fmt.Println("passing elapsed time to channel")

		workersChannel <- elapsed

		// fmt.Printf("Sending request to %s\n", host)
	}
	close(workersChannel)
}

func grpcWorker(host string, sleepTimeInMilliseconds []float64, workersChannel chan time.Duration) {
	conn, err := grpc.Dial(host, grpcOpts...) // format of host is f1.default.svc.cluster.local:80
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	// fmt.Println(host)
	client := pb.NewHandlerClient(conn)

	for _, sleepTime := range sleepTimeInMilliseconds {
		fmt.Printf("sleeping for %f ms\n", sleepTime)
		time.Sleep(time.Duration(sleepTime) * time.Millisecond)
		start := time.Now()
		_, err := client.Handle(ctx, &pb.HandleRequest{Message: message}, grpc.EmptyCallOption{})
		if err != nil {
			log.Fatalf("Failed to handle: %v", err)
		}
		elapsed := time.Since(start)

		workersChannel <- elapsed

		// fmt.Printf("Sending request to %s\n", host)
	}
	conn.Close()
	close(workersChannel)
}

func createPathIfNotExists(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Printf("Path %s does not exist, creating...\n", path)
		err := os.MkdirAll(path, 0755)
		if err != nil {
			log.Fatalf("Failed to create directory: %s", err)
		}
	}
}

func waitForHttpReady(host string) error {
	httpHost := fmt.Sprintf("http://%s", host)

	for i := 0; i < 180; i++ {
		statusCode, _, err := Post(httpHost, message)
		if err == nil && statusCode == 200 {
			fmt.Println("HTTP server is ready")
			return nil
		}
		fmt.Println("Waiting for HTTP server to be ready... sleeping for 1 second")
		time.Sleep(1 * time.Second)
	}

	return errors.New("HTTP server is not ready")
}

func waitForGrpcReady(host string) error {
	conn, err := grpc.Dial(host, grpcOpts...)
	if err != nil {
		return err
	}
	defer conn.Close()
	client := pb.NewHandlerClient(conn)

	for i := 0; i < 180; i++ {
		_, grpcCallError := client.Handle(ctx, &pb.HandleRequest{Message: message}, grpc.EmptyCallOption{})
		if grpcCallError == nil {
			fmt.Println("gRPC server is ready")
			return nil
		}
		fmt.Println("Waiting for gRPC server to be ready... sleeping for 1 second")
		time.Sleep(1 * time.Second)
	}
	return errors.New("gRPC server is not ready")
}

func GenerateTraffic(numberOfWorkers int, host string, workloadType string, numberOfEvents int, outputDir string, trafficType string, weibullShape float64, weibullScale float64, poissonLambda float64, paretoXmin float64, paretoAlpha float64) {
	createPathIfNotExists(outputDir)
	host = host + ":80"

	fmt.Printf(" number of workers: %d\n destination host: %s\n workload type: %s\n number of events: %d\n\n", numberOfWorkers, host, workloadType, numberOfEvents)

	var sleepTimeInMilliseconds [][]float64 = make([][]float64, numberOfWorkers)

	for i := 0; i < numberOfWorkers; i++ {
		if workloadType == "weibull" {
			if i == 0 {
				fmt.Printf(" weibull shape: %f\n weibull scale: %f\n", weibullShape, weibullScale)
			}
			sleepTimeInMilliseconds[i] = generateWeibullEvents(weibullShape, weibullScale, numberOfEvents)
		} else if workloadType == "poisson" {
			if i == 0 {
				fmt.Printf(" poisson lambda: %f\n", poissonLambda)
			}
			sleepTimeInMilliseconds[i] = generatePoissonProcess(poissonLambda, numberOfEvents)
		} else if workloadType == "pareto" {
			if i == 0 {
				fmt.Printf(" pareto xmin: %f\n pareto alpha: %f\n", paretoXmin, paretoAlpha)
			}
			sleepTimeInMilliseconds[i] = generateParetoEvents(paretoXmin, paretoAlpha, numberOfEvents)
		} else {
			log.Fatalf("Invalid workload type: %s", workloadType)
		}
	}

	// fmt.Println(sleepTimeInMilliseconds)

	if trafficType == "http" {
		if err := waitForHttpReady(host); err != nil {
			log.Fatalf("%s", err)
		}
	} else if trafficType == "grpc" {
		if err := waitForGrpcReady(host); err != nil {
			log.Fatalf("%s", err)
		}
	}

	var wg sync.WaitGroup = sync.WaitGroup{}
	for i := 0; i < numberOfWorkers; i++ {
		wg.Add(1)
		var workersChannel chan time.Duration = make(chan time.Duration, 10)

		if trafficType == "http" {
			go httpWorker(host, sleepTimeInMilliseconds[i], workersChannel)
			go responseTimeWorker(i+1, outputDir, workersChannel, &wg)
		} else if trafficType == "grpc" {
			go grpcWorker(host, sleepTimeInMilliseconds[i], workersChannel)
			go responseTimeWorker(i+1, outputDir, workersChannel, &wg)
		} else {
			log.Fatalf("Invalid traffic type: %s", trafficType)
		}
	}
	wg.Wait()
	fmt.Println("All workers finished")
}
