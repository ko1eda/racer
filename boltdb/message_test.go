package boltdb_test

import (
	// "fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/tinylttl/racer"

	"github.com/tinylttl/racer/boltdb"
)

type testrepo struct {
	db   *boltdb.Repo
	path string
}

func newRepo() *testrepo {
	// Retrieve a temporary path.
	f, err := ioutil.TempFile("", "")

	if err != nil {
		panic(err)
	}
	path := f.Name()

	// f.Close()
	// os.Remove(path)

	// Open the database.
	db, err := boltdb.NewRepo(path)

	if err != nil {
		panic(err)
	}
	// Return wrapped type.
	return &testrepo{db: db, path: path}
}

// Close and delete Bolt database.
func (repo *testrepo) Close() {
	defer os.Remove(repo.path)
	repo.db.Close()
}

func TestPut(t *testing.T) {
	t.Run("it puts a message into the database", func(t *testing.T) {
		tr := newRepo()

		defer tr.Close()

		want := &racer.Message{Body: "test"}

		msgs := []*racer.Message{want}

		if err := tr.db.Put("ID", msgs...); err != nil {
			t.Fatalf("test faild")
		}

		got, err := tr.db.FetchX("ID", 1)

		if err != nil {
			t.Fatal(err)
		}

		if got[0].Body != want.Body {
			t.Fatalf("got: %+v want: %+v", got[0], want)
		}
	})

	t.Run("it puts multiple messages into the database", func(t *testing.T) {
		tr := newRepo()

		defer tr.Close()

		want := 2

		msgs := []*racer.Message{
			&racer.Message{Timestamp: time.Now().UnixNano(), Body: "1"},
			&racer.Message{Timestamp: time.Now().Add(time.Hour * 22).UnixNano(), Body: "2"},
			&racer.Message{Timestamp: time.Now().Add(time.Hour * 24).UnixNano(), Body: "3"},
		}

		if err := tr.db.Put("ID", msgs...); err != nil {
			t.Fatalf("test faild")
		}

		got, _ := tr.db.FetchX("ID", 2)

		if len(got) != want {
			t.Fatalf("got: %+v want: %+v", len(got), want)
		}

	})
}

func TestFetchX(t *testing.T) {
	t.Run("it retrieves messages in reverse chronological order", func(t *testing.T) {
		tr := newRepo()

		defer tr.Close()

		msgs := []*racer.Message{
			&racer.Message{Timestamp: time.Now().Add(time.Hour * 48).UnixNano(), Body: "1"},
			&racer.Message{Timestamp: time.Now().Add(time.Hour * 24).UnixNano(), Body: "2"},
			&racer.Message{Timestamp: time.Now().UnixNano(), Body: "3"},
		}

		tr.db.Put("ID", msgs...)

		want := msgs[0]
		got, _ := tr.db.FetchX("ID", 3)

		if got[0].Body != want.Body {
			t.Fatalf("got: %+v want: %+v", got[0], want)
		}
	})
}

// func TestFetchX(t *testing.T) {
// 	// create new repo for test
// 	// test that fetchx works
// }
