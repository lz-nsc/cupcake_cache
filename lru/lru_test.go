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
		t.Fatalf("cache get \"hello\" failed, got: %s", string(val.(String)))
	}

	cache.Add("hello", String("Lynn"))
	if val, ok := cache.Get("hello"); !ok || string(val.(String)) != "Lynn" {
		t.Fatalf("cache update \"hello\" failed, got: %s", string(val.(String)))
	}
}

func TestRemoveOldest(t *testing.T) {
	key1 := "Hello"
	val1 := String("World")
	key2 := "Hi"
	val2 := String("Lynn")
	cap := len(key1) + val1.Len() + len(key2) + val2.Len()

	cache := New(int64(cap), nil)
	cache.Add(key1, val1)
	cache.Add(key2, val2)
	cache.Add("Bye", String("History"))

	if _, ok := cache.Get(key1); ok {
		t.Fatal("cache remove oldest failed")
	}

	if val, ok := cache.Get(key2); !ok || string(val.(String)) != string(val2) {
		t.Fatalf("cache remove oldest failed, get key2 got: %v", string(val.(String)))
	}

	if val, ok := cache.Get("Bye"); !ok || string(val.(String)) != "History" {
		t.Fatalf("cache update \"Bye\" failed, got: %s", string(val.(String)))
	}
}
