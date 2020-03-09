package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"

	dgo "github.com/dgraph-io/dgo/v2"

	"github.com/graymeta/stow"
	_ "github.com/graymeta/stow/google"
	"github.com/graymeta/stow/s3"
	"github.com/sergeyt/pandora/modules/config"
	"github.com/sergeyt/pandora/modules/dgraph"
	"github.com/sergeyt/pandora/modules/utils"
	log "github.com/sirupsen/logrus"
)

// InitStore initializes file store
func InitStore() {
	fs := NewStow()
	err := fs.EnsureBucket()
	if err != nil {
		log.Errorf("EnsureBucket fail: %v", err)
	}
}

// NewCloudStore creates new instance of cloud store
func NewCloudStore() CloudStore {
	return NewStow()
}

func env(name, defval string) string {
	val := os.Getenv(name)
	if val == "" {
		val = defval
	}
	return val
}

// NewStow creates new Stow instance
func NewStow() *Stow {
	endpoint := os.Getenv("AWS_S3_ENDPOINT")
	id := env("AWS_ACCESS_KEY_ID", os.Getenv("AWS_ACCESS_KEY"))
	secret := env("AWS_SECRET_ACCESS_KEY", os.Getenv("AWS_SECRET_KEY"))
	region := env("AWS_REGION", "eu-west-1")
	bucket := env("AWS_S3_BUCKET", "pandora")

	// TODO enable aws debug logging
	kind := "s3"
	config := stow.ConfigMap{
		s3.ConfigEndpoint:    endpoint,
		s3.ConfigAccessKeyID: id,
		s3.ConfigSecretKey:   secret,
		s3.ConfigRegion:      region,
		s3.ConfigDisableSSL:  "true",
	}
	return &Stow{
		kind:   kind,
		config: config,
		bucket: bucket,
	}
}

// CloudStore is store of any BLOB objects
type CloudStore interface {
	Download(ctx context.Context, id string, w io.Writer) error
	DownloadFile(ctx context.Context, file *FileInfo, w io.Writer) error
	Upload(ctx context.Context, path, mediaType string, r io.ReadCloser) (map[string]interface{}, error)
	Delete(ctx context.Context, id string) (string, interface{}, error)
	DeleteObject(ctx context.Context, path string) error
}

// Stow is cloud object store implemented using https://github.com/graymeta/stow
type Stow struct {
	kind   string
	config stow.ConfigMap
	bucket string
}

// EnsureBucket creates S3 bucket
func (fs *Stow) EnsureBucket() error {
	location, err := stow.Dial(fs.kind, fs.config)
	if err != nil {
		log.Errorf("stow.Dial fail: %v", err)
		return err
	}
	defer location.Close()

	_, err = location.CreateContainer(fs.bucket)
	return err
}

// Download object at given path
func (fs *Stow) Download(ctx context.Context, id string, w io.Writer) error {
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
func (fs *Stow) DownloadFile(ctx context.Context, file *FileInfo, w io.Writer) error {
	if file == nil {
		return fmt.Errorf("file not found")
	}

	location, err := stow.Dial(fs.kind, fs.config)
	if err != nil {
		log.Errorf("stow.Dial fail: %v", err)
		return err
	}
	defer location.Close()

	url, err := fs.fileUrl(file.Path)
	if err != nil {
		return err
	}

	item, err := location.ItemByURL(url)
	if err != nil {
		log.Errorf("store.ItemByURL fail: %v", err)
		return err
	}

	r, err := item.Open()
	if err != nil {
		log.Errorf("stow.Item.Open fail: %v", err)
		return err
	}
	defer r.Close()

	_, err = io.Copy(w, r)

	return err
}

func (fs *Stow) fileUrl(path string) (*url.URL, error) {
	url, err := url.Parse("s3://" + fs.bucket + "/" + path)
	if err != nil {
		log.Errorf("url.Parse fail: %v", err)
		return nil, err
	}
	return url, nil
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
func (fs *Stow) Upload(ctx context.Context, path, mediaType string, r io.ReadCloser) (map[string]interface{}, error) {
	_, err := fs.fileUrl(path)
	if err != nil {
		return nil, err
	}

	location, err := stow.Dial(fs.kind, fs.config)
	if err != nil {
		log.Errorf("stow.Dial fail: %v", err)
		return nil, err
	}
	defer location.Close()

	container, err := location.Container(fs.bucket)
	if err != nil {
		log.Errorf("stow.GetContainer fail: %v", err)
		return nil, err
	}

	// TODO check if put requires to load file content into memory
	// TODO add metadata
	_, err = container.Put(path, r, 0, nil)
	if err != nil {
		log.Errorf("container.Put fail: %v", err)
		return nil, err
	}

	dg, close, err := dgraph.NewClient()
	if err != nil {
		return nil, err
	}
	defer close()

	tx := dg.NewTxn()
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

func noop() {}

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

	dispose1 := noop
	dispose2 := noop
	if tx == nil {
		dg, close, err := dgraph.NewClient()
		if err != nil {
			return nil, err
		}

		dispose1 = close
		tx = dg.NewTxn()
		dispose2 = func() {
			tx.Discard(ctx)
		}
	}

	defer dispose1()
	defer dispose2()

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
func (fs *Stow) Delete(ctx context.Context, id string) (string, interface{}, error) {
	dg, close, err := dgraph.NewClient()
	if err != nil {
		return "", nil, err
	}
	defer close()

	tx := dg.NewTxn()
	defer tx.Discard(ctx)

	file, err := findFile(ctx, tx, id)
	if err != nil {
		return "", nil, err
	}

	if file == nil {
		return "", nil, fmt.Errorf("file not found: %s", id)
	}

	location, err := stow.Dial(fs.kind, fs.config)
	if err != nil {
		log.Errorf("stow.Dial fail: %v", err)
		return "", nil, err
	}
	defer location.Close()

	container, err := location.Container(fs.bucket)
	if err != nil {
		log.Errorf("stow.GetContainer fail: %v", err)
		return "", nil, err
	}

	err = container.RemoveItem(file.Path)
	if err != nil {
		log.Errorf("stow.Container.RemoveItem fail: %v", err)
		return id, nil, err
	}

	resp2, err := dgraph.DeleteNode(ctx, tx, file.ID)
	if err != nil {
		return file.ID, nil, err
	}

	return file.ID, resp2, nil
}

// DeleteObject by given path
func (fs *Stow) DeleteObject(ctx context.Context, path string) error {
	location, err := stow.Dial(fs.kind, fs.config)
	if err != nil {
		log.Errorf("stow.Dial fail: %v", err)
		return err
	}

	container, err := location.Container(fs.bucket)
	if err != nil {
		log.Errorf("stow.GetContainer fail: %v", err)
		return err
	}

	err = container.RemoveItem(path)
	if err != nil {
		log.Errorf("stow.Container.RemoveItem fail: %v", err)
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
	dg, close, err := dgraph.NewClient()
	if err != nil {
		return nil, err
	}
	defer close()

	tx := dg.NewTxn()
	defer tx.Discard(ctx)

	return findFile(ctx, tx, id)
}
