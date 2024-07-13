package main

import (
	// standard
	"encoding/json"
	"net/http"
)

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	validateContentType(w, r)

	// send request for scanning
	newRequest := NewHealthCheckRequest()
	healthcheckrequests <- newRequest

	response := <-newRequest.ResponseChan

	err := response.err

	if err != nil {
		elog.Println(err)
		errorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	output, err := json.Marshal(response.health)
	if err != nil {
		elog.Println(err)
		errorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendJsonResponse(w, output)
}
