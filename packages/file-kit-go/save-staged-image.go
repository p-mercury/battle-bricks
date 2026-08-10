package filekit

import (
	"bytes"
	"context"
	"fmt"
	"contracts/dist/filestaging/v1"
	"io"
	"mime"
	"net/http"
	"slices"
	"strings"

	"connectrpc.com/connect"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3Types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/disintegration/imaging"
)

type SaveStagedImageInput struct {
	Id      string
	Targets []SaveStagedImageTarget
}

type SaveStagedImageTarget struct {
	Key    string
	Width  int
	Height int
	Anchor imaging.Anchor
	Filter imaging.ResampleFilter
}

func (client *Client) SaveStagedImage(ctx context.Context, params *SaveStagedImageInput) error {
	var file *filestaging.File
	{
		rep, err := client.fileStaging.GetFile(ctx, connect.NewRequest(
			&filestaging.GetFileRequest{
				Id: params.Id,
			},
		))
		if err != nil {
			return err
		}

		file = rep.Msg.File
	}

	if file.Url == nil {
		return fmt.Errorf("missing file url")
	}

	{
		if file.ContentType == nil {
			return fmt.Errorf("Unsupported content type")
		}
		mediaType, _, err := mime.ParseMediaType(*file.ContentType)
		if err != nil {
			mediaType, _, _ = strings.Cut(*file.ContentType, ";")
			mediaType = strings.TrimSpace(mediaType)
		}
		mediaType = strings.ToLower(mediaType)
		if !slices.Contains([]string{"image/png", "image/jpeg", "image/gif"}, mediaType) {
			return fmt.Errorf("unsupported content type")
		}
	}

	resp, err := client.transport.Get(*file.Url)
	if err != nil {
		return fmt.Errorf("failed to download from presigned URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file, status: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	img, err := imaging.Decode(bytes.NewReader(data), imaging.AutoOrientation(true))
	if err != nil {
		return fmt.Errorf("failed to decode image: %w", err)
	}

	for _, target := range params.Targets {
		generated := imaging.Fill(
			img,
			target.Width,
			target.Height,
			target.Anchor,
			target.Filter,
		)

		{
			buf := new(bytes.Buffer)
			if err := imaging.Encode(buf, generated, imaging.JPEG, imaging.JPEGQuality(90)); err != nil {
				return fmt.Errorf("failed to encode large image: %w", err)
			}
			len := int64(buf.Len())
			_, err = client.s3.PutObject(ctx, &s3.PutObjectInput{
				Bucket:            &client.bucketName,
				Key:               new(target.Key),
				ContentType:       new("image/jpeg"),
				ACL:               s3Types.ObjectCannedACLBucketOwnerFullControl,
				ChecksumAlgorithm: s3Types.ChecksumAlgorithmSha256,
				Body:              buf,
				ContentLength:     &len,
			})
			if err != nil {
				return fmt.Errorf("failed to upload to bucket %s: %w", client.bucketName, err)
			}
		}
	}

	return nil
}
