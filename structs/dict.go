package structs

import (
	"sync"

	"github.com/VincentFF/thinredis/util"
)

const MaxTableSize = 2048
const MinTableSize = 16

// Dict manage a table slice with multiple hashmap shards.
// It is threads safe by using rwLock.
// it supports maximum table size = MaxConSize
type Dict struct {
	dict []*shard
	size int
}

type shard struct {
	mp   map[string]any
	rwMu *sync.RWMutex
}

func NewDict() *Dict {
	size := ShardNum
	if size < MinTableSize || size > MaxTableSize {
		size = 128
	}
	dict := &Dict{
		dict: make([]*shard, size),
		size: size,
	}
	for i := 0; i < size; i++ {
		dict.dict[i] = &shard{mp: make(map[string]any), rwMu: &sync.RWMutex{}}
	}
	return dict
}

func (m *Dict) getKeyPos(key string) int {
	return util.HashKey(key) % m.size
}

// Set key-value pair to the map.
// Return 1 if added, 0 if updated.
func (m *Dict) Set(key string, rob any) int {
	added := 0
	pos := m.getKeyPos(key)

	shard := m.dict[pos]
	shard.rwMu.Lock()
	defer shard.rwMu.Unlock()

	if _, ok := shard.mp[key]; !ok {
		added = 1
	}
	shard.mp[key] = rob
	return added
}

func (m *Dict) SetIfExist(key string, rob any) int {
	pos := m.getKeyPos(key)
	shard := m.dict[pos]
	shard.rwMu.Lock()
	defer shard.rwMu.Unlock()

	if _, ok := shard.mp[key]; ok {
		shard.mp[key] = rob
		return 1
	}
	return 0
}

func (m *Dict) SetIfNotExist(key string, rob any) int {
	pos := m.getKeyPos(key)
	shard := m.dict[pos]
	shard.rwMu.Lock()
	defer shard.rwMu.Unlock()

	if _, ok := shard.mp[key]; !ok {
		shard.mp[key] = rob
		return 1
	}
	return 0
}

func (m *Dict) Get(key string) (any, bool) {
	pos := m.getKeyPos(key)
	shard := m.dict[pos]
	shard.rwMu.RLock()
	defer shard.rwMu.RUnlock()

	value, ok := shard.mp[key]
	return value, ok
}

func (m *Dict) Delete(key string) int {
	pos := m.getKeyPos(key)
	shard := m.dict[pos]
	shard.rwMu.Lock()
	defer shard.rwMu.Unlock()

	if _, ok := shard.mp[key]; ok {
		delete(shard.mp, key)
		return 1
	}
	return 0
}

func (m *Dict) Len() int {
	sum := 0
	for _, shard := range m.dict {
		shard.rwMu.RLock()
		sum += len(shard.mp)
		shard.rwMu.RUnlock()
	}
	return sum
}

func (m *Dict) Clear() {
	*m = *NewDict()
}

func (m *Dict) Keys() []string {
	keys := make([]string, 0)
	for _, shard := range m.dict {
		shard.rwMu.RLock()
		for key := range shard.mp {
			keys = append(keys, key)
		}
		shard.rwMu.RUnlock()
	}
	return keys
}
