package main

import (
	// standard
	"net/http"
	"errors"
	"encoding/json"
	"log"
)

func _getQurantineFilesBucket(qurantineFilesBucket string) string{
	// input has more priority
	if (qurantineFilesBucket != "" ) {
		return qurantineFilesBucket
	} 
	if (quarantine_files_bucket != "" ) {
		return quarantine_files_bucket
	}
	return ""
}

func _getCleanFilesBucket(cleanFilesBucket string) string{
	// input has more priority
	if (cleanFilesBucket != "" ) {
		return cleanFilesBucket
	} 
	if (clean_files_bucket != "" ) {
		return clean_files_bucket
	}
	return ""
}

func validateInputBucket(w http.ResponseWriter, bucket string, bucketInterface BucketInterface) error {
	if (bucket == "") {
		errorResponse(w, "Invalid input bucket", http.StatusUnprocessableEntity)
		return errors.New("Invalid input bucket")
	}

	bucketExists, err := bucketInterface.bucketExists(bucket)

	if(err != nil) {
		errorResponse(w,  err.Error(), http.StatusInternalServerError)
		return err
	}
	if (!bucketExists) {
		errorResponse(w,  "Bucket: "+bucket+" does not exists", http.StatusUnprocessableEntity)
		return errors.New("Bucket: "+bucket+" does not exists")
	}
	return nil
}

func validateInputKey(w http.ResponseWriter, bucket string, key string, bucketInterface BucketInterface) error {
	if (key == "") {
		errorResponse(w, "Invalid input key", http.StatusUnprocessableEntity)
		return errors.New("Invalid input key")
	}

	keyExists, err := bucketInterface.keyExists(bucket,key)
	if(err != nil) {
		errorResponse(w,  err.Error(), http.StatusInternalServerError)
		return err
	}
	if (!keyExists) {
		errorResponse(w,  "Key: "+key+" does not exist in Bucket: "+bucket, http.StatusUnprocessableEntity)
		return errors.New("Key: "+key+" does not exist in Bucket: "+bucket)
	}
	return nil
}

func validateQrantineFilesBucket(w http.ResponseWriter, qurantineFilesBucket string, bucketInterface BucketInterface) error {
	var bucket = _getQurantineFilesBucket(qurantineFilesBucket)
	
	if (bucket == "" ) {
		errorResponse(w, "Invalid qurantine files bucket", http.StatusBadRequest)
		return errors.New("Invalid qurantine files bucket")

	} else {
		err := validateInputBucket(w,bucket, bucketInterface)
		if err != nil {
			return err
		}
	}
	return nil
}

func validateCleanFilesBucket(w http.ResponseWriter, cleanFilesBucket string, bucketInterface BucketInterface) error {

	var bucket = _getCleanFilesBucket(cleanFilesBucket)

	if (bucket == "" ) {
		errorResponse(w, "Invalid clean files bucket", http.StatusBadRequest)
		return errors.New("Invalid clean files bucket")

	} else {
		err := validateInputBucket(w,bucket, bucketInterface)
		if err != nil {
			return err
		}
	}
	return nil

}

func validateInputData(w http.ResponseWriter, data *ScanObject, bucketInterface BucketInterface) error {

	err := validateInputBucket(w,data.BucketName, bucketInterface)
	if err != nil {
		return err
	}

	err = validateInputKey(w,data.BucketName,data.Key, bucketInterface)
	if err != nil {
		return err
	}
	
	err = validateQrantineFilesBucket(w,data.QurantineFilesBucket, bucketInterface)
	if err != nil {
		return err
	}
	
	err = validateCleanFilesBucket(w,data.CleanFilesBucket, bucketInterface)
	if err != nil {
		return err
	}

	return nil
}

func ScanBucketObject(w http.ResponseWriter, r *http.Request, bucketInterface BucketInterface) () {

	data := new(ScanObject)
	err := decodeJSONBody(w, r, &data)
    if err != nil {
        var mr *malformedRequest
        if errors.As(err, &mr) {
            errorResponse(w, mr.msg, mr.status)
        } else {
            log.Println(err.Error())
            errorResponse(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
        }
        return
    }

	err = validateInputData(w,data, bucketInterface)
	if err != nil {
		elog.Println(" validateInputData failed " + err.Error())
		return
	}

	resp, _ := json.Marshal(data)
	info.Println(" Received ScanS3 request " + string(resp))
		
	byteData, err := bucketInterface.readFile(data.BucketName, data.Key)
	if err != nil {
		elog.Println(err)
		errorResponse(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// send request for scanning
	newRequest := NewScanStreamRequest(byteData)
	scanstreamrequests <- newRequest

	response := <-newRequest.ResponseChan

	err = response.err

	if err != nil {
		elog.Println(err)
		errorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	} else  {
		if response.data.Status == "INFECTED" {
			elog.Println("Key " +data.Key+ " from bucket "+data.BucketName+ " is Infected")
			err = bucketInterface.copyFile(data.BucketName, data.Key, _getQurantineFilesBucket(data.QurantineFilesBucket))
			if err != nil {
				elog.Println(err)
				errorResponse(w, err.Error(), http.StatusInternalServerError)
				return
			}
			err = bucketInterface.deleteFile(data.BucketName, data.Key)
			if err != nil {
				elog.Println(err)
				errorResponse(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else if response.data.Status == "CLEAN" {
			info.Println("Key " +data.Key+ " from bucket "+data.BucketName+ " is Clean")
			err = bucketInterface.copyFile(data.BucketName, data.Key, _getCleanFilesBucket(data.CleanFilesBucket))
			if err != nil {
				elog.Println(err)
				errorResponse(w, err.Error(), http.StatusInternalServerError)
				return 
			}
			err = bucketInterface.deleteFile(data.BucketName, data.Key)
			if err != nil {
				elog.Println(err)
				errorResponse(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	output, err := json.Marshal(response.data)
	if err != nil {
		elog.Println(err)
		errorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendJsonResponse(w, output)
	//fmt.Fprintf(w, string(output))
}
