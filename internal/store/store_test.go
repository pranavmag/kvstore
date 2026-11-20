package store

import (
	"fmt"
	"sync"
	"testing"
)

func TestStore(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		wantValue string
		wantOK    bool
	}{
		// Get operation tests
		{"existing key", "greeting", "hello", true},
		{"empty value", "empty", "", true},
		{"missing key", "missing", "", false},

		// Put operation overwrite test
		{"overwrite key", "overwrite", "second", true},

		// Delete operation test
		{"delete key", "delete", "", false},
	}

	s := New()
	s.Put("greeting", "hello")
	s.Put("empty", "")

	s.Put("overwrite", "first")
	s.Put("overwrite", "second")

	s.Put("delete", "val")
	s.Delete("delete")

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotValue, gotOK := s.Get(tc.key)

			if gotValue != tc.wantValue || gotOK != tc.wantOK {
				t.Errorf("Get(%q) = (%q, %t), want (%q, %t)",
					tc.key, gotValue, gotOK, tc.wantValue, tc.wantOK)
			}
		})
	}
}

func TestPutConcurrent(t *testing.T) {
	s := New()
	var wg sync.WaitGroup

	for i := range 1000 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			s.Put(fmt.Sprintf("key%d", i), "val")
		}(i)
	}

	wg.Wait()

	for i := range 1000 {
		v, ok := s.Get(fmt.Sprintf("key%d", i))
		if !ok || v != "val" {
			t.Errorf("key%d: got (%q, %t), want (\"val\", true)", i, v, ok)
		}
	}
}

func TestConcurrentReadWrite(t *testing.T) {
	s := New()

	allowed := map[string]bool{}
	for i := range 1000 {
		allowed[fmt.Sprintf("val%d", i)] = true
	}

	var wg sync.WaitGroup

	for i := range 1000 {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()
			s.Put("shared", fmt.Sprintf("val%d", i))

			v, ok := s.Get("shared")
			if !ok || !allowed[v] {
				t.Errorf("Get(%q) = (%q, %t), want a value written to the key", "shared", v, ok)
			}
		}(i)
	}

	wg.Wait()
}
