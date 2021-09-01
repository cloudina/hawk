package main

import (
	// standard
	"encoding/json"
	"net/http"
	// external
	"github.com/gorilla/mux"
)

func RuleListHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ruleset := vars["ruleset"]

	req := NewRuleListRequest(ruleset)
	rulerequests<- req

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
