package main

// Defining an interface
type ScannerInterface interface {
	Scan(data [] byte) (*ScanReport, error)
}

// struct to handle matches
type ScanMatch struct {
	Rule      string `json:"rule"`
	Namespace string `json:"namespace"`
	Tags      []string `json:"tags"`
}

type ListResponse struct {
	Files []string `json:"files"`
}

type ScanReport struct {
	Filename string  `json:"filename"`
	Matches  []ScanMatch `json:"matches"`
	Status string `json:"status"`
}

type ScanErrorData struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

type ScanErrorResponse struct {
	Error ScanErrorData `json:"error"`
}

func ScanStream(scannerIF ScannerInterface, data []byte) (*ScanReport, error) {
	return scannerIF.Scan(data)
}
