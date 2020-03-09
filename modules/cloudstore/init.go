package cloudstore

import log "github.com/sirupsen/logrus"

// InitStore initializes file store
func InitStore() {
	fs := NewStow()
	err := fs.EnsureBucket()
	if err != nil {
		log.Errorf("EnsureBucket fail: %v", err)
	}
}
