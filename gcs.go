package main

import (
	"context"
	"errors"
	"path"
	"time"
	"io/ioutil"
	"cloud.google.com/go/storage"
)

type GCS_Manager struct {
	BucketMgr
}

func (self *GCS_Manager) bucketExists(bucket string) (bool, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)

	if err != nil {
		elog.Println(time.Now().Format(time.RFC3339) + "bucketExists: storage.NewClient " + bucket + " error : " + err.Error())
		return false, errors.New("storage.NewClient Failed ")
	}
	defer client.Close()

	bucketObj := client.Bucket(bucket)
	_, err = bucketObj.Attrs(ctx)

	if err != nil {
		elog.Println(time.Now().Format(time.RFC3339) + "bucketExists: bucketObj.Attrs " + bucket + " error : " + err.Error())
		return false, errors.New("bucketObj.Attrs Failed ")
	} else {
		return true, nil
	}
}

func (self *GCS_Manager) keyExists(bucket string, key string) (bool, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)

	if err != nil {
		elog.Println(time.Now().Format(time.RFC3339) + "keyExists: storage.NewClient " + bucket + " error : " + err.Error())
		return false, errors.New("storage.NewClient Failed ")
	}
	defer client.Close()

	bucketObj := client.Bucket(bucket)
	object := bucketObj.Object(key)
	_, err = object.Attrs(ctx)

	if err != nil {
		elog.Println(time.Now().Format(time.RFC3339) + "keyExists: object.Attrs " + bucket + " error : " + err.Error())
		return false, errors.New("object.Attrs Failed ")
	} else {
		return true, nil
	}
}

func (self *GCS_Manager) readFile(bucket string, item string) ([]byte, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		elog.Println(time.Now().Format(time.RFC3339) + "readFile: storage.NewClient " + bucket + " key " +item +  " error : " + err.Error())
		return nil, errors.New("storage.NewClient Failed ")
	}
	defer client.Close()

	rc, err := client.Bucket(bucket).Object(item).NewReader(ctx)
	if err != nil {
		elog.Println(time.Now().Format(time.RFC3339) + "readFile: client.Bucket " + bucket + " key " +item +  " error : " + err.Error())
		return nil, errors.New("Bucket.Object.NewReader Failed ")
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		elog.Println(time.Now().Format(time.RFC3339) + "readFile: client.Bucket " + bucket + " key " +item +  " error : " + err.Error())
		return nil, errors.New("ioutil ReadAll Failed ")
	}

	info.Println(time.Now().Format(time.RFC3339) + " Downloaded object " + item + " from bucket " + bucket)

	return data, nil
}

func (self *GCS_Manager) copyFile(srcBucket string, srcObject string, other string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		elog.Println(time.Now().Format(time.RFC3339) + "copyFile: storage.NewClient " + srcBucket + " error : " + err.Error())
		return errors.New("storage.NewClient Failed ")
	}
	defer client.Close()

	dstBucket := path.Dir(other)
	dstObject := path.Base(other)

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
		elog.Println(time.Now().Format(time.RFC3339) + " Unable to copy object " + srcObject + " from bucket " + srcBucket + " to bucket " + dstBucket + " error : " + err.Error())
		return errors.New("Unable to copy file")
	}

	return nil
}

func (self *GCS_Manager) deleteFile(bucket string, item string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		elog.Println(time.Now().Format(time.RFC3339) + "deleteFile: storage.NewClient " + bucket + " error : " + err.Error())
		return errors.New("storage.NewClient Failed ")
	}
	defer client.Close()

	//ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	//defer cancel()

	o := client.Bucket(bucket).Object(item)

	// Optional: set a generation-match precondition to avoid potential race
	// conditions and data corruptions. The request to delete the file is aborted
	// if the object's generation number does not match your precondition.
	attrs, err := o.Attrs(ctx)
	if err != nil {
		elog.Println(time.Now().Format(time.RFC3339) + "deleteFile: bucketObj.Attrs " + bucket + " error : " + err.Error())
		return errors.New("bucketObj.Attrs Failed ")	
	}
	o = o.If(storage.Conditions{GenerationMatch: attrs.Generation})

	if err := o.Delete(ctx); err != nil {
		elog.Println(time.Now().Format(time.RFC3339) + "deleteFile: Delete " + bucket + " object "+item + " error : " + err.Error())
		return errors.New("Object Delete Failed ")	
	}
	info.Println(time.Now().Format(time.RFC3339) + " deleteFile " + item + " successfully deleted from bucket " + bucket)
	return nil
}
