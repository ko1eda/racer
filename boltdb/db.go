package boltdb

import (
	"os"
	"path/filepath"

	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
)

// DB wraps a database with some extra utility information
type DB struct {
	*bolt.DB
	path string
}

// NewDB returns a DB struct intialized with the default path
func NewDB(opts ...func(*DB)) *DB {
	db := &DB{path: defaultPath()}

	for _, opt := range opts {
		opt(db)
	}

	return db
}

// WithPath can be used as an argument to NewDB,
// it sets the path that the repo will use to open the database
func WithPath(path string) func(*DB) {
	return func(db *DB) {
		db.path = path
	}
}

// Open a database connection
func (db *DB) Open() error {
	// make the necessary parent directories if they do not exist
	if err := os.MkdirAll(filepath.Dir(db.path), 0700); err != nil {
		return err
	}

	d, err := bolt.Open(db.path, 0600, nil)

	if err != nil {
		return errors.Wrap(err, "DB open error")
	}

	db.DB = d

	return nil
}

// defaultPath sets the default path to a racer directory
// in the users home directory
func defaultPath() string {
	home, err := os.UserHomeDir()

	if err != nil {
		panic(err)
	}

	return filepath.FromSlash(filepath.Join(home, "racer", "racer.db"))
}

// Close closes the database connection
func (db *DB) Close() { defer db.DB.Close() }
