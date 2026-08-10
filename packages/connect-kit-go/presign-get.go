package connectkit

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type PresignGetInput struct {
	Signer             *s3.PresignClient
	BucketName         *string
	Key                string
	Expires            time.Duration
	Name               *string
	ContentDisposition *string
}

func PresignGet(ctx context.Context, input PresignGetInput) (string, error) {
	getObjectAclInput := &s3.GetObjectInput{
		Bucket: input.BucketName,
		Key:    new(input.Key),
	}

	if input.ContentDisposition != nil {
		disposition := *input.ContentDisposition

		if input.Name != nil {
			disposition += `; filename="` + *input.Name + `"`
		}

		getObjectAclInput.ResponseContentDisposition = &disposition
	}

	presigned, err := input.Signer.PresignGetObject(ctx,
		getObjectAclInput,
		s3.WithPresignExpires(input.Expires))
	if err != nil {
		return "", err
	}

	return presigned.URL, nil
}
