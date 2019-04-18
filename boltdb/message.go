package boltdb

import (
	"encoding/binary"
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
	"github.com/tinylttl/racer"
)

// var _ racer.MessageRepo = (*Repo)(nil)

// MessageRepo provides an interface for interacting with a storage solution
// type MessageRepo interface {
// 	Fetch(ID string) []*racer.Message
// 	FetchX(ID string, x int) []*racer.Message
// 	Put(ID string, msgs ...*racer.Message) error
// 	Delete(ID string) error
// }

// Repo implements racer.MessageRepo
type Repo struct {
	db *bolt.DB
}

// NewRepo returns a new repository intialized with an open
// DB connection
// TODO: Make This configureable
func NewRepo(path string) (*Repo, error) {
	// /boltdb/racer.db
	db, err := bolt.Open(path, 0600, nil)

	// wrap error and return,
	// in client, type switch based on returned error type
	// if DBError shut down the system bc
	// we have a serious problem
	if err != nil {
		return nil, err
	}

	r := &Repo{db: db}

	return r, nil
}

// Close closes the repos database connection
// you can no longer use it to connect to the database after calling this.
func (r *Repo) Close() { defer r.db.Close() }

// func (r *Repo) Fetch(ID string) []*racer.Message {

// }

// FetchX fetches the latest x messages from a given bucket identified by ID
// func (r *Repo) FetchX(ID string, x int) ([]*racer.Message, error) {
// 	res := make(chan [][]byte, x)

// 	// This may possibly be bad code goroutine accessing db after view function ends
// 	go func() {
// 		err := r.db.View(func(tx *bolt.Tx) error {
// 			b := tx.Bucket([]byte(ID))
// 			c := b.Cursor()

// 			// wrap Cursors Next method to return nil after a certain thrshold limit
// 			// i is 1 here because c.First() returns first key val pair
// 			// if i was 0 based x =1 would return c.First() i = 0 bc this is bultin func, c.NextLimit() i = 0
// 			i := 1
// 			cNextLimit := func() (key []byte, value []byte) {

// 				if i == x {
// 					return nil, nil
// 				}

// 				key, value = c.Next()

// 				i++

// 				return key, value
// 			}

// 			// seek from newest date(unix.now().nano()) to oldest date to find most recent messages
// 			// for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
// 			// 	fmt.Printf("%s: %s\n", k, v)
// 			// }

// 			// format [ [[]byte(key), []byte(value)], ... ]
// 			for k, v := c.First(); k != nil; k, v = cNextLimit() {
// 				res <- [][]byte{k, v}
// 			}

// 			close(res)

// 			return nil
// 		})

// 		if err != nil {
// 			panic(err)
// 		}
// 	}()

// 	// read values from the
// 	msgs := make([]*racer.Message, 0, x)
// 	done := make(chan struct{}, 1)
// 	go func() {
// 		var msg *racer.Message
// 		for item := range res {
// 			err := json.Unmarshal(item[1], &msg)

// 			if err != nil {
// 				panic(err) // unmarshall error
// 			}

// 			msgs = append(msgs, msg)
// 		}

// 		close(done)
// 	}()

// 	<-done

// 	return msgs, nil
// }

// Put stores any number of messages to the bucket identified with ID
func (r *Repo) Put(ID string, msgs ...*racer.Message) error {
	err := r.db.Batch(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(ID))

		if err != nil {
			return errors.Wrap(err, "could not find or create bucket")
		}

		for _, msg := range msgs {
			marshalledbytes, err := json.Marshal(msg)

			if err != nil {
				return errors.Wrap(err, "could not marshall msg")
			}

			// store the timestamp converted to bytes askey, marshalled *racer.Message as data
			err = b.Put([]byte(i64tob(msg.Timestamp)), marshalledbytes)

			if err != nil {
				return errors.Wrap(err, "could not store msg to database")
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// u64tob converts a uint64 into an 8-byte slice.
func i64tob(v int64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

// convert bytes pack to timestamps
func btoi64(b []byte) int64 {
	return int64(binary.BigEndian.Uint64(b[0:8]))
}

// FetchX fetches the latest x messages
func (r *Repo) FetchX(ID string, x int) ([]*racer.Message, error) {
	var msg *racer.Message
	msgs := make([]*racer.Message, 0, x)

	err := r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ID))
		c := b.Cursor()

		// wrap Cursors Next method to return nil after a certain thrshold limit
		i := 1
		cNextLimit := func() (key []byte, value []byte) {

			if i == x {
				return nil, nil
			}

			key, value = c.Next()

			i++

			return key, value
		}

		for k, v := c.First(); k != nil; k, v = cNextLimit() {
			err := json.Unmarshal(v, &msg)

			if err != nil {
				return errors.Wrap(err, "could not marshall msg")
			}

			msgs = append(msgs, msg)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return msgs, nil
}

// func (r *Repo) Delete(ID string) error {

// }
