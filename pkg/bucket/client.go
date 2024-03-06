package bucket

import (
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Bucket struct {
	s3Srv  *s3.S3
	bucket string
}

// NewBucket return a new bucket instance
func NewBucket(id, secret, bucket, region string) (*Bucket, error) {
	cfg := aws.NewConfig()
	session, err := session.NewSession(
		cfg.WithRegion(region),
		cfg.WithCredentials(credentials.NewStaticCredentials(
			id,
			secret,
			"",
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	srv := s3.New(session)
	return &Bucket{srv, bucket}, nil
}

// PutImage adds an image to a bucket.
func (b *Bucket) PutImage(filename string, content io.ReadCloser) error {
	input := new(s3.PutObjectInput)
	input.SetBucket(b.bucket)
	input.SetKey(fmt.Sprintf("imgo/%s.png", filename))
	input.SetBody(aws.ReadSeekCloser(content))

	_, err := b.s3Srv.PutObject(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			return errors.New("aws err: " + aerr.Error())
		}
		return fmt.Errorf("failed to put object: %w", err)
	}

	return nil
}

// GetImageURL get url of image from bucket.
// url is temporarily public and it has expiration time of 20 second.
func (b *Bucket) GetImageURL(filename string) (string, error) {
	input := new(s3.GetObjectInput)
	input.SetBucket(b.bucket)
	input.SetKey(fmt.Sprintf("imgo/%s.png", filename))

	req, _ := b.s3Srv.GetObjectRequest(input)
	url, err := req.Presign(20 * time.Second)
	if err != nil {
		return "", err
	}

	return url, nil
}
