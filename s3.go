package main

import (
	"context"
	"net/http"
	"github.com/aws/smithy-go"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"net/url"
	"time"
	"fmt"
	"os"
	"strconv"
	"errors"
)


type S3_Manager struct {
	*BucketMgr
}

func (self* S3_Manager) getPartSize() int64 {
	var partSize int64

	strSizeInMb, err := os.LookupEnv("DOWNLOAD_PART_SIZE")
	
	if !err {
		elog.Println(time.Now().Format(time.RFC3339) + "DOWNLOAD_PART_SIZE is not present..using DefaultDownloadPartSize ")
		partSize = manager.DefaultDownloadPartSize
	} else {
		sizeInMb, err := strconv.Atoi(strSizeInMb)
		if err != nil {
			elog.Println(time.Now().Format(time.RFC3339) + " DOWNLOAD_PART_SIZE conversion issue..using DefaultDownloadPartSize ")
			partSize = manager.DefaultDownloadPartSize
		} else {
			partSize = int64(sizeInMb) * 1024 * 1204
		}
	}
	return partSize
}

func (self* S3_Manager) getRegion() string {
	region, err := os.LookupEnv("AWS_REGION")
	if !err {
		elog.Println(time.Now().Format(time.RFC3339) + " AWS_REGION is not present..using us-east-1")
		region = "us-east-1"
	}
	return region
}

// check if a bucket exists.
func (self *S3_Manager ) bucketExists(bucket string) (bool, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(self.getRegion()),
	)
	if err != nil {
		elog.Println( time.Now().Format(time.RFC3339) + " bucketExists: Filed to load config for bucket "+bucket + " error : " + err.Error())
		return false, errors.New("Filed to load config")			
	}

	s3client := s3.NewFromConfig(cfg)

	_, err = s3client.HeadBucket(context.TODO(),&s3.HeadBucketInput{Bucket: aws.String(bucket)})

	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {			
			var httpResponseErr *awshttp.ResponseError
			if  errors.As(err, &httpResponseErr) {
				switch httpResponseErr.HTTPStatusCode() {
				case http.StatusMovedPermanently:
					elog.Println( time.Now().Format(time.RFC3339) + " bucketExists: failed for bucket "+bucket + " error : " + err.Error())
					return false, errors.New("Bucket StatusMovedPermanently ")	
				case http.StatusForbidden:
					elog.Println( time.Now().Format(time.RFC3339) + " bucketExists: failed for bucket "+bucket + " error : " + err.Error())
					return false, errors.New("Bucket StatusForbidden")	
				case http.StatusNotFound:
					elog.Println( time.Now().Format(time.RFC3339) + " bucketExists: failed for bucket "+bucket + " error : " + err.Error())
					return false, nil				
				default:
					elog.Println(time.Now().Format(time.RFC3339) + " bucketExists: ResponseError failed for bucket "+bucket + "with error: "+err.Error())
					return false, errors.New("Filed to find bucket")
				}
			} else {
				elog.Println(time.Now().Format(time.RFC3339) + " bucketExists: ApiError failed for bucket "+bucket + "with error: "+err.Error())
				return false, errors.New("Filed to find bucket")
			}
		} else {
			elog.Println(time.Now().Format(time.RFC3339) + " bucketExists: failed for bucket "+bucket + "with error: "+err.Error())
			return false, errors.New("Filed to find bucket")
		}
	}

	return true,nil
}

func (self* S3_Manager) getHeadObject(bucket string, key string) (*s3.HeadObjectOutput, error) {

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(self.getRegion()),
	)
	if err != nil {
		elog.Println( time.Now().Format(time.RFC3339) + " getHeadObject: Filed to load config for bucket "+bucket + " error : " + err.Error())
		return nil, errors.New("Filed to load config")			
	}

	s3client := s3.NewFromConfig(cfg)

	headObjectResponse, err := s3client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {			
			var httpResponseErr *awshttp.ResponseError
			if  errors.As(err, &httpResponseErr) {
				switch httpResponseErr.HTTPStatusCode() {
				case http.StatusMovedPermanently:
					elog.Println( time.Now().Format(time.RFC3339) + " getHeadObject: failed for bucket "+bucket +" key "+key+" error : " + err.Error())
					return nil, errors.New("Bucket StatusMovedPermanently ")	
				case http.StatusForbidden:
					elog.Println( time.Now().Format(time.RFC3339) + " getHeadObject: failed for bucket "+bucket +" key "+key+" error : " + err.Error())
					return nil, errors.New("Bucket StatusForbidden")	
				case http.StatusNotFound:
					elog.Println( time.Now().Format(time.RFC3339) + " getHeadObject: failed for bucket "+bucket +" key "+key+" error : " + err.Error())
					return nil, errors.New("Bucket StatusNotFound")					
				default:
					elog.Println(time.Now().Format(time.RFC3339) + " getHeadObject: ResponseError failed for bucket "+bucket +" key "+key+" with error: "+err.Error())
					return nil, errors.New("Filed to find object")
				}
			} else {
				elog.Println(time.Now().Format(time.RFC3339) + " getHeadObject: APIError failed for bucket "+bucket +" key "+key+" with error: "+err.Error())
				return nil, errors.New("Filed to find object")
			}
		} else {
			elog.Println(time.Now().Format(time.RFC3339) + " getHeadObject: failed for bucket "+bucket +" key "+key+" with error: "+err.Error())
			return nil, errors.New("Filed to find object")
		}
	}

	return headObjectResponse, nil
}

