package filekit

import (
	"context"
	"contracts/dist/filestaging/v1"
	"fmt"
	"mime"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3Types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
)

type SaveStagedFileInput struct {
	File *filestaging.File
	Path string
	Tags map[string]string
}

func (client *Client) SaveStagedFile(ctx context.Context, params *SaveStagedFileInput) (string, error) {
	if params.File == nil {
		return "", fmt.Errorf("missing file")
	}

	if params.File.Url == nil {
		return "", fmt.Errorf("missing file url")
	}

	if params.File.ContentType == nil {
		return "", fmt.Errorf("unsupported content type")
	}

	var contentType string
	{
		mediaType, _, err := mime.ParseMediaType(*params.File.ContentType)
		if err != nil {
			mediaType, _, _ = strings.Cut(*params.File.ContentType, ";")
			mediaType = strings.TrimSpace(mediaType)
		}
		mediaType = strings.ToLower(mediaType)
		if !slices.Contains([]string{
			"image/png",
			"image/gif",
			"image/jpeg",
			"application/pdf",
			"application/zip",
			"application/gzip",
			"application/x-7z-compressed",
			"application/json",
			"application/msword",
			"application/vnd.ms-excel",
			"application/xml",
			"text/plain",
			"text/markdown",
			"text/csv",
		}, mediaType) {
			return "", fmt.Errorf("unsupported content type")
		}
		contentType = mediaType
	}

	resp, err := client.transport.Get(*params.File.Url)
	if err != nil {
		return "", fmt.Errorf("failed to download from presigned URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download file, status: %s", resp.Status)
	}

	key := strings.Trim(params.Path, "/") + "/" + uuid.New().String() + "/" + params.File.Name

	var tagging *string
	if len(params.Tags) > 0 {
		if len(params.Tags) > 10 {
			return "", fmt.Errorf("too many tags: max 10 allowed, got %d", len(params.Tags))
		}
		values := url.Values{}
		for k, v := range params.Tags {
			if k == "" {
				return "", fmt.Errorf("tag key cannot be empty")
			}
			if len(k) > 128 {
				return "", fmt.Errorf("tag key %q exceeds 128 characters", k)
			}
			if len(v) > 256 {
				return "", fmt.Errorf("tag value for key %q exceeds 256 characters", k)
			}
			values.Set(k, v)
		}
		tagging = new(values.Encode())
	}

	_, err = client.s3.PutObject(ctx, &s3.PutObjectInput{
		Bucket:         &client.bucketName,
		Key:            &key,
		ContentType:    &contentType,
		ACL:            s3Types.ObjectCannedACLBucketOwnerFullControl,
		ChecksumSHA256: &params.File.Checksum,
		Body:           resp.Body,
		ContentLength:  &resp.ContentLength,
		Tagging:        tagging,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to bucket %s: %w", client.bucketName, err)
	}

	return key, nil
}
