package main

import (
	"container/heap"
	"encoding/json"
	"fmt"
	"os"
)

// CacheEntry represents one cached item on an OCA server.
type CacheEntry struct {
	TitleID     string
	VariantKey  string
	SizeBytes   int64
	AccessCount int64
	index       int // position in heap
}

// popularityHeap implements heap.Interface for min-popularity eviction.
type popularityHeap []*CacheEntry

func (h popularityHeap) Len() int            { return len(h) }
func (h popularityHeap) Less(i, j int) bool  { return h[i].AccessCount < h[j].AccessCount }
func (h popularityHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}
func (h *popularityHeap) Push(x any) {
	e := x.(*CacheEntry)
	e.index = len(*h)
	*h = append(*h, e)
}
func (h *popularityHeap) Pop() any {
	old := *h
	n := len(old)
	e := old[n-1]
	old[n-1] = nil
	*h = old[:n-1]
	return e
}

// OCACache simulates a popularity-weighted LFU cache for one OCA server.
type OCACache struct {
	capacityBytes int64
	usedBytes     int64
	entries       map[string]*CacheEntry
	h             popularityHeap
}

// NewOCACache creates a cache with the given byte capacity.
func NewOCACache(capacityBytes int64) *OCACache {
	c := &OCACache{
		capacityBytes: capacityBytes,
		entries:       make(map[string]*CacheEntry),
	}
	heap.Init(&c.h)
	return c
}

func entryKey(titleID, variantKey string) string {
	return titleID + "/" + variantKey
}

// Put adds or updates an entry, evicting least-popular entries as needed.
func (c *OCACache) Put(titleID, variantKey string, sizeBytes int64) bool {
	k := entryKey(titleID, variantKey)
	if _, exists := c.entries[k]; exists {
		c.entries[k].AccessCount++
		heap.Fix(&c.h, c.entries[k].index)
		return true
	}
	// Evict until there is room.
	for c.usedBytes+sizeBytes > c.capacityBytes && len(c.h) > 0 {
		evicted := heap.Pop(&c.h).(*CacheEntry)
		ek := entryKey(evicted.TitleID, evicted.VariantKey)
		delete(c.entries, ek)
		c.usedBytes -= evicted.SizeBytes
	}
	if c.usedBytes+sizeBytes > c.capacityBytes {
		return false // not enough space even after eviction
	}
	e := &CacheEntry{TitleID: titleID, VariantKey: variantKey, SizeBytes: sizeBytes, AccessCount: 1}
	heap.Push(&c.h, e)
	c.entries[k] = e
	c.usedBytes += sizeBytes
	return true
}

// Get records an access hit.
func (c *OCACache) Get(titleID, variantKey string) bool {
	k := entryKey(titleID, variantKey)
	e, ok := c.entries[k]
	if !ok {
		return false
	}
	e.AccessCount++
	heap.Fix(&c.h, e.index)
	return true
}

// Stats returns cache utilization info.
func (c *OCACache) Stats() map[string]any {
	return map[string]any{
		"capacity_bytes": c.capacityBytes,
		"used_bytes":     c.usedBytes,
		"entry_count":    len(c.entries),
		"fill_percent":   int(c.usedBytes * 100 / c.capacityBytes),
	}
}

func main() {
	// 1 GB cache, add a few entries
	cache := NewOCACache(1024 * 1024 * 1024)
	cache.Put("tt-001", "h264/1080p/5800", 512*1024*1024)
	cache.Put("tt-002", "h264/720p/2350", 256*1024*1024)
	cache.Get("tt-001", "h264/1080p/5800")
	cache.Get("tt-001", "h264/1080p/5800")
	// This should evict tt-002 (lower access count) to fit tt-003
	cache.Put("tt-003", "h264/1080p/5800", 400*1024*1024)

	stats := cache.Stats()
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(stats); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
