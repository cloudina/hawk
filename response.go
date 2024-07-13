package main

type ScanResponse struct {
	data *ScanReport
	err  error
}

type HealthCheckResponse struct {
	health string //OK,ERROR
	err    error
}

// struct to handle namespace requests
type RuleSetResponseObject struct {
	Names []string
}

type RuleSetResponse struct {
	data *RuleSetResponseObject
	err  error
}

// sturc to handle
type RuleListResponseObject struct {
	Rules []string
}

type RuleListResponse struct {
	data *RuleListResponseObject
	err  error
}
