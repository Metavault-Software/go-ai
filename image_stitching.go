package main

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"image"
)

type ImageStitchingAgent struct {
	BucketName  string
	ObjectNames []string
	OutputPath  string
}

func (a *ImageStitchingAgent) Execute(ctx context.Context, task *Task) error {
	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create storage client: %v", err)
	}
	defer storageClient.Close()

	images := make([]image.Image, len(a.ObjectNames))

	for i, objectName := range a.ObjectNames {
		rc, err := storageClient.Bucket(a.BucketName).Object(objectName).NewReader(ctx)
		if err != nil {
			return fmt.Errorf("failed to read object: %v", err)
		}

		img, _, err := image.Decode(rc)
		rc.Close()
		if err != nil {
			return fmt.Errorf("failed to decode image: %v", err)
		}

		images[i] = img
	}

	//// Stitch the images together horizontally
	//result := imaging.Collage(images, nil, len(images), 1)
	//
	//// Save the result
	//err = imaging.Save(result, a.OutputPath)
	//if err != nil {
	//	return fmt.Errorf("failed to save result: %v", err)
	//}

	return nil
}
