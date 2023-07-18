package main

import (
	"cloud.google.com/go/storage"
	vision "cloud.google.com/go/vision/apiv1"
	"context"
	"fmt"
	pb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

type GoogleOCRStorageAgent struct {
	BucketName string
	ObjectName string
}

func (a *GoogleOCRStorageAgent) Execute(ctx context.Context, task *Task) error {
	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create client: %v", err)
	}
	defer client.Close()

	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create storage client: %v", err)
	}
	defer storageClient.Close()

	rc, err := storageClient.Bucket(a.BucketName).Object(a.ObjectName).NewReader(ctx)
	if err != nil {
		return fmt.Errorf("failed to read object: %v", err)
	}
	defer rc.Close()

	image, err := vision.NewImageFromReader(rc)
	if err != nil {
		return fmt.Errorf("failed to create image: %v", err)
	}

	resp, err := client.DetectTexts(ctx, image, &pb.ImageContext{}, 10)
	if err != nil {
		return fmt.Errorf("failed to detect texts: %v", err)
	}

	for _, annotation := range resp {
		fmt.Println(annotation.Description)
	}

	return nil
}
