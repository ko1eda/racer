package id_test

import (
	"sync"
	"testing"

	"github.com/tinylttl/racer/id"
)

// Test that ID's are always unique
func TestGenerator_NewID(t *testing.T) {
	t.Run("It generates IDs of the proper length", func(t *testing.T) {
		want := 14

		gen, _ := id.NewGenerator(id.WithLength(want))
		id, err := gen.NewID()
		if err != nil {
			t.Fatalf("%v", err)
		}

		if got := len(id); got != want {
			t.Fatalf("got: %d want: %d", got, want)
		}
	})

	t.Run("It generates only unique ids", func(t *testing.T) {
		want := 10000
		res := make(chan string, 10000)

		var wg sync.WaitGroup
		wg.Add(10000)
		go func() {
			wg.Wait()
			close(res)
		}()

		// 100 goroutines each generating 100 id's for a total of 10000
		gen, _ := id.NewGenerator()
		var i int
		for i < 100 {
			go func() {
				var n int
				for n < 100 {
					id, err := gen.NewID()

					if err != nil {
						t.Fatalf("%v", err)
					}

					res <- id
					n++
					wg.Done()
				}
			}()
			i++
		}

		m := map[string]struct{}{}
		for id := range res {
			m[id] = struct{}{}
		}

		if got := len(m); got != want {
			t.Fatalf("got: %d want: %d", got, want)
		}
	})

}
