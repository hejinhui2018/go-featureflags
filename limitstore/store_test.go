package limitstore

import (
	"errors"
	"fmt"
	"sync"
	"testing"
)

func TestPutAndGet(t *testing.T) {
	s := New(2)
	if err := s.Put("theme", "dark"); err != nil {
		t.Fatal(err)
	}
	if got, ok := s.Get("theme"); !ok || got != "dark" {
		t.Fatalf("Get(theme) = %q, %v", got, ok)
	}
}

func TestCapacityRejectsNewKey(t *testing.T) {
	s := New(1)
	if err := s.Put("theme", "dark"); err != nil {
		t.Fatal(err)
	}
	if err := s.Put("layout", "compact"); !errors.Is(err, ErrCapacity) {
		t.Fatalf("Put(layout) error = %v", err)
	}
}

func TestDeleteFreesCapacity(t *testing.T) {
	s := New(1)
	if err := s.Put("theme", "dark"); err != nil {
		t.Fatal(err)
	}
	if !s.Delete("theme") {
		t.Fatal("Delete(theme) = false")
	}
	if err := s.Put("layout", "compact"); err != nil {
		t.Fatal(err)
	}
}

func TestRejectsBlankKey(t *testing.T) {
	s := New(1)
	if err := s.Put("  ", "value"); !errors.Is(err, ErrInvalidKey) {
		t.Fatalf("Put(blank) error = %v", err)
	}
}

// ---- regression tests for capacity-bounded update ----

// TestUpdateExistingKeyAtCapacity verifies that an existing key can be
// updated even when the store has reached its capacity limit.
func TestUpdateExistingKeyAtCapacity(t *testing.T) {
	s := New(2)

	// Fill the store to capacity.
	if err := s.Put("theme", "dark"); err != nil {
		t.Fatalf("Put(theme, dark) error = %v", err)
	}
	if err := s.Put("density", "comfortable"); err != nil {
		t.Fatalf("Put(density, comfortable) error = %v", err)
	}
	if got := s.Len(); got != 2 {
		t.Fatalf("Len() = %d, want 2", got)
	}

	// Update an existing key while at capacity -- must succeed.
	if err := s.Put("theme", "light"); err != nil {
		t.Fatalf("Put(theme, light) at capacity error = %v", err)
	}

	// The count must remain 2 (no new record created).
	if got := s.Len(); got != 2 {
		t.Fatalf("Len() after update = %d, want 2", got)
	}

	// The value must reflect the update.
	if got, ok := s.Get("theme"); !ok || got != "light" {
		t.Fatalf("Get(theme) = %q, %v; want %q, true", got, ok, "light")
	}

	// The other key must be untouched.
	if got, ok := s.Get("density"); !ok || got != "comfortable" {
		t.Fatalf("Get(density) = %q, %v; want %q, true", got, ok, "comfortable")
	}

	// Update the second existing key as well.
	if err := s.Put("density", "compact"); err != nil {
		t.Fatalf("Put(density, compact) at capacity error = %v", err)
	}
	if got := s.Len(); got != 2 {
		t.Fatalf("Len() after second update = %d, want 2", got)
	}
	if got, ok := s.Get("density"); !ok || got != "compact" {
		t.Fatalf("Get(density) = %q, %v; want %q, true", got, ok, "compact")
	}
}

// TestNewKeyStillRejectedAtCapacity verifies that writing a brand-new key
// when the store is full still returns ErrCapacity and does not create a
// record.
func TestNewKeyStillRejectedAtCapacity(t *testing.T) {
	s := New(2)
	if err := s.Put("theme", "dark"); err != nil {
		t.Fatal(err)
	}
	if err := s.Put("density", "comfortable"); err != nil {
		t.Fatal(err)
	}

	// A third, brand-new key must be rejected.
	if err := s.Put("lang", "en"); !errors.Is(err, ErrCapacity) {
		t.Fatalf("Put(lang, en) error = %v, want ErrCapacity", err)
	}

	// No record should have been created.
	if got := s.Len(); got != 2 {
		t.Fatalf("Len() = %d, want 2", got)
	}
	if _, ok := s.Get("lang"); ok {
		t.Fatal("Get(lang) should be not-found")
	}

	// Existing keys must be unchanged.
	if got, ok := s.Get("theme"); !ok || got != "dark" {
		t.Fatalf("Get(theme) = %q, %v; want %q, true", got, ok, "dark")
	}
	if got, ok := s.Get("density"); !ok || got != "comfortable" {
		t.Fatalf("Get(density) = %q, %v; want %q, true", got, ok, "comfortable")
	}
}

