package main

import (
	"net/http"
	"path/filepath"

	bolt "go.etcd.io/bbolt"
)

type SVGHandler struct {
	maps storage
}

func (s *SVGHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	svgFile := filepath.Base(r.URL.Path)
	content, err := s.maps.get(svgFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if content == nil {
		http.Error(w, "image not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "image/svg+xml")
	w.Write(content)
}

type storage interface {
	save(key string, value []byte) error
	get(key string) ([]byte, error)
}

type memoryStorage map[string][]byte

func (m memoryStorage) save(key string, value []byte) error {
	m[key] = value
	return nil
}

func (m memoryStorage) get(key string) ([]byte, error) {
	if val, ok := m[key]; ok {
		return val, nil
	}
	return nil, nil
}

type persistentStorage struct {
	db *bolt.DB
}

func (persistentstorage *persistentStorage) save(key string, value []byte) error {
	return persistentstorage.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("SVG"))
		err := b.Put([]byte(key), value)
		return err
	})
}
func (persistentstorage *persistentStorage) get(key string) ([]byte, error) {
	var v []byte
	err := persistentstorage.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("SVG"))
		v = b.Get([]byte(key))
		return nil
	})
	return v, err
}

func newPersistentStorage(path string) (*persistentStorage, error) {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("SVG"))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &persistentStorage{
		db: db,
	}, nil

}
