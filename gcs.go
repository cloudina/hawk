package main

type GCS_Manager struct {
	*BucketMgr
}

func (self *GCS_Manager ) bucketExists(bucket string) (bool, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
			return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	bucket := client.Bucket(bucket)
	exists,err := bucket.Attrs(ctx)
	if err != nil {
		log.Fatalf("Message: %v",err)
		return false, err
	} else {
		return true, nil
	}
}

func (self *GCS_Manager ) keyExists(bucket string, key string) (bool, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
			return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	bucket := client.Bucket(bucket)
	object = bucket.Object(key)
	exists,err := object.Attrs(ctx)
	if err != nil {
		log.Fatalf("Message: %v",err)
		return false, err
	} else {
		return true, nil
	}
}

func (self *GCS_Manager ) readFile(bucket string, item string) ([] byte, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
			return nil, fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	rc, err := client.Bucket(bucket).Object(item).NewReader(ctx)
	if err != nil {
			return nil, fmt.Errorf("Object(%q).NewReader: %w", object, err)
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
			return nil, fmt.Errorf("ioutil.ReadAll: %w", err)
	}
	fmt.Fprintf(w, "Blob %v downloaded.\n", object)
	return data, nil
}

func (self *GCS_Manager ) copyFile(srcBucket string, srcObject string, other string) (error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
			return fmt.Errorf("storage.NewClient: %w", err)
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
			return fmt.Errorf("Object(%q).CopierFrom(%q).Run: %w", dstObject, srcObject, err)
	}
	fmt.Fprintf(w, "Blob %v in bucket %v copied to blob %v in bucket %v.\n", srcObject, srcBucket, dstObject, dstBucket)
	return nil
}

func (self *GCS_Manager ) deleteFile(bucket string, item string) (error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
			return fmt.Errorf("storage.NewClient: %w", err)
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
			return fmt.Errorf("object.Attrs: %w", err)
	}
	o = o.If(storage.Conditions{GenerationMatch: attrs.Generation})

	if err := o.Delete(ctx); err != nil {
			return fmt.Errorf("Object(%q).Delete: %w", item, err)
	}
	fmt.Fprintf(w, "Blob %v deleted.\n", item)
	return nil
}
