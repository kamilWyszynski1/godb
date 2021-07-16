package godb

// LRU provides lru cache for storing
type LRU interface {
	Add(int, node)
	Get(int) node
}

// LRUCache help tree to store cache with deleting last recently used values
type LRUCache struct {
	maxEntries int
	Cache      map[int]node

	lru int // stores last recently used value
}
