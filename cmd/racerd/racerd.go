package main

import (
	"net/http"

	"github.com/tinylttl/racer/boltdb"
	rhttp "github.com/tinylttl/racer/http"
)

func main() {
	db := boltdb.NewDB()

	if err := db.Open(); err != nil {
		panic(err)
	}
	
	repo := boltdb.NewMessageRepo(db)

	// TODO: switch to actual Server struct because these default settings are bad
	http.ListenAndServe(":80", rhttp.NewHandler(repo))
}
