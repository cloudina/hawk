package main

import (
	// standard
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func ScanStreamHandler(w http.ResponseWriter, r *http.Request) {

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		elog.Println(err)
		errorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// send request for scanning
	newRequest := NewScanStreamRequest(buf)
	scanstreamrequests <- newRequest

	response := <-newRequest.ResponseChan

	err = response.err

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
