package main

type ABS_Manager struct {
	*BucketMgr
}

func (self *ABS_Manager ) bucketExists(bucket string) (bool, error)
func (self *ABS_Manager ) keyExists(bucket string, key string) (bool, error)
func (self *ABS_Manager ) readFile(bucket string, item string) ([] byte, error) 
func (self *ABS_Manager ) copyFile(bucket string, item string, other string) (error)
func (self *ABS_Manager ) deleteFile(bucket string, item string) (error) 
