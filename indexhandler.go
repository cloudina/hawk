package main

import (
	// standard
	"net/http"
	"encoding/json"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	output, err := json.Marshal("hello")
	if err != nil {
		elog.Println(err)
		errorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendJsonResponse(w, output)
}