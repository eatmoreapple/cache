package cache

import (
	"io"
	"sync"
	"time"

	gocache "github.com/patrickmn/go-cache"
)

const (
	// NoExpiration is used to indicate that an item should never expire.
	NoExpiration time.Duration = -1
	// DefaultExpiration is used to indicate that an item should use the
	DefaultExpiration time.Duration = 0
)

// GenericCache is a generic cache that can be used with any type.
type GenericCache[T any] struct {
	cache *gocache.Cache
}

// Set add an item to the cache, replacing any existing item. If the duration is 0
// (DefaultExpiration), the cache's default expiration time is used. If it is -1
// (NoExpiration), the item never expires.
func (g *GenericCache[T]) Set(key string, v T) {
	g.SetWithExpireIn(key, v, DefaultExpiration)
}

// SetWithExpireIn add an item to the cache, replacing any existing item. If the duration is 0
func (g *GenericCache[T]) SetWithExpireIn(key string, value T, expireIn time.Duration) {
	g.cache.Set(key, value, expireIn)
}

// Get returns the value of the item associated with the key, or nil if no item
func (g *GenericCache[T]) Get(key string) (result T, exists bool) {
	v, ok := g.cache.Get(key)
	if !ok {
		return
	}
	return v.(T), true
}

// Delete removes the provided key from the cache.
func (g *GenericCache[T]) Delete(key string) {
	g.cache.Delete(key)
}

// DeleteExpired removes all expired items from the cache.
func (g *GenericCache[T]) DeleteExpired() {
	g.cache.DeleteExpired()
}

// Add adds an item to the cache, only if the key does not already exist.
// otherwise, it returns false and does nothing.
func (g *GenericCache[T]) Add(key string, value T) bool {
	return g.AddWithExpireIn(key, value, DefaultExpiration)
}

// AddWithExpireIn adds an item to the cache, only if the key does not already exist.
// otherwise, it returns false and does nothing.
func (g *GenericCache[T]) AddWithExpireIn(key string, value T, expireIn time.Duration) bool {
	err := g.cache.Add(key, value, expireIn)
	return err == nil
}

// SetIfNotExists sets the value of the item associated with the key, only if the key does not already exist.
// otherwise, it returns an error.
func (g *GenericCache[T]) SetIfNotExists(key string, value T) bool {
	return g.SetIfNotExistsWithExpireIn(key, value, DefaultExpiration)
}

// SetIfNotExistsWithExpireIn sets the value of the item associated with the key, only if the key does not already exist.
// otherwise, it returns an error.
func (g *GenericCache[T]) SetIfNotExistsWithExpireIn(key string, value T, expireIn time.Duration) bool {
	err := g.cache.Add(key, value, expireIn)
	return err == nil
}

// Replace replaces an item in the cache, only if the key already exists.
// otherwise, does nothing and returns false.
func (g *GenericCache[T]) Replace(key string, value T) bool {
	return g.ReplaceWithExpireIn(key, value, DefaultExpiration)
}

// ReplaceWithExpireIn replaces an item in the cache, only if the key already exists.
// otherwise, does nothing and returns false.
func (g *GenericCache[T]) ReplaceWithExpireIn(key string, value T, expireIn time.Duration) bool {
	err := g.cache.Replace(key, value, expireIn)
	return err == nil
}

// SetIfExists sets the value of the item associated with the key, only if the key already exists.
// otherwise, does nothing and returns false.
func (g *GenericCache[T]) SetIfExists(key string, value T) bool {
	return g.SetIfExistsWithExpireIn(key, value, DefaultExpiration)
}

// SetIfExistsWithExpireIn sets the value of the item associated with the key, only if the key already exists.
// otherwise, does nothing and returns false.
func (g *GenericCache[T]) SetIfExistsWithExpireIn(key string, value T, expireIn time.Duration) bool {
	return g.ReplaceWithExpireIn(key, value, expireIn)
}

// Flush removes all items from the cache.
func (g *GenericCache[T]) Flush() {
	g.cache.Flush()
}

// DumpTo dumps the cache to the given writer.
func (g *GenericCache[T]) DumpTo(writer io.Writer) error {
	return g.cache.Save(writer)
}

// LoadFrom loads the cache from the given reader.
func (g *GenericCache[T]) LoadFrom(reader io.Reader) error {
	return g.cache.Load(reader)
}

// New returns a new GenericCache[T] with the given default expiration duration and cleanup interval.
func New[T any](defaultExpiration, cleanupInterval time.Duration) *GenericCache[T] {
	cache := gocache.New(defaultExpiration, cleanupInterval)
	return &GenericCache[T]{cache: cache}
}

// Numeric is a numeric type.
// it could be int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr, float32, float64.
type Numeric interface {
	// signed unsigned float are all supported.
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

// NumericCache is a cache that can be used with any numeric type.
type NumericCache[T Numeric] struct {
	*GenericCache[T]
	// mu is used to protect the cache from concurrent access.
	mu sync.RWMutex
}

// Increment increments the value of the item associated with the key by delta.
// if the key does not exist, it returns false and zero.
// otherwise, it returns true and the incremented value.
func (n *NumericCache[T]) Increment(key string, delta T) (T, bool) {
	n.mu.Lock()
	v, ok := n.Get(key)
	if !ok {
		n.mu.Unlock()
		return v, false
	}
	v += delta
	n.Set(key, v)
	n.mu.Unlock()
	return v, true
}

// Decrement decrements the value of the item associated with the key by delta.
// if the key does not exist, it returns false and zero.
// otherwise, it returns true and the decremented value.
func (n *NumericCache[T]) Decrement(key string, delta T) (T, bool) {
	return n.Increment(key, -delta)
}

// NewNumericCache returns a new NumericCache[T] with the given default expiration duration and cleanup interval.
func NewNumericCache[T Numeric](defaultExpiration, cleanupInterval time.Duration) *NumericCache[T] {
	return &NumericCache[T]{
		GenericCache: New[T](defaultExpiration, cleanupInterval),
	}
}
