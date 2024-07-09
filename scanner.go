package main

import (
	"encoding/json"
	"os"
	"strconv"
	"time"
)

// struct to hold compiler and channels
type Scanner struct {
	yarascanner         YaraScanner
	clamscanner         ClamScanner
	healthcheckrequests chan *HealthCheckRequest
	scanstreamrequests  chan *ScanStreamRequest
	namerequests        chan *RuleSetRequest
	rulerequests        chan *RuleListRequest
}

func (self *Scanner) healthcheck() *HealthCheckResponse {

	healthCheckResponse := new(HealthCheckResponse)
	clamDHealth := self.clamscanner.isClamdReady()

	if clamDHealth {
		healthCheckResponse.health = "OK"
	} else {
		healthCheckResponse.health = "ERROR"
	}

	return healthCheckResponse
}

func (self *Scanner) scanstream(data []byte) *ScanResponse {

	info.Println("Running yarascan")

	scanResponse := new(ScanResponse)

	yaraScannerResponse, yaraerr := ScanStream(&self.yarascanner, data)
	scanResponse.data = yaraScannerResponse
	scanResponse.err = yaraerr

	yaraRespJson, _ := json.Marshal(yaraScannerResponse)
	info.Println(time.Now().Format(time.RFC3339) + " yarascan scan result " + string(yaraRespJson))

	if (yaraerr == nil) && len(yaraScannerResponse.Matches) > 0 {
		info.Println(time.Now().Format(time.RFC3339) + " Found matches with yara " + string(yaraRespJson))
	}

	info.Println("Running clamscan on addr: " + clamdaddr)

	clamScannerResponse, clamerr := ScanStream(&self.clamscanner, data)

	clamRespJson, _ := json.Marshal(clamScannerResponse)
	info.Println(time.Now().Format(time.RFC3339) + " clamav scan result " + string(clamRespJson))

	if (clamerr == nil) && len(clamScannerResponse.Matches) > 0 {
		info.Println(time.Now().Format(time.RFC3339) + " Found matches with clamav" + string(clamRespJson))
		scanResponse.data = clamScannerResponse
	}

	if clamerr != nil {
		scanResponse.err = clamerr
	}

	return scanResponse
}

func (self *Scanner) warmUp() {

	info.Println("Warming Up")

	var yaraHealth = bool(false)
	var clamDHealth = bool(false)

	yaraScannerResponse, yaraerr := ScanStream(&self.yarascanner, eicar)

	if (yaraerr == nil) && len(yaraScannerResponse.Matches) > 0 {
		yaraHealth = true
	}

	clamDHealth = self.clamscanner.warmUp()

	if yaraHealth && clamDHealth {
		info.Println("Warmed Up")
	} else {
		info.Println(time.Now().Format(time.RFC3339) + " Warm up failed exiting.. Yara Health" + strconv.FormatBool(yaraHealth) + "ClamD Health" + strconv.FormatBool(clamDHealth))
		os.Exit(1)
	}
}

func (self *Scanner) LoadIndex(indexPath string) error {
	return self.yarascanner.LoadIndex(indexPath)
}

func (self *Scanner) listRuleSets() *RuleSetResponse {
	response, err := self.yarascanner.ListRuleSets()
	ruleSetResponse := new(RuleSetResponse)
	ruleSetResponse.err = err
	ruleSetResponse.data = response
	return ruleSetResponse
}

func (self *Scanner) listRules(rulesetname string) *RuleListResponse {

	response, err := self.yarascanner.ListRules(rulesetname)
	ruleListResponse := new(RuleListResponse)
	ruleListResponse.err = err
	ruleListResponse.data = response

	return ruleListResponse
}

func (self *Scanner) Run() {
	info.Println("Waiting for scan requests")
	for {
		select {
		case healthcheckmsg := <-healthcheckrequests:
			response := self.healthcheck()
			healthcheckmsg.ResponseChan <- response
		case scanstreammsg := <-scanstreamrequests:
			response := self.scanstream(scanstreammsg.data)
			scanstreammsg.ResponseChan <- response
		case setmsg := <-namerequests:
			response := self.listRuleSets()
			setmsg.ResponseChan <- response
		case rulemsg := <-rulerequests:
			response := self.listRules(rulemsg.RuleSet)
			rulemsg.ResponseChan <- response
		}
	}
}

func NewScanner(healthcheckreq chan *HealthCheckRequest, scanstream chan *ScanStreamRequest, name chan *RuleSetRequest, list chan *RuleListRequest) (*Scanner, error) {
	scanner := new(Scanner)
	scanner.healthcheckrequests = healthcheckreq
	scanner.scanstreamrequests = scanstream
	scanner.namerequests = name
	scanner.rulerequests = list
	scanner.yarascanner = YaraScanner{}
	scanner.clamscanner = ClamScanner{clamdaddr}
	return scanner, nil
}