// TestDeleteThenAddAtCapacity verifies that after deleting a record, a new
// key can be added (capacity is genuinely freed).
func TestDeleteThenAddAtCapacity(t *testing.T) {
	s := New(2)
	_ = s.Put("theme", "dark")
	_ = s.Put("density", "comfortable")

	if !s.Delete("theme") {
		t.Fatal("Delete(theme) = false")
	}
	if got := s.Len(); got != 1 {
		t.Fatalf("Len() after delete = %d, want 1", got)
	}

	// Now there is room for one more key.
	if err := s.Put("lang", "en"); err != nil {
		t.Fatalf("Put(lang, en) error = %v", err)
	}
	if got := s.Len(); got != 2 {
		t.Fatalf("Len() after add = %d, want 2", got)
	}
	if got, ok := s.Get("lang"); !ok || got != "en" {
		t.Fatalf("Get(lang) = %q, %v", got, ok)
	}

	// Store is full again -- new key must be rejected.
	if err := s.Put("extra", "fail"); !errors.Is(err, ErrCapacity) {
		t.Fatalf("Put(extra, fail) error = %v, want ErrCapacity", err)
	}

	// But updating existing keys still works.
	if err := s.Put("lang", "fr"); err != nil {
		t.Fatalf("Put(lang, fr) error = %v", err)
	}
	if got, ok := s.Get("lang"); !ok || got != "fr" {
		t.Fatalf("Get(lang) = %q, %v; want %q, true", got, ok, "fr")
	}
}

// TestBlankKeyAtCapacity verifies that blank-key validation still fires even
// when the store is full (capacity check must not shadow the key check).
func TestBlankKeyAtCapacity(t *testing.T) {
	s := New(2)
	_ = s.Put("theme", "dark")
	_ = s.Put("density", "comfortable")

	if err := s.Put("", "x"); !errors.Is(err, ErrInvalidKey) {
		t.Fatalf("Put(empty) error = %v, want ErrInvalidKey", err)
	}
	if err := s.Put("   ", "x"); !errors.Is(err, ErrInvalidKey) {
		t.Fatalf("Put(whitespace) error = %v, want ErrInvalidKey", err)
	}
}

// TestNormalReadWrite verifies basic Put/Get/Delete round-trip after the fix.
func TestNormalReadWrite(t *testing.T) {
	s := New(3)
	_ = s.Put("a", "1")
	_ = s.Put("b", "2")
	_ = s.Put("c", "3")

	for k, want := range map[string]string{"a": "1", "b": "2", "c": "3"} {
		got, ok := s.Get(k)
		if !ok || got != want {
			t.Fatalf("Get(%s) = %q, %v; want %q, true", k, got, ok, want)
		}
	}

	// Update
	_ = s.Put("a", "99")
	if got, _ := s.Get("a"); got != "99" {
		t.Fatalf("Get(a) after update = %q, want 99", got)
	}

	// Delete
	if !s.Delete("b") {
		t.Fatal("Delete(b) = false")
	}
	if _, ok := s.Get("b"); ok {
		t.Fatal("Get(b) should be not-found after delete")
	}
}

// TestConcurrentPutUpdate exercises concurrent access: multiple goroutines
// update existing keys and attempt to add new keys simultaneously.
func TestConcurrentPutUpdate(t *testing.T) {
	s := New(100)
	// Pre-populate 50 keys.
	for i := 0; i < 50; i++ {
		_ = s.Put(fmt.Sprintf("key%d", i), "init")
	}

	var wg sync.WaitGroup
	// 100 goroutines updating existing keys.
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", idx%50)
			val := fmt.Sprintf("val%d", idx)
			_ = s.Put(key, val)
		}(i)
	}
	// 100 goroutines trying to add new keys (capacity 100, 50 used, room for 50).
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			_ = s.Put(fmt.Sprintf("new%d", idx), "new")
		}(i)
	}
	wg.Wait()

	// Len must never exceed capacity.
	if got := s.Len(); got > 100 {
		t.Fatalf("Len() = %d, exceeds capacity 100", got)
	}

	// All 50 pre-existing keys must still be present (values may differ).
	for i := 0; i < 50; i++ {
		if _, ok := s.Get(fmt.Sprintf("key%d", i)); !ok {
			t.Fatalf("key%d missing after concurrent access", i)
		}
	}
}
