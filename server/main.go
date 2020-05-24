package main

import (
	"flag"
	"log"
	"net/http"
)

// DB + cors Settings
var db_host, db_prefix, cors_allowed string

// Rate Limit + Proxy Settings
var rl_max_req, rl_max_ipr, rl_bigquery, rl_proxy int

// main - Start the server
func main() {
	// Get Flags
	flag.StringVar(&db_host, "db_host", "localhost:6379", "[Database] DB host")
	flag.StringVar(&db_prefix, "db_prefix", "rkt:", "[Database] Entries prefix")
	flag.StringVar(&cors_allowed, "cors_allowed", "http://localhost", "[CORS] Allowed domain(s)")
	flag.IntVar(&rl_max_req, "rl_max_req", 2048, "[Rate Limit] Max data added in Kb per request")
	flag.IntVar(&rl_max_ipr, "rl_max_ipr", 10240, "[Rate Limit] Max data added in Kb per IP block every 10 minutes")
	flag.IntVar(&rl_bigquery, "rl_bigquery", 512, "[Rate Limit] Data added in Kb for query to be considered big (limit to 10 min, penalty)")
	flag.IntVar(&rl_proxy, "rl_proxy", -1, "[Rate Limit] -1 No Proxy, 1 Trust Proxy, 2 Trust CloudFlare")
	flag.Parse()
	// Prep Flags
	rl_max_req = rl_max_req * 1024
	rl_max_ipr = rl_max_ipr * 1024
	rl_bigquery = rl_bigquery * 1024
	// Start Server
	log.Println("Starting Server...")
	pool = newPool(db_host)
	router := NewRouter()
	log.Println("Listening...")
	log.Fatal(http.ListenAndServe(":8080", router))
}

// HandleError - Handles unexpected errors
func HandleError(err error) {
	if err != nil {
		panic(err)
	}
}
