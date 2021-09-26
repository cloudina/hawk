package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"net/url"
	"time"
	"fmt"
	"os"
	"errors"
)

// check if a bucket exists.
func bucketExists(bucket string) (bool, error) {
	awsSession, _ := session.NewSession(&aws.Config{
		Region: aws.String(getRegion())},
	)

	svc := s3.New(awsSession)
	input := &s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	}

	_, err := svc.HeadBucket(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
				case s3.ErrCodeNoSuchBucket:
					return false,nil
				default:
					elog.Println( time.Now().Format(time.RFC3339) + " bucketExists failed for bucket "+bucket + " error : " + err.Error())
					return false, errors.New("Filed to find bucket")			
			}
		}
		elog.Println( time.Now().Format(time.RFC3339) + " bucketExists got unknown error for bucket "+bucket + " error : " + err.Error())
		return false, errors.New("Filed to find bucket")			
	}

	return true,nil
}

// check if a file exists.
func keyExists(bucket string, key string) (bool, error) {
	awsSession, _ := session.NewSession(&aws.Config{
		Region: aws.String(getRegion())},
	)

	svc := s3.New(awsSession)

	_, err := svc.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {            
			case "NotFound": // s3.ErrCodeNoSuchKey does not work, aws is missing this error code so we hardwire a string
				elog.Println( time.Now().Format(time.RFC3339) + " keyExists got NotFound error for " +key+ " bucket "+bucket + " error : " + err.Error())
				return false, nil
			default:
				elog.Println( time.Now().Format(time.RFC3339) + " keyExists failed for " +key+ " bucket "+bucket + " error : " + err.Error())
				return false, errors.New("Filed to find file")
			}
		}
		elog.Println( time.Now().Format(time.RFC3339) + " keyExists got unknown error for " +key+ " bucket "+bucket + " error : " + err.Error())
		return false, errors.New("Filed to find file")
	}
	return true, nil
}


func getRegion() string {
	region, err := os.LookupEnv("AWS_REGION")
	if !err {
		fmt.Println("AWS_REGION is not present..using us-east-1")
		region = "us-east-1"
	}
	return region
}

func readFile(bucket string, item string) ([] byte, error) {

	awsSession, _ := session.NewSession(&aws.Config{
		Region: aws.String(getRegion())},
	)
	
	buff := &aws.WriteAtBuffer{}

	s3dl := s3manager.NewDownloader(awsSession)

	_, err := s3dl.Download(buff, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(item),
	})

	if err != nil {
		elog.Println( time.Now().Format(time.RFC3339) + " Unable to read file " +item+ " from bucket "+bucket + " error : " + err.Error())
		return nil, errors.New("Unable to read file")
	}
	
	info.Println(time.Now().Format(time.RFC3339) +" Downloaded file "+item+ " from bucket "+bucket)

	return buff.Bytes(), nil
}

func copyFile(bucket string, item string, other string) (error){

	awsSession, _ := session.NewSession(&aws.Config{
		Region: aws.String(getRegion())},
	)

	// Create S3 service client
	svc := s3.New(awsSession)

	source := bucket + "/" + item

	// Copy the file
	_, err := svc.CopyObject(&s3.CopyObjectInput{Bucket: aws.String(other),
	CopySource: aws.String(url.PathEscape(source)), Key: aws.String(item),  ACL: aws.String("bucket-owner-full-control")})

	if err != nil {
		elog.Println( time.Now().Format(time.RFC3339) + " Unable to read file " +item+ " from bucket "+bucket+ " to bucket "+other+" error : " + err.Error())
		return errors.New("Unable to copy file")
	}

	// Wait to see if the file got copied
	err = svc.WaitUntilObjectExists(&s3.HeadObjectInput{Bucket: aws.String(other), Key: aws.String(item)})
	if err != nil {
		elog.Println( time.Now().Format(time.RFC3339) + " Error occurred while waiting for file " +item+ " to be copied to bucket "+other+ " error: "+ fmt.Sprint(err))
		return errors.New("Error while  waiting for file to copy")
	}

	info.Println( time.Now().Format(time.RFC3339) + " File "+ item+ " successfully copied from bucket "+bucket+ " to bucket "+other)

	return nil
}

func deleteFile(bucket string, item string) (error) {
	awsSession, _ := session.NewSession(&aws.Config{
		Region: aws.String(getRegion())},
	)

	// Create S3 service client
	svc := s3.New(awsSession)

	params := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(item),
	}
	_, err := svc.DeleteObject(params)
	if err != nil {
		elog.Println( time.Now().Format(time.RFC3339) + " Error occurred while deleting file " +item+ " from bucket "+bucket+" err: "+ fmt.Sprint(err))
		return errors.New("Error occurred while deleting file")
	}
	return nil
}
