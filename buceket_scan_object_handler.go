package main

import (
	// standard
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

func BucketScanObjectHandler(w http.ResponseWriter, r *http.Request) {
	switch cloud_provider {
		case "AWS":
			retrun S3ScanFileHandler(w,r)
		case "GCP":
			retrun GCSScanFileHandler(w,r)
		case "AZURE":
			retrun ABSScanFileHandler(w,r)
		default:
			panic(fmt.Errorf("unwknown cloud_provider: %s", s))
    }
}
