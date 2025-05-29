package cache

import (
	"context"
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type item struct {
	Value      interface{}
	Expiration int64
}

// FileStore implements file-based caching
type FileStore struct {
	dir     string
	prefix  string
	items   map[string]item
	mu      sync.RWMutex
	janitor *janitor
}

// NewFileStore creates a new file cache store
func NewFileStore(config *Config) (*FileStore, error) {
	dir := config.FilePath
	if dir == "" {
		dir = "storage/cache"
	}

	// Create cache directory if it doesn't exist
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	store := &FileStore{
		dir:    dir,
		prefix: config.Prefix,
		items:  make(map[string]item),
	}

	// Start janitor to clean expired items
	store.janitor = newJanitor(store)
	store.janitor.start()

	// Load existing cache files
	if err := store.load(); err != nil {
		return nil, err
	}

	return store, nil
}

// Get retrieves a cached value
func (s *FileStore) Get(ctx context.Context, key string) (interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.get(key)
}

// Set stores a value in cache
func (s *FileStore) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.set(key, value, expiration)
}

// Delete removes a value from cache
func (s *FileStore) Delete(ctx context.Context, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.delete(key)
}

// Clear removes all values from cache
func (s *FileStore) Clear(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.items = make(map[string]item)
	return s.save()
}

// Has checks if a key exists and is not expired
func (s *FileStore) Has(ctx context.Context, key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, err := s.get(key)
	return err == nil
}

// Remember gets from cache or stores the result of getter
func (s *FileStore) Remember(ctx context.Context, key string, expiration time.Duration, getter func() (interface{}, error)) (interface{}, error) {
	// Try to get from cache first
	if value, err := s.Get(ctx, key); err == nil {
		return value, nil
	}

	// Get fresh value
	value, err := getter()
	if err != nil {
		return nil, err
	}

	// Store in cache
	if err := s.Set(ctx, key, value, expiration); err != nil {
		return nil, err
	}

	return value, nil
}

// Increment increments a number value
func (s *FileStore) Increment(ctx context.Context, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	value, err := s.get(key)
	if err != nil {
		return err
	}

	switch v := value.(type) {
	case int:
		return s.set(key, v+1, NoExpiration)
	case int64:
		return s.set(key, v+1, NoExpiration)
	case float64:
		return s.set(key, v+1, NoExpiration)
	default:
		return fmt.Errorf("value is not a number")
	}
}

// Decrement decrements a number value
func (s *FileStore) Decrement(ctx context.Context, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	value, err := s.get(key)
	if err != nil {
		return err
	}

	switch v := value.(type) {
	case int:
		return s.set(key, v-1, NoExpiration)
	case int64:
		return s.set(key, v-1, NoExpiration)
	case float64:
		return s.set(key, v-1, NoExpiration)
	default:
		return fmt.Errorf("value is not a number")
	}
}

// Close stops the janitor and saves cache to disk
func (s *FileStore) Close() error {
	s.janitor.stopJanitor()
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.save()
}

// Internal methods

func (s *FileStore) get(key string) (interface{}, error) {
	item, found := s.items[key]
	if !found {
		return nil, fmt.Errorf("key not found: %s", key)
	}

	if item.Expiration > 0 && item.Expiration < time.Now().UnixNano() {
		s.delete(key)
		return nil, fmt.Errorf("key expired: %s", key)
	}

	return item.Value, nil
}

func (s *FileStore) set(key string, value interface{}, expiration time.Duration) error {
	var exp int64
	if expiration > 0 {
		exp = time.Now().Add(expiration).UnixNano()
	}

	s.items[key] = item{
		Value:      value,
		Expiration: exp,
	}

	return s.save()
}

func (s *FileStore) delete(key string) error {
	delete(s.items, key)
	return s.save()
}

func (s *FileStore) save() error {
	file, err := os.Create(s.getCacheFile())
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	return encoder.Encode(s.items)
}

func (s *FileStore) load() error {
	file, err := os.Open(s.getCacheFile())
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	return decoder.Decode(&s.items)
}

func (s *FileStore) getCacheFile() string {
	return filepath.Join(s.dir, "cache.gob")
}

// Janitor for cleaning expired items

type janitor struct {
	store    *FileStore
	interval time.Duration
	stopChan chan bool
}

func newJanitor(store *FileStore) *janitor {
	return &janitor{
		store:    store,
		interval: time.Hour,
		stopChan: make(chan bool),
	}
}

func (j *janitor) start() {
	go func() {
		ticker := time.NewTicker(j.interval)
		for {
			select {
			case <-ticker.C:
				j.clean()
			case <-j.stopChan:
				ticker.Stop()
				return
			}
		}
	}()
}

func (j *janitor) stopJanitor() {
	j.stopChan <- true
}

func (j *janitor) clean() {
	j.store.mu.Lock()
	defer j.store.mu.Unlock()

	now := time.Now().UnixNano()
	for key, item := range j.store.items {
		if item.Expiration > 0 && item.Expiration < now {
			delete(j.store.items, key)
		}
	}
	j.store.save()
}
