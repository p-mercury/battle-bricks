package filekit

import (
	"contracts/dist/filestaging/v1/filestagingconnect"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Client struct {
	fileStaging filestagingconnect.InternalClient
	transport   *http.Client
	s3          *s3.Client
	bucketName  string
}

type NewClientInput struct {
	FileStaging filestagingconnect.InternalClient
	Transport   *http.Client
	S3          *s3.Client
	BucketName  string
}

func NewClient(params *NewClientInput) *Client {

	transport := http.DefaultClient
	if params.Transport != nil {
		transport = params.Transport
	}

	return new(Client{
		fileStaging: params.FileStaging,
		transport:   transport,
		s3:          params.S3,
		bucketName:  params.BucketName,
	})
}
