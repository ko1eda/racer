package main

import (
	"net/http"

	rhttp "github.com/tinylttl/racer/http"
)

func main() {
	http.ListenAndServe(":80", rhttp.NewHandler())
}
