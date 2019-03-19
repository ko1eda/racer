package broker

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/speps/go-hashids"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Generator generates random hash id strings.
// These are used to uniquely indentify topics started by a broker.
type Generator struct {
	encoder *hashids.HashID
	salt    []byte
	minlen  int // the minimum length of an id
}

// NewGenerator returns a Generator with default options.
// It can take a variadic number of functional options.
func NewGenerator(opts ...func(*Generator)) (*Generator, error) {
	gen := Generator{salt: make([]byte, 16), minlen: 12}

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

// NewID returns a unique id string
func (g *Generator) NewID() (string, error) {
	id, err := g.encoder.Encode([]int{rand.Int()})

	if err != nil {
		return "", fmt.Errorf("Failed encoding ints to id %v", err)
	}

	return id, nil
}
