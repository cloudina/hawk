package main

type ScanObject struct {
	BucketName           string `json:"bucketname"`
	Key                  string `json:"key"`
	CleanFilesBucket     string `json:"clean_files_bucket,omitempty"`
	QurantineFilesBucket string `json:"qurantine_files_bucket,omitempty"`
}

type HealthCheckRequest struct {
	ResponseChan chan *HealthCheckResponse
}

type ScanStreamRequest struct {
	data         []byte
	ResponseChan chan *ScanResponse
}

func NewHealthCheckRequest() *HealthCheckRequest {
	healthcheckreq := new(HealthCheckRequest)
	healthcheckreq.ResponseChan = make(chan *HealthCheckResponse)
	return healthcheckreq
}

func NewScanStreamRequest(data []byte) *ScanStreamRequest {
	scan := new(ScanStreamRequest)
	scan.data = data
	scan.ResponseChan = make(chan *ScanResponse)

	return scan
}

type RuleSetRequest struct {
	ResponseChan chan *RuleSetResponse
}

func NewRuleSetRequest() *RuleSetRequest {
	rule := new(RuleSetRequest)
	rule.ResponseChan = make(chan *RuleSetResponse)

	return rule
}

type RuleListRequest struct {
	RuleSet      string
	ResponseChan chan *RuleListResponse
}

func NewRuleListRequest(ruleset string) *RuleListRequest {
	rule := new(RuleListRequest)
	rule.RuleSet = ruleset
	rule.ResponseChan = make(chan *RuleListResponse)

	return rule
}
