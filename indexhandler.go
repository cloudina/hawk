package main

import (
	// standard
	"encoding/json"
	"net/http"
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
