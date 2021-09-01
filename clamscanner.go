package main

import (
	"github.com/dutchcoders/go-clamd"
	"bytes"
	"encoding/json"
	"time"
	"fmt"
)

var eicar = []byte(`X5O!P%@AP[4\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*`)

type ClamScanner struct {
	clamdaddr string
}

func NewClamScanner(clamdaddr string) (*ClamScanner, error) {
	scanner := new(ClamScanner)
	scanner.clamdaddr = clamdaddr
	return scanner, nil
}

func (self *ClamScanner) Scan(data [] byte) (*ScanReport,error) {
	var matches []ScanMatch
	response := new(ScanReport)
	response.Filename = "stream"

	clamdScanner := clamd.NewClamd(self.clamdaddr)
	abortChannel := make(chan bool)
	defer close(abortChannel) // necessary to not leak a goroutine. See https://github.com/dutchcoders/go-clamd/issues/9	ioreader := bytes.NewReader(data)
	ioreader := bytes.NewReader(data)
	ch, err := clamdScanner.ScanStream(ioreader, abortChannel)

	if err != nil {
		response.Status = "ERROR"
		return response, err
	}

	r := (<-ch)	//defer close(response)

	respJson, err := json.Marshal(&r)
	if err != nil {
		response.Status = "ERROR"
		return response, err
	}
	fmt.Printf(time.Now().Format(time.RFC3339)+" Scan result :  %v\n", string(respJson))

	switch r.Status {
		case clamd.RES_OK:
			response.Status = "CLEAN"
		case clamd.RES_FOUND:
			response.Status = "INFECTED"
			var match ScanMatch
			match.Namespace = ""
			match.Tags = nil
			match.Rule = r.Description
			matches = append(matches, match)
		case clamd.RES_ERROR:
		case clamd.RES_PARSE_ERROR:
		default:
			response.Status = "ERROR"
	}

	if len(matches) <= 0 {
		matches = [] ScanMatch{}
	}
	
	response.Matches = matches
	fmt.Printf(time.Now().Format(time.RFC3339) + " Finished scanning: " + "\n")

	go func() {
		for range ch {
		} // empty the channel so the goroutine from go-clamd/*CLAMDConn.readResponse() doesn't get stuck
	}()

	return response,nil
	
}

func (self *ClamScanner) ping() error {
	clamdScanner := clamd.NewClamd(self.clamdaddr)
	return clamdScanner.Ping()
}

func (self *ClamScanner) version() (string, error) {
	clamdScanner := clamd.NewClamd(self.clamdaddr)
	ch, err := clamdScanner.Version()
	if err != nil {
		return "", err
	}

	r := (<-ch)
	return r.Raw, nil
}

func (self *ClamScanner) isClamdReady() bool {
	if err := self.ping(); err != nil {
		fmt.Printf("ClamD ping failed.. error [%v]\n", err)
		return false
	} 

	fmt.Printf("Connectted to ClamD Server\n")
	if response, err := self.version(); err != nil {
			fmt.Printf("ClamD version check failed.. error [%v]\n", err)
			return false
	} else {
			fmt.Printf("ClamD version: %#v\n", response)
	}
	
	return true
		
}

func (self *ClamScanner) runScanCheck() bool {
	if ! self.isClamdReady() {
		return false
	}

	if _, err := self.Scan(eicar); err != nil {
		fmt.Printf("ClamD EICAR scan check failed.. error [%v]\n", err)
		return false
	}

	if _, err := self.Scan([]byte("hello world... how are you")); err != nil {
		fmt.Printf("ClamD sample text scan check failed.. error [%v]\n", err)
		return false
	} 
	return true

}

func (self *ClamScanner) warmUp() bool {
	for i:=0; i < 24 ; i++ {
		if ! self.runScanCheck() { 
			time.Sleep(time.Second * 5)
		} else {
			return true
		}
	}
	return false
}