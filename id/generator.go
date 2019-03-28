package id

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/speps/go-hashids"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const (
	defaultSaltSize = 16 // default saltSize, corresponds to 16 bytes or 128 bits
	defaultMinLen   = 14 // default hashid length, note that ids are truncated by this package to only reach the minimum length
)

// Generator generates random hash id strings.
// These are used to uniquely indentify topics started by a broker.
type Generator struct {
	encoder *hashids.HashID
	salt    []byte
	minlen  int
}

// NewGenerator returns a Generator with default options.
// It can take a variadic number of functional options.
func NewGenerator(opts ...func(*Generator)) (*Generator, error) {
	gen := Generator{salt: make([]byte, defaultSaltSize), minlen: defaultMinLen}

	if _, err := rand.Read(gen.salt); err != nil {
		return nil, fmt.Errorf("Failed reading random bytes %v", err)
	}

	for _, opt := range opts {
		opt(&gen)
	}

	data := hashids.NewData()
	data.Salt = string(gen.salt)
	data.MinLength = gen.minlen

	h, err := hashids.NewWithData(data)
	if err != nil {
		return nil, err
	}

	gen.encoder = h

	return &gen, nil
}

// WithSalt passes a custom salt string to the generator. Use with NewGenerator()
func WithSalt(s []byte) func(*Generator) {
	return func(g *Generator) {
		g.salt = s
	}
}

// WithLength changes the minimum length of the hashID. Use with NewGenerator()
func WithLength(len int) func(*Generator) {
	return func(g *Generator) {
		g.minlen = len
	}
}

// NewID returns a unique id string using two random integers
// TODO: http://blog.sgmansfield.com/2016/01/the-hidden-dangers-of-default-rand/
// read this blog post and access the benefits of using a new rand.Source()
// each time the generator is created, instead of using the default source that has a lock
// See Also: https://blog.gopheracademy.com/advent-2017/a-tale-of-two-rands/
func (g *Generator) NewID() (string, error) {
	id, err := g.encoder.Encode([]int{rand.Int()})

	if err != nil {
		return "", fmt.Errorf("Failed encoding ints to ID %v", err)
	}

	return id[:g.minlen], nil
}

// // ParseID takes a hashID created by the generator and returns
// // the slice of n integers used to create it.
// func (g *Generator) ParseID(ID string) ([]int, error) {
// 	nums, err := g.encoder.DecodeWithError(ID)

// 	if err != nil {
// 		return nil, fmt.Errorf("Failed decoding ID %v", err)
// 	}

// 	return nums, nil
// }
