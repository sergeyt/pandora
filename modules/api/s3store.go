package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	dgo "github.com/dgraph-io/dgo/v2"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/sergeyt/pandora/modules/config"
	"github.com/sergeyt/pandora/modules/dgraph"
	"github.com/sergeyt/pandora/modules/utils"
	log "github.com/sirupsen/logrus"
)

// InitStore initializes file store
func InitStore() {
	fs := NewS3Store()
	err := fs.EnsureBucket()
	if err != nil {
		log.Errorf("s3.EnsureBucket fail: %v", err)
	}
}

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
		Endpoint:         aws.String(os.Getenv("AWS_S3_ENDPOINT")),
		Region:           aws.String(os.Getenv("AWS_REGION")),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
		LogLevel:         &logLevel,
	}
}

func getBucket() string {
	return os.Getenv("AWS_S3_BUCKET")
}

// ObjectStore is store of any BLOB objects
type ObjectStore interface {
	Download(ctx context.Context, id string, w io.Writer) error
	DownloadFile(ctx context.Context, file *FileInfo, w io.Writer) error
	Upload(ctx context.Context, path, mediaType string, r io.ReadCloser) (map[string]interface{}, error)
	Delete(ctx context.Context, id string) (string, interface{}, error)
	DeleteObject(ctx context.Context, path string) error
}

// S3Store is Amazon S3 ObjectStore
type S3Store struct {
	config *aws.Config
	bucket string
}

// EnsureBucket creates AWS_S3_BUCKET
func (fs *S3Store) EnsureBucket() error {
	sess, err := session.NewSession(fs.config)
	if err != nil {
		log.Errorf("aws.session.NewSession fail: %v", err)
		return err
	}
	svc := s3.New(sess)
	_, err = svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(fs.bucket),
	})
	return err
}

// Download object at given path
func (fs *S3Store) Download(ctx context.Context, id string, w io.Writer) error {
	file, err := findFileTx(ctx, id)
	if err != nil {
		return err
	}

	if file == nil {
		return fmt.Errorf("file not found: %s", id)
	}

	return fs.DownloadFile(ctx, file, w)
}

// DownloadFile downloads given file
func (fs *S3Store) DownloadFile(ctx context.Context, file *FileInfo, w io.Writer) error {
	if file == nil {
		return fmt.Errorf("file not found")
	}

	path := file.Path
	s, err := session.NewSession(fs.config)
	if err != nil {
		log.Errorf("aws.seesion.NewSession fail: %v", err)
		return err
	}
	d := s3manager.NewDownloader(s)
	_, err = d.DownloadWithContext(ctx, &s3Writer{w, 0}, &s3.GetObjectInput{
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
	offset int64
}

func (w *s3Writer) WriteAt(p []byte, off int64) (n int, err error) {
	if w.offset == off {
		n, err := w.Write(p)
		if err != nil {
			return n, err
		}
		w.offset += int64(n)
		return n, err
	}
	return 0, fmt.Errorf("write at any offset is not supported")
}

// Upload object at given path
func (fs *S3Store) Upload(ctx context.Context, path, mediaType string, r io.ReadCloser) (map[string]interface{}, error) {
	sess, err := session.NewSession(fs.config)
	if err != nil {
		log.Errorf("aws.session.NewSession fail: %v", err)
		return nil, err
	}
	up := s3manager.NewUploader(sess, func(u *s3manager.Uploader) {
		u.PartSize = 5 * 1024 * 1024
		u.LeavePartsOnError = true
	})
	_, err = up.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket: aws.String(fs.bucket),
		Key:    aws.String(path),
		Body:   r,
	})

	if err != nil {
		log.Errorf("s3.Upload fail: %v", err)
		return nil, err
	}

	dc, err := dgraph.NewClient()
	if err != nil {
		return nil, err
	}

	tx := dc.NewTxn()
	defer tx.Discard(ctx)

	file, err := findFile(ctx, tx, path)
	if err != nil {
		return nil, err
	}

	if file == nil {
		file = &FileInfo{}
	}
	file.Path = path
	file.MediaType = mediaType

	return addFile(ctx, tx, file)
}

