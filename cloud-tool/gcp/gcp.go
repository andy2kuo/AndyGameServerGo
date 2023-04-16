package gcp

import (
	"context"
	"os"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

var storage_client *storage.Client

func Init(json_file_name string) {
	ctx := context.Background()
	var err error

	var json_file []byte
	json_file, err = os.ReadFile(json_file_name)
	if err != nil {
		panic(err)
	}

	storage_client, err = storage.NewClient(ctx, option.WithCredentialsJSON(json_file))
	if err != nil {
		panic(err)
	}
}
