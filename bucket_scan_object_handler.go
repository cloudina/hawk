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

func initialiseBucketInterface(cloud_provider string, interface* BucketInterface ) (){
	switch cloud_provider {
		case CloudProviderAWS:
			interface := S3_Manager()
		case CloudProviderAzure:
			interface = ABS_Manager()
		case CloudProviderGCP:
			interface = GCP_Manager()
		default:
			panic(fmt.Errorf("unwknown cloud_provider: %s", cloud_provider))
    }
}

func BucketScanObjectHandler(w http.ResponseWriter, r *http.Request) {
	ScanObject(bucketInterface, w, r)
}