func addFile(ctx context.Context, tx *dgo.Txn, file *FileInfo) (map[string]interface{}, error) {
	in := make(utils.OrderedJSON)
	id := file.ID
	if len(id) > 0 {
		in["uid"] = id
	}

	in["content_type"] = file.MediaType

	if len(file.URL) > 0 {
		in["url"] = file.URL
	}

	if len(file.Path) > 0 {
		in["path"] = file.Path
	}

	if len(file.URL) == 0 && len(file.Path) > 0 {
		baseURL := config.ServerURL()
		in["url"] = fmt.Sprintf("%s/api/file/%s", baseURL, file.Path)
	}

	discardHere := false
	if tx == nil {
		dc, err := dgraph.NewClient()
		if err != nil {
			return nil, err
		}

		tx = dc.NewTxn()
	}

	defer func() {
		if discardHere {
			tx.Discard(ctx)
		}
	}()

	results, err := dgraph.Mutate(ctx, tx, dgraph.Mutation{
		Input:     in,
		NodeLabel: dgraph.NodeLabel("file"),
		ID:        id,
	})
	if err != nil {
		return nil, err
	}
	if len(results) != 1 {
		return nil, fmt.Errorf("unexpected mutation results: %v", results)
	}

	return results[0], nil
}

// Delete object by given path or file id
func (fs *S3Store) Delete(ctx context.Context, id string) (string, interface{}, error) {
	dc, err := dgraph.NewClient()
	if err != nil {
		return "", nil, err
	}

	tx := dc.NewTxn()
	defer tx.Discard(ctx)

	file, err := findFile(ctx, tx, id)
	if err != nil {
		return "", nil, err
	}

	if file == nil {
		return "", nil, fmt.Errorf("file not found: %s", id)
	}

	path := file.Path
	id = file.ID
	sess, err := session.NewSession(fs.config)
	if err != nil {
		log.Errorf("aws.session.NewSession fail: %v", err)
		return "", nil, err
	}

	svc := s3.New(sess)
	_, err = svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(fs.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		log.Errorf("s3.Delete fail: %v", err)
		return id, nil, err
	}

	resp2, err := dgraph.DeleteNode(ctx, tx, id)
	if err != nil {
		return file.ID, nil, err
	}

	return file.ID, resp2, nil
}

// DeleteObject by given path
func (fs *S3Store) DeleteObject(ctx context.Context, path string) error {
	sess, err := session.NewSession(fs.config)
	if err != nil {
		log.Errorf("aws.session.NewSession fail: %v", err)
		return err
	}
	svc := s3.New(sess)
	_, err = svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(fs.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		log.Errorf("s3.Delete fail: %v", err)
		return err
	}
	return nil
}

type FileInfo struct {
	ID        string `json:"uid"`
	Path      string `json:"path"`
	URL       string `json:"url"`
	MediaType string `json:"content_type"`
}

func findFile(ctx context.Context, tx *dgo.Txn, id string) (*FileInfo, error) {
	filter := "eq(path, $id)"
	if dgraph.IsUID(id) {
		filter = "uid($id)"
	}

	query := fmt.Sprintf(`query file($id: string) {
		files(func: %s) {
			uid
			path
		}
	  }`, filter)
	resp, err := tx.QueryWithVars(ctx, query, map[string]string{
		"$id": id,
	})
	if err != nil {
		log.Errorf("dgraph.Txn.Mutate fail: %v", err)
		return nil, err
	}

	var result struct {
		Files []FileInfo `json:"files"`
	}
	err = json.Unmarshal(resp.GetJson(), &result)
	if err != nil {
		log.Errorf("json.Unmarshal fail: %v", err)
		return nil, err
	}

	if len(result.Files) == 0 {
		return nil, nil
	}

	if len(result.Files) > 1 {
		return nil, fmt.Errorf("inconsistent db state: found multiple file nodes")
	}

	file := result.Files[0]
	return &file, nil
}

func findFileTx(ctx context.Context, id string) (*FileInfo, error) {
	dc, err := dgraph.NewClient()
	if err != nil {
		return nil, err
	}

	tx := dc.NewTxn()
	defer tx.Discard(ctx)

	return findFile(ctx, tx, id)
}
