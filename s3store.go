package main

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
	"github.com/spf13/viper"
)

func makeFileStore() (FileStore, error) {
	return &S3Store{
		config: &aws.Config{
			Endpoint:    aws.String(viper.GetString("s3.endpoint")),
			Region:      aws.String(viper.GetString("s3.region")),
			Credentials: credentials.NewEnvCredentials(),
		},
		bucket: getBucket(),
	}, nil
}

func getBucket() string {
	return os.Getenv("S3_BUCKET")
}

type FileStore interface {
	Download(c context.Context, path string, w io.Writer) error
	Upload(c context.Context, path string, r io.ReadCloser) (interface{}, error)
	Delete(c context.Context, path string) error
}

type S3Store struct {
	config *aws.Config
	bucket string
}

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

func (fs *S3Store) Delete(c context.Context, path string) error {
	sess := session.New(fs.config)
	svc := s3.New(sess)
	// TODO how to delete by key?
	iter := s3manager.NewDeleteListIterator(svc, &s3.ListObjectsInput{
		Bucket: aws.String(fs.bucket),
		// Key:    path,
	})
	err := s3manager.NewBatchDeleteWithClient(svc).Delete(c, iter)
	if err != nil {
		log.Errorf("s3.Delete fail: %v", err)
	}
	return err
}
