package lru

import "testing"

type String string

func (str String) Len() int {
	return len(str)
}

func TestGet(t *testing.T) {
	cache := New(int64(0), nil)

	cache.Add("hello", String("World"))
	if val, ok := cache.Get("hello"); !ok || string(val.(String)) != "World" {
		t.Fatalf("cache get \"hello\" failed, got: %v,%v", ok, string(val.(String)))
	}

	cache.Add("hello", String("Lynn"))
	if val, ok := cache.Get("hello"); !ok || string(val.(String)) != "Lynn" {
		t.Fatal("cache update \"hello\" failed")
	}
}
