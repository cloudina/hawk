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

var CloudProviderMap = map[string]CloudProvider{
	"AWS":   CloudProviderAWS,
	"AZURE": CloudProviderAzure,
	"GCP":   CloudProviderGCP,
}

func ParseCloudProviderString(str string) (CloudProvider, bool) {
	c, ok := CloudProviderMap[strings.ToUpper(str)]
	return c, ok
}

func BucketScanObjectHandler(w http.ResponseWriter, r *http.Request) {

	switch cloud_provider {
	case CloudProviderAWS:
		s3_Mgr := &S3_Manager{}
		ScanBucketObject(w, r, s3_Mgr)
	case CloudProviderAzure:
		abs_Mgr := &ABS_Manager{}
		ScanBucketObject(w, r, abs_Mgr)
	case CloudProviderGCP:
		gcs_Mgr := &GCS_Manager{}
		ScanBucketObject(w, r, gcs_Mgr)
	default:
		panic(fmt.Errorf("unwknown cloud_provider: %s", cloud_provider))
	}

}
