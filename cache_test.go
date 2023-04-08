package cache

import (
	"testing"
	"time"
)

// GenericCache is a generic cache that can be used with any type.
func TestGenericCache(t *testing.T) {
	c := New[string](time.Second*3, time.Second)
	c.Set("foo", "bar")
	if v, ok := c.Get("foo"); !ok || v != "bar" {
		t.Errorf("expected foo to be bar, got %v", v)
	}
	time.Sleep(time.Second * 3)
	if _, ok := c.Get("foo"); ok {
		t.Errorf("expected foo to be expired")
	}

	c.Set("foo", "bar")
	if exists := c.SetIfExists("foo", "baz"); !exists {
		t.Errorf("expected foo to be set")
	}
	if exists := c.SetIfExists("bar", "baz"); exists {
		t.Errorf("expected bar to not be set")
	}
}

func TestNewNumericCache(t *testing.T) {
	c := NewNumericCache[int64](time.Second*3, time.Second)
	c.Set("foo", 123)
	if v, ok := c.Get("foo"); !ok || v != 123 {
		t.Errorf("expected foo to be 123, got %v", v)
	}
	c.Increment("foo", 1)
	if v, ok := c.Get("foo"); !ok || v != 124 {
		t.Errorf("expected foo to be 124, got %v", v)
	}
	time.Sleep(time.Second * 3)
	if _, ok := c.Get("foo"); ok {
		t.Errorf("expected foo to be expired")
	}
}
