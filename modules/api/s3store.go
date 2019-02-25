package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/dgraph-io/dgo"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gocontrib/auth"
	"github.com/sergeyt/pandora/modules/dgraph"
	"github.com/sergeyt/pandora/modules/utils"
	log "github.com/sirupsen/logrus"
)

// NewFileStore creates new ObjectStore
func NewFileStore() ObjectStore {
	return NewS3Store()
}

// NewS3Store creates new S3 ObjectStore
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
	Download(ctx context.Context, path string, w io.Writer) error
	Upload(ctx context.Context, path, mediaType string, r io.ReadCloser) (map[string]interface{}, error)
	Delete(ctx context.Context, path string) (string, interface{}, error)
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
func (fs *S3Store) Download(ctx context.Context, path string, w io.Writer) error {
	s := session.New(fs.config)
	d := s3manager.NewDownloader(s)
	_, err := d.DownloadWithContext(ctx, &s3Writer{w}, &s3.GetObjectInput{
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
func (fs *S3Store) Upload(ctx context.Context, path, mediaType string, r io.ReadCloser) (map[string]interface{}, error) {
	sess := session.New(fs.config)
	up := s3manager.NewUploader(sess, func(u *s3manager.Uploader) {
		u.PartSize = 5 * 1024 * 1024
		u.LeavePartsOnError = true
	})
	out, err := up.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket: aws.String(fs.bucket),
		Key:    aws.String(path),
		Body:   r,
	})

	if err != nil {
		log.Errorf("s3.Upload fail: %v", err)
		return nil, err
	}

	client, err := dgraph.NewClient()
	if err != nil {
		return nil, err
	}

	tx := client.NewTxn()
	defer tx.Discard(ctx)

	id, err := findFile(ctx, tx, path)
	if err != nil {
		return nil, err
	}

	user := auth.GetContextUser(ctx)

	in := make(utils.OrderedJSON)
	if len(id) > 0 {
		in["uid"] = id
	}
	in["path"] = path
	in["url"] = out.Location
	in["upload_id"] = out.UploadID
	in["version_id"] = out.VersionID
	in["content_type"] = mediaType

	label := dgraph.NodeLabel("file")
	i := strings.Index(mediaType, "/")
	if i >= 0 {
		label = dgraph.NodeLabel(mediaType[0:i])
	}

	results, err := dgraph.Mutate(ctx, tx, dgraph.Mutation{
		Input:     in,
		NodeLabel: label,
		ID:        id,
		By:        user.GetID(),
	})
	if err != nil {
		return nil, err
	}
	if len(results) != 1 {
		return nil, fmt.Errorf("unexpected mutation results: %v", results)
	}

	return results[0], nil
}

// Delete object at given path
func (fs *S3Store) Delete(ctx context.Context, path string) (string, interface{}, error) {
	id := ""
	sess := session.New(fs.config)
	svc := s3.New(sess)
	_, err := svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(fs.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		log.Errorf("s3.Delete fail: %v", err)
		return id, nil, err
	}

	client, err := dgraph.NewClient()
	if err != nil {
		return id, nil, err
	}

	tx := client.NewTxn()
	defer tx.Discard(ctx)

	id, err = findFile(ctx, tx, path)
	if err != nil {
		return id, nil, err
	}

	resp2, err := dgraph.DeleteNode(ctx, tx, id)
	if err != nil {
		return id, nil, err
	}

	return id, resp2, nil
}

func findFile(ctx context.Context, tx *dgo.Txn, path string) (string, error) {
	query := `query file($path: string) {
		files(func: eq(path, $path)) {
			uid
		}
	  }`
	resp, err := tx.QueryWithVars(ctx, query, map[string]string{
		"$path": path,
	})
	if err != nil {
		log.Errorf("dgraph.Txn.Mutate fail: %v", err)
		return "", err
	}

	var result struct {
		Files []struct {
			ID string `json:"uid"`
		} `json:"files"`
	}
	err = json.Unmarshal(resp.GetJson(), &result)
	if err != nil {
		log.Errorf("json.Unmarshal fail: %v", err)
		return "", err
	}

	if len(result.Files) == 0 {
		return "", nil
	}

	if len(result.Files) > 1 {
		return "", fmt.Errorf("inconsistent db state: found multiple file nodes")
	}

	id := result.Files[0].ID
	return id, nil
}
