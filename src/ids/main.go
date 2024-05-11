package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/kitabisa/teler-waf"
	"github.com/kitabisa/teler-waf/request"
	"github.com/kitabisa/teler-waf/threat"
)

func handleError(w http.ResponseWriter, successor string, err string) {
	errMsg := "Error sending POST request to " + successor + " : " + err
	fmt.Println(errMsg)
	http.Error(w, "Error sending POST request to "+successor+" : "+err, http.StatusInternalServerError)
}

var myHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

})

var rejectHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// This is the handler function for the route that we want to be rejected
	// if the teler-waf's security measures are triggered.
	http.Error(w, "Sorry, your request has been denied for security reasons.", http.StatusForbidden)
})

func main() {
	// Create a new instance of the Teler type using the New function
	// and configure it using the Options struct.
	telerMiddleware := teler.New(teler.Options{
		Excludes: []threat.Threat{
			threat.BadReferrer,
			threat.BadCrawler,
		},
		Whitelists: []string{
			`request.Headers matches "(curl|Go-http-client|okhttp)/*" && threat == BadCrawler`,
			`request.URI startsWith "/wp-login.php"`,
			`request.IP in ["127.0.0.1", "::1", "0.0.0.0"]`,
			`request.Headers contains "authorization" && request.Method == "POST"`,
		},
		// Specify custom rules for the teler-waf to follow.
		Customs: []teler.Rule{
			{
				// Give the rule a name for easy identification.
				Name: "Log4j Attack",
				// Specify the logical operator to use when evaluating the rule's conditions.
				Condition: "or",
				// Specify the conditions that must be met for the rule to trigger.
				Rules: []teler.Condition{
					{
						// Specify the HTTP method that the rule applies to.
						Method: request.GET,
						// Specify the element of the request that the rule applies to
						// (e.g. URI, headers, body).
						Element: request.URI,
						// Specify the pattern to match against the element of the request.
						Pattern: `\$\{.*:\/\/.*\/?\w+?\}`,
					},
				},
			},
			{
				// Give the rule a name for easy identification.
				Name: `Headers Contains "curl" String`,
				// Specify the conditions that must be met for the rule to trigger.
				Rules: []teler.Condition{
					{
						// Specify the DSL expression that the rule applies to.
						DSL: `request.Headers contains "curl"`,
					},
				},
			},
		},
		// Specify the file path to use for logging.
		LogFile: "/tmp/teler.log",
	})

	telerMiddleware.SetHandler(rejectHandler)
	app := telerMiddleware.Handler(myHandler)

	fmt.Println("server started at :8080")
	http.ListenAndServe(":8080", app)
}
