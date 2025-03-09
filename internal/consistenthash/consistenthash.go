package consistenthash

import (
	"fmt"
	"slices"
	"zencache/internal/config"
)

// Hash 是一个哈希函数类型，接受 []byte 并返回 uint32 哈希值。
type Hash func([]byte) uint32

// Map 表示一致性哈希地图。
type Map struct {
	hash    Hash           // 哈希函数
	replica int            // 副本数量
	keys    []int          // 存储所有哈希键，用于二分查找
	hashMap map[int]string // 哈希键到原始键的映射
	conf    *config.Config // 配置信息
}

// New 创建一个新的 Map 实例。
func New(conf *config.Config, hash Hash) *Map {
	return &Map{
		hash:    hash,
		replica: conf.Hash.Replicas,
		keys:    make([]int, 0),
		hashMap: make(map[int]string),
		conf:    conf,
	}
}

// Add 将一组键添加到一致性哈希中。
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for idx := range m.replica {
			replicaKey := fmt.Sprint(idx) + key
			hashKey := m.hash([]byte(replicaKey))
			m.hashMap[int(hashKey)] = key
			m.keys = append(m.keys, int(hashKey))
		}
	}
	slices.Sort(m.keys)
}

// Get 根据给定的键，返回一致性哈希中最接近的键。
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}
	hashKey := m.hash([]byte(key))
	idx, _ := slices.BinarySearch(m.keys, int(hashKey))
	return m.hashMap[m.keys[idx%len(m.keys)]]
}


func (m *Map) Delete(key string) {
	for idx := range m.replica {
		replicaKey := fmt.Sprint(idx) + key
		hashKey := m.hash([]byte(replicaKey))
		delete(m.hashMap, int(hashKey))
		m.keys = slices.DeleteFunc(m.keys, func(v int) bool { return v == int(hashKey) })
	}
}
