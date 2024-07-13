package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/bloberror"
	"path"
	"log"
	"os"
)

type ABS_Manager struct {
	BucketMgr
}

func (self *ABS_Manager) handleError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func (self *ABS_Manager) getServiceClient() *azblob.Client {
	// Create a new service client with token credential
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}

	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)

	credential, err := azidentity.NewDefaultAzureCredential(nil)
	self.handleError(err)

	client, err := azblob.NewClient(serviceURL, credential, nil)
	self.handleError(err)
	return client
}

func (self *ABS_Manager) getContainerClient(containerName string) *container.Client {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}

	containerURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	self.handleError(err)

	containerClient, err := container.NewClient(containerURL, cred, nil)
	self.handleError(err)
	return containerClient
}

func (self *ABS_Manager) getBlobClient(containerName string, blobName string) *blob.Client {
	// From the Azure portal, get your Storage account blob service URL endpoint.
	accountName, accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"), os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")

	blobURL := fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s", accountName, containerName, blobName)
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	self.handleError(err)
	blobClient, err := blob.NewClientWithSharedKeyCredential(blobURL, credential, nil)
	self.handleError(err)
	return blobClient
}

func (self *ABS_Manager) listBuckets() []string {

	client := self.getServiceClient()

	pager := client.NewListContainersPager(&azblob.ListContainersOptions{
		Include: azblob.ListContainersInclude{Metadata: true, Deleted: false},
	})

	var buckets []string

	for pager.More() {
		resp, err := pager.NextPage(context.TODO())
		self.handleError(err) // if err is not nil, break the loop.
		for _, _container := range resp.ContainerItems {
			buckets = append(buckets, *_container.Name)
		}
	}
	return buckets
}

func (self *ABS_Manager) bucketExists(bucket string) (bool, error) {
	client := self.getContainerClient(bucket)
	_, err := client.GetProperties(context.TODO(), nil)

	if bloberror.HasCode(err, bloberror.ContainerNotFound) {
		return false, err
	} else {
		return true, nil
	}
}

func (self *ABS_Manager) keyExists(bucket string, key string) (bool, error) {
	client := self.getBlobClient(bucket, key)
	_, err := client.GetProperties(context.TODO(), nil)
	
	if bloberror.HasCode(err, bloberror.BlobNotFound) {
		return false, err
	} else {
		return true, nil
	}
}

func (self *ABS_Manager) readFile(bucket string, item string) ([]byte, error) {

	client := self.getServiceClient()
	// Download the blob
	downloadResponse, err := client.DownloadStream(context.TODO(), bucket, item, nil)
	self.handleError(err)

	downloadedData := bytes.Buffer{}
	retryReader := downloadResponse.NewRetryReader(context.TODO(), &azblob.RetryReaderOptions{})
	_, err = downloadedData.ReadFrom(retryReader)
	self.handleError(err)

	err = retryReader.Close()
	self.handleError(err)

	return downloadedData.Bytes(), nil
}

func (self *ABS_Manager) copyFile(bucket string, item string, other string) error {

	data, _ := self.readFile(bucket, item)

	client := self.getServiceClient()

	_, err := client.UploadBuffer(context.TODO(), path.Dir(other), path.Base(other), data, &azblob.UploadBufferOptions{})
	self.handleError(err)
	return err
}

func (self *ABS_Manager) deleteFile(bucket string, item string) error {
	client := self.getServiceClient()
	// Delete the blob.
	_, err := client.DeleteBlob(context.TODO(), bucket, item, nil)
	self.handleError(err)
	return err
}
