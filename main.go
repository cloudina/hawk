package main

import (
	// standard
	"flag"
	"log"
	"net/http"
	"os"
	// external
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// config options
	index_files StringArgs
	address     string
	port        string
	addrport    string
	clamdaddr   string
	clean_files_bucket string
	quarantine_files_bucket string
	cloud_provider string
	// channels
	healthcheckrequests chan *HealthCheckRequest
	scanstreamrequests chan *ScanStreamRequest
	namerequests chan *RuleSetRequest
	rulerequests chan *RuleListRequest

	// loggers
	info *log.Logger
	elog *log.Logger
)

func init() {

	flag.Var(&index_files, "i", "path to yara rules")
	flag.StringVar(&address, "address", "0.0.0.0", "address to bind to")
	flag.StringVar(&port, "port", "9999", "port to bind to")
	flag.StringVar(&clamdaddr, "clamaddr", "tcp://localhost:3310", "clamd address to bind to")

	flag.Parse()

	// initialize logger
	info = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	elog = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	//build address string
	addrport = address + ":" + port

	clean_files_bucket = getEnv("CLEAN_FILES_BUCKET", "")
	quarantine_files_bucket = getEnv("QUARANTINE_FILES_BUCKET", "")
	cloud_provider = getEnv("CLOUD_PROVIDER", "")

	info.Println("reading CLEAN_FILES_BUCKET value as " +clean_files_bucket)
	info.Println("reading QUARANTINE_FILES_BUCKET value as " +quarantine_files_bucket)
	info.Println("reading CLOUD_PROVIDER value as " +cloud_provider)

}

func main() {
	// create channels
	info.Println("Initializing channels")
	healthcheckrequests = make(chan *HealthCheckRequest)
	scanstreamrequests = make(chan *ScanStreamRequest)
	namerequests = make(chan *RuleSetRequest)
	rulerequests = make(chan *RuleListRequest)
	// create scanner
	info.Println("Initializing scanner")
	scanner, err := NewScanner(healthcheckrequests, scanstreamrequests, namerequests, rulerequests)
	if err != nil {
		panic(err)
	}

	// load indexes
	for _, index := range index_files {
		info.Println("Loading index: " + index)
		err = scanner.LoadIndex(index)
		if err != nil {
			panic(err)
		}
	}

	// warmup the scanner
	scanner.warmUp()

	// launch scanner
	go scanner.Run()

	// setup http server and begin serving traffic
	r := mux.NewRouter()
	// helmet := CustomHelmet()
	// r.Use(helmet.Secure)

	helmet := SimpleHelmet{}
	helmet.Default()
	r.Use(helmet.Secure)
	r.NotFoundHandler = Handle404(helmet)

	r.HandleFunc("/", IndexHandler).Methods("GET")
	r.HandleFunc("/health", HealthCheckHandler).Methods("GET")
	r.HandleFunc("/scanstream", ScanStreamHandler).Methods("POST")
	// Prometheus metrics
	r.Handle("/metrics", promhttp.Handler())

	bucket_sub := r.PathPrefix("/bucket").Subrouter()
	bucket_sub.HandleFunc("/scanobject", BucketScanObjectHandler).Methods("POST")
	ruleset_sub := r.PathPrefix("/ruleset").Subrouter()
	ruleset_sub.HandleFunc("", RuleSetListHandler).Methods("GET")
	ruleset_sub.HandleFunc("/", RuleSetListHandler).Methods("GET")
	ruleset_sub.HandleFunc("/{ruleset}", RuleListHandler).Methods("GET")
	
	loggedRouter := handlers.CombinedLoggingHandler(os.Stdout, r)
	log.Fatal(http.ListenAndServe(addrport, loggedRouter))
	//log.Fatal(http.ListenAndServe(addrport, r))
}
