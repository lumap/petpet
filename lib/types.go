package lib

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
)

// Most of, if not all, of these types and constants are based on Discord's API documentation. However, I have stripped them down
// to what my bot needs. This is, and will, not be complete.

const (
	DISCORD_API_URL   = "https://discord.com/api/v10"
	DISCORD_CDN_URL   = "https://cdn.discordapp.com"
	DISCORD_EPOCH     = 1420070400000 // Discord epoch in milliseconds
	CONTENT_TYPE_JSON = "application/json"
	ROOT_PLACEHOLDER  = "-"
)

type Snowflake uint64
func (snowflake *Snowflake) String() string {
	return fmt.Sprintf("%d", uint64(*snowflake))
}
func (s Snowflake) MarshalJSON() ([]byte, error) {
	b := strconv.FormatUint(uint64(s), 10)
	return json.Marshal(b)
}

func (s *Snowflake) UnmarshalJSON(b []byte) error {
	str, err := strconv.Unquote(string(b))
	if err != nil {
		return err
	}

	i, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return err
	}

	*s = Snowflake(i)
	return nil
}

type BitSet uint64

type SharedMap[K comparable, V any] struct {
	mu    sync.RWMutex
	cache map[K]V
}
func NewSharedMap[K comparable, V any]() *SharedMap[K, V] {
	return &SharedMap[K, V]{
		mu:    sync.RWMutex{},
		cache: make(map[K]V),
	}
}
func (sm *SharedMap[K, V]) Has(key K) bool {
	sm.mu.RLock()
	_, available := sm.cache[key]
	sm.mu.RUnlock()
	return available
}
func (sm *SharedMap[K, V]) Set(key K, value V) {
	sm.mu.Lock()
	sm.cache[key] = value
	sm.mu.Unlock()
}
func (sm *SharedMap[K, V]) Get(key K) (V, bool) {
	sm.mu.RLock()
	item, available := sm.cache[key]
	sm.mu.RUnlock()
	return item, available
}