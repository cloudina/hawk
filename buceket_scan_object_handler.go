package main

import (
	// standard
	"fmt"
	"net/http"
	"strings"
)

type CloudProvider int
const (
    CloudProviderAWS = iota
    CloudProviderAzure
    CloudProviderGCP
)

var CloudProviderMap = map[string] CloudProvider{
    "AWS": CloudProviderAWS ,
    "AZURE": CloudProviderAzure ,
    "GCP": CloudProviderGCP ,
}

func ParseCloudProviderString(str string) (CloudProvider, bool) {
    c, ok := CloudProviderMap[strings.ToUpper(str)]
    return c, ok
}

func BucketScanObjectHandler(w http.ResponseWriter, r *http.Request) {
	switch cloud_provider {
		case CloudProviderAWS:
			S3ScanFileHandler(w,r)
		case CloudProviderAzure:
			ABSScanFileHandler(w,r)
		case CloudProviderGCP:
			GCSScanFileHandler(w,r)
		default:
			panic(fmt.Errorf("unwknown cloud_provider: %s", cloud_provider))
    }
}
