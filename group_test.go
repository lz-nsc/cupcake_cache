package cupcake_cache

import (
	"testing"
)

// fake database
var fakeDatabase = map[string]string{
	"Hello": "World",
	"Hi":    "Lynn",
	"Bye":   "History",
}

func TestGroup(t *testing.T) {
	// Use counter to count how many times application get data from database
	counter := map[string]int{}
	// Create a new group
	group := NewGroup("test", 0, GetterFunc(func(key string) ([]byte, error) {
		// get data from fake database
		if val, ok := fakeDatabase[key]; ok {
			counter[key]++
			return []byte(val), nil
		}
		return nil, nil
	}))

	for key, val := range fakeDatabase {
		// first time get data
		bv, err := group.Get(key)
		if err != nil {
			t.Fatalf("failed to get data with key %s, err: %v", key, err)
		}
		if bv.String() != val {
			t.Fatalf("failed to get correct data with key %s, want: %s, got %s", key, val, bv.String())
		}

		// second time get data
		bv, err = group.Get(key)
		if err != nil {
			t.Fatalf("failed to get data with key %s, err: %v", key, err)
		}
		if counter[key] != 1 {
			t.Fatalf("counter does not match for key %s [want: %d, got %d]", key, 1, counter[key])
		}
		if bv.String() != val {
			t.Fatalf("failed to get correct data with key %s [want: %s, got %s]", key, val, bv.String())
		}
	}
}
