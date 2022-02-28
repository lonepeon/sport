package s3

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var (
	ErrGeneric = errors.New("something wrong happened")
)

type Bucket struct {
	name   string
	region string

	// Endpoint is an optional endpoint URL (hostname only or fully qualified URI) that overrides the default one from s3
	Endpoint   string
	HTTPClient *http.Client
}

// NewBucket represents an S3 bucket. Credentials are loaded using environment variables:
// - Access Key ID:     AWS_ACCESS_KEY_ID or AWS_ACCESS_KEY
// - Secret Access Key: AWS_SECRET_ACCESS_KEY or AWS_SECRET_KEY
func NewBucket(name string, region string) *Bucket {
	return &Bucket{
		name:       name,
		region:     region,
		HTTPClient: http.DefaultClient,
	}
}

func (b *Bucket) StoreAsset(r io.Reader, dest string) error {
	sess, err := b.openSession()
	if err != nil {
		return fmt.Errorf("can't initialize s3 client: %w: %v", ErrGeneric, err)
	}

	uploader := s3manager.NewUploader(sess)
	input := s3manager.UploadInput{Bucket: aws.String(b.name), Key: aws.String(dest), Body: r}
	if _, err := uploader.Upload(&input); err != nil {
		return fmt.Errorf("can't upload file: %w: %v", ErrGeneric, err)
	}

	return nil
}

func (b *Bucket) DeleteAsset(path string) error {
	sess, err := b.openSession()
	if err != nil {
		return fmt.Errorf("can't initialize s3 client: %w: %v", ErrGeneric, err)
	}

	svc := s3.New(sess)
	if _, err := svc.DeleteObject(&s3.DeleteObjectInput{Bucket: &b.name, Key: aws.String(path)}); err != nil {
		return fmt.Errorf("can't delete s3 file: %v", err)
	}

	return nil
}

func (b *Bucket) openSession() (*session.Session, error) {
	return session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Credentials:      credentials.NewEnvCredentials(),
			Region:           aws.String(b.region),
			Endpoint:         &b.Endpoint,
			HTTPClient:       b.HTTPClient,
			S3ForcePathStyle: aws.Bool(true),
		},
	})
}
