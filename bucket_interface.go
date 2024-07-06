package main

// Defining an interface
type BucketInterface interface {
	ScanObject(w http.ResponseWriter, r *http.Request)
}

func validateInputBucket(w http.ResponseWriter, bucket string) error {
	if (bucket == "") {
		errorResponse(w, "Invalid input bucket", http.StatusUnprocessableEntity)
		return errors.New("Invalid input bucket")
	}

	bucketExists, err := bucketExists(bucket)

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

func validateInputKey(w http.ResponseWriter, bucket string, key string) error {
	if (key == "") {
		errorResponse(w, "Invalid input key", http.StatusUnprocessableEntity)
		return errors.New("Invalid input key")
	}

	keyExists, err := keyExists(bucket,key)
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

func getQurantineFilesBucket(qurantineFilesBucket string) string{
	// input has more priority
	if (qurantineFilesBucket != "" ) {
		return qurantineFilesBucket
	} 
	if (quarantine_files_bucket != "" ) {
		return quarantine_files_bucket
	}
	return ""
}

func getCleanFilesBucket(cleanFilesBucket string) string{
	// input has more priority
	if (cleanFilesBucket != "" ) {
		return cleanFilesBucket
	} 
	if (clean_files_bucket != "" ) {
		return clean_files_bucket
	}
	return ""
}

func validateQrantineFilesBucket(w http.ResponseWriter, qurantineFilesBucket string) error {
	var bucket = getQurantineFilesBucket(qurantineFilesBucket)
	
	if (bucket == "" ) {
		errorResponse(w, "Invalid qurantine files bucket", http.StatusBadRequest)
		return errors.New("Invalid qurantine files bucket")

	} else {
		err := validateInputBucket(w,bucket)
		if err != nil {
			return err
		}
	}
	return nil
}

func validateCleanFilesBucket(w http.ResponseWriter, cleanFilesBucket string) error {

	var bucket = getCleanFilesBucket(cleanFilesBucket)

	if (bucket == "" ) {
		errorResponse(w, "Invalid clean files bucket", http.StatusBadRequest)
		return errors.New("Invalid clean files bucket")

	} else {
		err := validateInputBucket(w,bucket)
		if err != nil {
			return err
		}
	}
	return nil

}

func validateInputData(w http.ResponseWriter, data *ScanObject) error {

	err := validateInputBucket(w,data.BucketName)
	if err != nil {
		return err
	}

	err = validateInputKey(w,data.BucketName,data.Key)
	if err != nil {
		return err
	}
	
	err = validateQrantineFilesBucket(w,data.QurantineFilesBucket)
	if err != nil {
		return err
	}
	
	err = validateCleanFilesBucket(w,data.CleanFilesBucket)
	if err != nil {
		return err
	}

	return nil
}

func ScanObject(bucketIF BucketInterface, w http.ResponseWriter, r *http.Request) () {

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

	err = validateInputData(w,data)
	if err != nil {
		elog.Println(" validateInputData failed " + err.Error())
		return
	}

	resp, _ := json.Marshal(data)
	info.Println(" Received ScanS3 request " + string(resp))
		
	byteData, err := bucketIF.readFile(data.BucketName, data.Key)
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
			err = bucketIF.copyFile(data.BucketName, data.Key, getQurantineFilesBucket(data.QurantineFilesBucket))
			if err != nil {
				elog.Println(err)
				errorResponse(w, err.Error(), http.StatusInternalServerError)
				return
			}
			err = bucketIF.deleteFile(data.BucketName, data.Key)
			if err != nil {
				elog.Println(err)
				errorResponse(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else if response.data.Status == "CLEAN" {
			info.Println("Key " +data.Key+ " from bucket "+data.BucketName+ " is Clean")
			err = bucketIF.copyFile(data.BucketName, data.Key, getCleanFilesBucket(data.CleanFilesBucket))
			if err != nil {
				elog.Println(err)
				errorResponse(w, err.Error(), http.StatusInternalServerError)
				return 
			}
			err = bucketIF.deleteFile(data.BucketName, data.Key)
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
