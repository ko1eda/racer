package boltdb_test

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/tinylttl/racer"

	"github.com/tinylttl/racer/boltdb"
)

type testrepo struct {
	repo *boltdb.MessageRepo
	db   *boltdb.DB
	path string
}

func newRepo() *testrepo {
	// Retrieve a temporary path.
	f, err := ioutil.TempFile("", "")

	if err != nil {
		panic(err)
	}

	path := f.Name()

	db := boltdb.NewDB(boltdb.WithPath(path))

	if err = db.Open(); err != nil {
		panic(err)
	}
	// Open the database.
	repo := boltdb.NewMessageRepo(db)

	// Return wrapped type.
	return &testrepo{
		repo: repo,
		path: path,
		db:   db,
	}
}

// Close and delete Bolt database.
func (repo *testrepo) close() {
	defer os.Remove(repo.path)
	repo.db.Close()
}

func TestPut(t *testing.T) {
	t.Run("it puts a message into the database", func(t *testing.T) {
		tr := newRepo()

		defer tr.close()

		want := &racer.Message{Body: "test"}

		msgs := []*racer.Message{want}

		if err := tr.repo.Put("ID", msgs...); err != nil {
			t.Fatalf("test faild")
		}

		got, err := tr.repo.FetchX("ID", 1)

		if err != nil {
			t.Fatal(err)
		}

		if got[0].Body != want.Body {
			t.Fatalf("got: %+v want: %+v", got[0], want)
		}
	})

	t.Run("it puts multiple messages into the database", func(t *testing.T) {
		tr := newRepo()

		defer tr.close()

		want := 2

		msgs := []*racer.Message{
			&racer.Message{Timestamp: time.Now().UnixNano(), Body: "1"},
			&racer.Message{Timestamp: time.Now().Add(time.Hour * 22).UnixNano(), Body: "2"},
			&racer.Message{Timestamp: time.Now().Add(time.Hour * 24).UnixNano(), Body: "3"},
		}

		if err := tr.repo.Put("ID", msgs...); err != nil {
			t.Fatalf("test faild")
		}

		got, _ := tr.repo.FetchX("ID", 2)

		if len(got) != want {
			t.Fatalf("got: %+v want: %+v", len(got), want)
		}

	})
}

func TestFetchX(t *testing.T) {
	t.Run("it retrieves messages in reverse chronological order", func(t *testing.T) {
		tr := newRepo()

		defer tr.close()

		msgs := []*racer.Message{
			&racer.Message{Timestamp: time.Now().Add(time.Hour * 48).UnixNano(), Body: "1"},
			&racer.Message{Timestamp: time.Now().Add(time.Hour * 24).UnixNano(), Body: "2"},
			&racer.Message{Timestamp: time.Now().UnixNano(), Body: "3"},
		}

		tr.repo.Put("ID", msgs...)

		want := msgs[0]
		got, _ := tr.repo.FetchX("ID", 3)

		if got[0].Body != want.Body {
			t.Fatalf("got: %+v want: %+v", got[0], want)
		}
	})
}
