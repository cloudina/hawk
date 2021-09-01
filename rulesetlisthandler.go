package main

import (
	// standard
	"encoding/json"
	"net/http"
)

func RuleSetListHandler(w http.ResponseWriter, r *http.Request) {
	req := NewRuleSetRequest()
	namerequests <- req

	response := <-req.ResponseChan
	var err error = response.err

	if err != nil {
		elog.Println(err)
		errorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	output, err := json.Marshal(response.data)
	if err != nil {
		elog.Println(err)
		errorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendJsonResponse(w, output)
}
