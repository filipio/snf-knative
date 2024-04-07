package main

import (
	"fmt"
	"net/http"

	"github.com/kitabisa/teler-waf"
	"github.com/kitabisa/teler-waf/request"
	"github.com/kitabisa/teler-waf/threat"
)

var myHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	fmt.Println("received request")
	// This is the handler function for the route that we want to protect
	// with teler-waf's security measures.
	// Print headers of request
	for k, v := range r.Header {
		fmt.Fprintf(w, "%s: %v\n", k, v)
	}
	fmt.Println("Request URI: ", r.RequestURI)
	w.Write([]byte("hello darkness"))
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
		// Exclude specific threats from being checked by the teler-waf.
		Excludes: []threat.Threat{
			threat.BadReferrer,
			threat.BadCrawler,
		},
		// Specify whitelisted URIs (path & query parameters), headers,
		// or IP addresses that will always be allowed by the teler-waf
		// with DSL expressions.
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

	// Set the rejectHandler as the handler for the telerMiddleware.
	telerMiddleware.SetHandler(rejectHandler)

	// Create a new handler using the handler method of the Teler instance
	// and pass in the myHandler function for the route we want to protect.
	app := telerMiddleware.Handler(myHandler)
	// todo: check if this is actually working (meaning image is not stale)
	fmt.Println("server started at :8081")
	// Use the app handler as the handler for the route.
	http.ListenAndServe(":8081", app)
}
