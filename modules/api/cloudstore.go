package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/araddon/gou"
	dgo "github.com/dgraph-io/dgo/v2"

	"github.com/sergeyt/pandora/modules/config"
	"github.com/sergeyt/pandora/modules/dgraph"
	"github.com/sergeyt/pandora/modules/utils"
	log "github.com/sirupsen/logrus"

	"github.com/lytics/cloudstorage"
	"github.com/lytics/cloudstorage/awss3"
)

// ObjectStore is store of any BLOB objects
type ObjectStore interface {
	Download(ctx context.Context, id string, w io.Writer) error
	DownloadFile(ctx context.Context, file *FileInfo, w io.Writer) error
	Upload(ctx context.Context, path, mediaType string, r io.ReadCloser) (map[string]interface{}, error)
	Delete(ctx context.Context, id string) (string, interface{}, error)
	DeleteObject(ctx context.Context, path string) error
}

// NewCloudStore creates new cloud store based on environment variables
func NewCloudStore() *CloudStore {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-east-1"
	}
	// TODO log level
	config := &cloudstorage.Config{
		BaseUrl:    os.Getenv("AWS_S3_ENDPOINT"),
		Type:       awss3.StoreType,
		AuthMethod: awss3.AuthAccessKey,
		Bucket:     os.Getenv("AWS_S3_BUCKET"),
		TmpDir:     "/tmp/localcache/aws",
		Settings:   make(gou.JsonHelper),
		Region:     region,
	}
	config.Settings[awss3.ConfKeyAccessKey] = os.Getenv("AWS_ACCESS_KEY_ID")
	config.Settings[awss3.ConfKeyAccessSecret] = os.Getenv("AWS_SECRET_ACCESS_KEY")
	config.Settings[awss3.ConfKeyDisableSSL] = "true"
	return &CloudStore{
		config: config,
	}
}

// CloudStore is Amazon S3 ObjectStore
type CloudStore struct {
	config *cloudstorage.Config
}

// EnsureBucket creates AWS_S3_BUCKET
func (s *CloudStore) EnsureBucket() error {
	return nil
}

// Download object at given path
func (s *CloudStore) Download(ctx context.Context, id string, w io.Writer) error {
	file, err := findFileTx(ctx, id)
	if err != nil {
		return err
	}

	if file == nil {
		return fmt.Errorf("file not found: %s", id)
	}

	return s.DownloadFile(ctx, file, w)
}

// DownloadFile downloads given file
func (s *CloudStore) DownloadFile(ctx context.Context, file *FileInfo, w io.Writer) error {
	if file == nil {
		return fmt.Errorf("file not found")
	}

	path := file.Path
	store, err := s.newStore()
	if err != nil {
		return err
	}

	src, err := store.NewReaderWithContext(ctx, path)
	if err != nil {
		log.Errorf("store.NewReaderWithContext fail: %v", err)
		return err
	}

	return copy(src, w)
}

func (s *CloudStore) newStore() (cloudstorage.Store, error) {
	store, err := cloudstorage.NewStore(s.config)
	if err != nil {
		log.Errorf("cloudstorage.NewStore fail: %v", err)
		return nil, err
	}
	log.Info("!!!CloudStore is created!!!")
	return store, nil
}

func copy(r io.Reader, w io.Writer) error {
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		if _, err := w.Write(buf[:n]); err != nil {
			return err
		}
	}
	return nil
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
func (s *CloudStore) Upload(ctx context.Context, path, mediaType string, r io.ReadCloser) (map[string]interface{}, error) {
	store, err := s.newStore()
	if err != nil {
		return nil, err
	}

	obj, err := store.Get(ctx, path)
	if err != nil {
		// log.Errorf("store.Get fail: %v", err)
		// return nil, err
		obj, err = store.NewObject(path)
	}

	f, err := obj.Open(cloudstorage.ReadWrite)
	if err != nil {
		log.Errorf("cloudstorage.Object.Open fail: %v", err)
		return nil, err
	}

	// meta := map[string]string{}
	// w, err := store.NewWriterWithContext(ctx, path, meta, cloudstorage.Opts{
	// 	IfNotExists: true,
	// })
	// if err != nil {
	// 	log.Errorf("store.NewWriterWithContext fail: %v", err)
	// 	return nil, err
	// }

	err = copy(r, f)
	if err != nil {
		log.Errorf("store.Upload fail: %v", err)
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
func (s *CloudStore) Delete(ctx context.Context, id string) (string, interface{}, error) {
	store, err := s.newStore()
	if err != nil {
		return "", nil, err
	}

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

	err = store.Delete(ctx, file.Path)
	if err != nil {
		log.Errorf("store.Delete fail: %v", err)
		return "", nil, err
	}

	resp, err := dgraph.DeleteNode(ctx, tx, file.ID)
	if err != nil {
		return file.ID, nil, err
	}

	return file.ID, resp, nil
}

// DeleteObject by given path
func (s *CloudStore) DeleteObject(ctx context.Context, path string) error {
	store, err := s.newStore()
	if err != nil {
		return nil
	}
	err = store.Delete(ctx, path)
	if err != nil {
		log.Errorf("store.Delete fail: %v", err)
		return err
	}
	return nil
}

// FileInfo descriptor
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
