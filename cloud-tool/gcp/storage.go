package gcp

import (
	"context"
)

func UploadFromMemory(ctx context.Context, bucket_name string, file_name string, data []byte, context_type string) (err error) {
	bkt := storage_client.Bucket(bucket_name)
	obj := bkt.Object(file_name)

	writer := obj.NewWriter(ctx)
	defer writer.Close()

	writer.ContentType = context_type
	_, err = writer.Write(data)
	return err
}
