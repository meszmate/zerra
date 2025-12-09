package file

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/getsentry/sentry-go"
	"github.com/meszmate/zerra/internal/errx"
)

type Client struct {
	bucket string
	client *s3.Client
}

func NewClient(ctx context.Context, bucket string) (*Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	s3client := s3.NewFromConfig(cfg)

	return &Client{
		bucket: bucket,
		client: s3client,
	}, nil
}

func (c *Client) PutObject(ctx context.Context, key string, contentType string, body io.Reader) *errx.Error {
	_, err := c.client.PutObject(ctx, &s3.PutObjectInput{
		Key:         aws.String(key),
		Bucket:      aws.String(c.bucket),
		Body:        body,
		ContentType: aws.String(contentType),
	})

	if err != nil {
		sentry.CaptureException(err)
		return errx.InternalError()
	}

	return nil
}

func (c *Client) DeleteObject(ctx context.Context, key string) *errx.Error {
	_, err := c.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		sentry.CaptureException(err)
		return errx.InternalError()
	}

	return nil
}
