package main

// Defining an interface
type BucketInterface interface {
	bucketExists(bucket string) (bool, error)
	keyExists(bucket string, key string) (bool, error)
	readFile(bucket string, item string) ([] byte, error) 
	copyFile(bucket string, item string, other string) (error)
	deleteFile(bucket string, item string) (error) 
}
