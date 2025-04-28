package cache

import (
	"testing"
)

type testStruct struct {
	ID      int
	Name    string
	Address []string
}

func TestCache(t *testing.T) {
	c := NewCache[int, int]()
	c.Set(1, 1)
	if c.Len() != 1 {
		t.Errorf("cache len should be 1, but got %d", c.Len())
	}
	if _, ok := c.Get(2); ok {
		t.Errorf("cache should not have key 2")
	}
}

func TestUpdateCache(t *testing.T) {
	c := NewCache[string, testStruct]()
	c.Set("a", testStruct{ID: 1, Name: "a"})
	c.Update("a", func(v *testStruct) error {
		v.Name = "b"
		return nil
	})
	if v, ok := c.Get("a"); !ok || v.Name != "b" {
		t.Errorf("cache should have key a with value b, but got %v", v)
	}

	c.Update("c", func(v *testStruct) error {
		v.Name = "d"
		return nil
	})
	if v, ok := c.Get("c"); !ok || v.Name != "d" {
		t.Errorf("cache should have key c with value d, but got %v", v)
	}
}

func TestConcurrentCache(t *testing.T) {
	c := NewCache[int, int]()

	write := func() {
		for i := 0; i < 1000; i++ {
			c.Set(i, i)
		}
	}
	read := func() {
		for i := 0; i < 1000; i++ {
			if _, ok := c.Get(i); !ok {
				t.Errorf("cache should have key %d", i)
			}
		}
	}

	go write()
	go read()
}
