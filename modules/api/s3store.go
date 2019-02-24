package api

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	log "github.com/sirupsen/logrus"
)

func NewFileStore() ObjectStore {
	return NewS3Store()
}

func NewS3Store() *S3Store {
	return &S3Store{
		config: awsConfig(),
		bucket: getBucket(),
	}
}

func awsConfig() *aws.Config {
	logLevel := aws.LogDebug
	return &aws.Config{
		Credentials:      credentials.NewEnvCredentials(),
		Endpoint:         aws.String(os.Getenv("S3_ENDPOINT")),
		Region:           aws.String(os.Getenv("AWS_REGION")),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
		LogLevel:         &logLevel,
	}
}

func getBucket() string {
	return os.Getenv("S3_BUCKET")
}

// ObjectStore is store of any BLOB objects
type ObjectStore interface {
	Download(c context.Context, path string, w io.Writer) error
	Upload(c context.Context, path string, r io.ReadCloser) (interface{}, error)
	Delete(c context.Context, path string) (interface{}, error)
}

// S3Store is Amazon S3 ObjectStore
type S3Store struct {
	config *aws.Config
	bucket string
}

// EnsureBucket creates S3_BUCKET
func (fs *S3Store) EnsureBucket() error {
	sess := session.New(fs.config)
	svc := s3.New(sess)
	_, err := svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(fs.bucket),
	})
	return err
}

// Download object at given path
func (fs *S3Store) Download(c context.Context, path string, w io.Writer) error {
	s := session.New(fs.config)
	d := s3manager.NewDownloader(s)
	_, err := d.DownloadWithContext(c, &s3Writer{w}, &s3.GetObjectInput{
		Bucket: aws.String(fs.bucket),
		Key:    aws.String(path),
	})

	if err != nil {
		log.Errorf("s3.Download fail: %v", err)
	}

	return err
}

type s3Writer struct {
	io.Writer
}

func (w *s3Writer) WriteAt(p []byte, off int64) (n int, err error) {
	return 0, fmt.Errorf("not supported")
}

// Upload object at given path
func (fs *S3Store) Upload(c context.Context, path string, r io.ReadCloser) (interface{}, error) {
	sess := session.New(fs.config)
	up := s3manager.NewUploader(sess, func(u *s3manager.Uploader) {
		u.PartSize = 5 * 1024 * 1024
		u.LeavePartsOnError = true
	})
	out, err := up.UploadWithContext(c, &s3manager.UploadInput{
		Bucket: aws.String(fs.bucket),
		Key:    aws.String(path),
		Body:   r,
	})

	if err != nil {
		log.Errorf("s3.Upload fail: %v", err)
		return nil, err
	}

	return out, nil
}

// Delete object at given path
func (fs *S3Store) Delete(c context.Context, path string) (interface{}, error) {
	sess := session.New(fs.config)
	svc := s3.New(sess)
	out, err := svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(fs.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		log.Errorf("s3.Delete fail: %v", err)
		return nil, err
	}
	return out, nil
}
