package boltdb_test

import (
	"io/ioutil"
	"os"
	"testing"

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
	// make repo call put
	// expect output

	cases := []struct {
		name string
		want *racer.Message
	}{
		{name: "it puts message into the database", want: &racer.Message{Body: "test"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			tr := newRepo()
			msgs := []*racer.Message{tc.want}

			if err := tr.db.Put("ID", msgs...); err != nil {
				t.Fatalf("test faild")
			}

			got := tr.db.FetchX("ID", 1)

			if got[0].Body != tc.want.Body {
				t.Fatalf("got: %+v want: %+v", got[0], tc.want)
			}

		})

	}

}

// func TestFetchX(t *testing.T) {
// 	// create new repo for test
// 	// test that fetchx works
// }
