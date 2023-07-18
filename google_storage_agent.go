package main

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"io/ioutil"
	"regexp"
)

type GoogleStorageAgent struct {
	BucketName   string
	ObjectFilter string
}

func (a *GoogleStorageAgent) Execute(ctx context.Context, task *Task) error {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create client: %v", err)
	}

	bucket := client.Bucket(a.BucketName)

	// Compile regex
	r, err := regexp.Compile(a.ObjectFilter)
	if err != nil {
		return fmt.Errorf("failed to compile regex: %v", err)
	}

	// Iterate over objects in the bucket
	it := bucket.Objects(ctx, nil)
	for {
		attrs, err := it.Next()
		//if err == storage.Done {
		//	break
		//}
		if err != nil {
			return fmt.Errorf("failed to list objects: %v", err)
		}

		// Match object name with regex
		if r.MatchString(attrs.Name) {
			// Download object
			rc, err := bucket.Object(attrs.Name).NewReader(ctx)
			if err != nil {
				return fmt.Errorf("failed to create object reader: %v", err)
			}

			data, err := ioutil.ReadAll(rc)
			rc.Close()
			if err != nil {
				return fmt.Errorf("failed to read object: %v", err)
			}

			// Print the object data
			fmt.Println(string(data))
		}
	}

	return nil
}
