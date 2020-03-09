package cloudstore

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"

	"github.com/graymeta/stow"
	_ "github.com/graymeta/stow/google"
	"github.com/graymeta/stow/s3"
	"github.com/sergeyt/pandora/modules/dgraph"
	log "github.com/sirupsen/logrus"
)

// NewCloudStore creates new instance of cloud store
func NewCloudStore() CloudStore {
	return NewStow()
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
	file, err := FindFile(ctx, id)
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
	// TODO what about google storage URLs
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

	file, err := FindFileImpl(ctx, tx, path)
	if err != nil {
		return nil, err
	}

	if file == nil {
		file = &FileInfo{}
	}
	file.Path = path
	file.MediaType = mediaType

	return AddFile(ctx, tx, file)
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

	file, err := FindFileImpl(ctx, tx, id)
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
