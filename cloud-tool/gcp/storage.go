package gcp

import (
	"context"
)

type UploadRequest struct {
	BucketName  string
	FileName    string
	Data        []byte
	ContentType string
}

func UploadFromMemory(ctx context.Context, request UploadRequest) (err error) {
	bkt := storage_client.Bucket(request.BucketName)
	obj := bkt.Object(request.FileName)

	writer := obj.NewWriter(ctx)
	defer writer.Close()

	writer.ContentType = request.ContentType
	_, err = writer.Write(request.Data)

	return err
}
