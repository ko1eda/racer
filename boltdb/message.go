package badger

import (
	"github.com/boltdb/bolt"

	"github.com/tinylttl/racer"
)

// MessageRepo provides an interface for interacting with a storage solution
type MessageRepo interface {
	Fetch(ID string) []*racer.Message
	FetchX(ID string, x int) []*racer.Message
	Put(ID string, msgs ...*racer.Message) error
	Delete(ID string) error
}

// Repo implements racer.MessageRepo
type Repo struct {
	db *bolt.DB
}

// func (r *Repo) Fetch(ID string) []*racer.Message {

// }

// FetchX fetches the latest x messages
// func (r *Repo) FetchX(ID string, x int) []*racer.Message {
// 	// use prefix to make a "bucket"
// 	// store messages in sequential order
// }

// func (r *Repo) Put(ID string, msgs ...*racer.Message) error {

// }

// func (r *Repo) Delete(ID string) error {

// }