// check if a file exists.
func (self* S3_Manager) keyExists(bucket string, key string) (bool, error) {

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(self.getRegion()),
	)
	if err != nil {
		elog.Println( time.Now().Format(time.RFC3339) + "keyExists: Filed to load config for bucket "+bucket +" key "+key+" error : " + err.Error())
		return false, errors.New("Filed to load config")			
	}

	s3client := s3.NewFromConfig(cfg)

	_, err = s3client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {			
			var httpResponseErr *awshttp.ResponseError
			if  errors.As(err, &httpResponseErr) {
				switch httpResponseErr.HTTPStatusCode() {	
				case http.StatusNotFound:
					elog.Println( time.Now().Format(time.RFC3339) + " keyExists: failed for bucket "+bucket +" key "+key+" error : " + err.Error())
					return false, nil				
				default:
					elog.Println(time.Now().Format(time.RFC3339) + " keyExists: ResponseError failed for bucket "+bucket +" key "+key+" with error: "+err.Error())
					return false, errors.New("Filed to find key")
				}
			}  else {
				elog.Println(time.Now().Format(time.RFC3339) + " keyExists: APIErrorfailed for bucket "+bucket +" key "+key+" with error: "+err.Error())
				return false, errors.New("Filed to find key")
			}
		} else {
			elog.Println(time.Now().Format(time.RFC3339) + " keyExists: failed for bucket "+bucket +" key "+key+" with error: "+err.Error())
			return false, errors.New("Filed to find key")
		}
	}

	return true, nil
}

func (self *S3_Manager ) readFile(bucket string, item string) ([] byte, error) {

	// Load AWS Config
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(self.getRegion()),
	)
	if err != nil {
		elog.Println( time.Now().Format(time.RFC3339) + " readFile: Filed to load config to read file " +item+ " from bucket "+bucket + " error : " + err.Error())
		return nil, errors.New("Filed to load config")			
	}

	// Create an S3 client using the loaded configuration
	s3client := s3.NewFromConfig(cfg)

	// Create a downloader with the client and custom downloader options
	downloader := manager.NewDownloader(s3client, func(d *manager.Downloader) {
		d.PartSize = self.getPartSize()
	})

	headObject, err := self.getHeadObject(bucket,item)
	if err != nil {
		elog.Println( time.Now().Format(time.RFC3339) + " readFile: getHeadObject failed " +item+ " from bucket "+bucket + " error : " + err.Error())
		return nil, errors.New("Filed to read file")			
	}
	// pre-allocate in memory buffer, where headObject type is *s3.HeadObjectOutput
	buff := make([]byte, int(*headObject.ContentLength))
	// wrap with aws.WriteAtBuffer
	w := manager.NewWriteAtBuffer(buff)
	// download file into the memory
	_, err = downloader.Download(context.TODO(), w, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(item),
	})

	if err != nil {
		elog.Println( time.Now().Format(time.RFC3339) + " Unable to read file " +item+ " from bucket "+bucket + " error : " + err.Error())
		return nil, errors.New("Unable to read file")
	}
	
	info.Println(time.Now().Format(time.RFC3339) +" Downloaded file "+item+ " from bucket "+bucket)

	return buff, nil
}

func (self *S3_Manager ) copyFile(bucket string, item string, other string) (error){

	// Load AWS Config
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(self.getRegion()),
	)
	if err != nil {
		elog.Println( time.Now().Format(time.RFC3339) + " copyFile: Filed to load config to read file " +item+ " from bucket "+bucket + " error : " + err.Error())
		return errors.New("Filed to load config")			
	}

	// Create an S3 client using the loaded configuration
	s3client := s3.NewFromConfig(cfg)

	source := bucket + "/" + item

	_, err = s3client.CopyObject(context.TODO(), &s3.CopyObjectInput{
		Bucket:     aws.String(other),
		CopySource: aws.String(url.PathEscape(source)),
		Key:        aws.String(item),
	})

	if err != nil {
		elog.Println( time.Now().Format(time.RFC3339) + " Unable to read file " +item+ " from bucket "+bucket+ " to bucket "+other+" error : " + err.Error())
		return errors.New("Unable to copy file")
	}

	info.Println( time.Now().Format(time.RFC3339) + " File "+ item+ " successfully copied from bucket "+bucket+ " to bucket "+other)

	return nil
}

func (self *S3_Manager ) deleteFile(bucket string, item string) (error) {

	// Load AWS Config
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(self.getRegion()),
	)
	if err != nil {
		elog.Println( time.Now().Format(time.RFC3339) + " deleteFile: Filed to load config to read file " +item+ " from bucket "+bucket + " error : " + err.Error())
		return errors.New("Filed to load config")			
	}

	// Create an S3 client using the loaded configuration
	s3client := s3.NewFromConfig(cfg)

	_, err = s3client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(item),
	})

	if err != nil {
		elog.Println( time.Now().Format(time.RFC3339) + " Error occurred while deleting file " +item+ " from bucket "+bucket+" err: "+ fmt.Sprint(err))
		return errors.New("Error occurred while deleting file")
	}
	return nil
}
