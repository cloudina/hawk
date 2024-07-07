package main

type GCS_Manager struct {
	*BucketMgr
}
func (self *GCS_Manager ) bucketExists(bucket string) (bool, error)
func (self *GCS_Manager ) keyExists(bucket string, key string) (bool, error)
func (self *GCS_Manager ) readFile(bucket string, item string) ([] byte, error) 
func (self *GCS_Manager ) copyFile(bucket string, item string, other string) (error)
func (self *GCS_Manager ) deleteFile(bucket string, item string) (error) 



import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
)

// copyFile copies an object into specified bucket.
func copyFile(w io.Writer, dstBucket, srcBucket, srcObject string) error {
	// dstBucket := "bucket-1"
	// srcBucket := "bucket-2"
	// srcObject := "object"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
			return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	dstObject := srcObject + "-copy"
	src := client.Bucket(srcBucket).Object(srcObject)
	dst := client.Bucket(dstBucket).Object(dstObject)

	// Optional: set a generation-match precondition to avoid potential race
	// conditions and data corruptions. The request to copy is aborted if the
	// object's generation number does not match your precondition.
	// For a dst object that does not yet exist, set the DoesNotExist precondition.
	dst = dst.If(storage.Conditions{DoesNotExist: true})
	// If the destination object already exists in your bucket, set instead a
	// generation-match precondition using its generation number.
	// attrs, err := dst.Attrs(ctx)
	// if err != nil {
	//      return fmt.Errorf("object.Attrs: %w", err)
	// }
	// dst = dst.If(storage.Conditions{GenerationMatch: attrs.Generation})

	if _, err := dst.CopierFrom(src).Run(ctx); err != nil {
			return fmt.Errorf("Object(%q).CopierFrom(%q).Run: %w", dstObject, srcObject, err)
	}
	fmt.Fprintf(w, "Blob %v in bucket %v copied to blob %v in bucket %v.\n", srcObject, srcBucket, dstObject, dstBucket)
	return nil
}