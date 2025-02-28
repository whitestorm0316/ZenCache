package lru

import (
	"testing"
)

type testValue struct {
	size int
}

func (v testValue) Len() int { return v.size }

func TestCache_AddAndGet(t *testing.T) {
	c := New(100, nil)
	key := "key1"
	value := testValue{size: 10}

	c.Add(key, value)
	if c.Len() != 1 {
		t.Fatalf("expected cache size 1, got %d", c.Len())
	}

	retrieved, ok := c.Get(key)
	if !ok {
		t.Fatal("failed to retrieve added item")
	}
	if retrieved != value {
		t.Errorf("expected %v, got %v", value, retrieved)
	}
}

func TestCache_EvictOldest(t *testing.T) {
	// 每个条目大小: len(key)=2, value.Len=5 → 总大小7
	// 设置缓存容量为14（允许两个条目）
	c := New(14, nil)

	c.Add("k1", testValue{5}) // 大小7
	c.Add("k2", testValue{5}) // 大小7 → 总14
	c.Add("k3", testValue{5}) // 触发淘汰 → 总14（k2 + k3）

	if c.Len() != 2 {
		t.Errorf("expected 2 items, got %d", c.Len())
	}

	// 验证k1被淘汰
	if _, ok := c.Get("k1"); ok {
		t.Error("k1 should have been evicted")
	}

	// 验证缓存总大小正确
	expectedSize := int64(2*2 + 2*5) // 两个条目，每个key长2，value 5
	if c.nBytes != expectedSize {
		t.Errorf("expected size %d, got %d", expectedSize, c.nBytes)
	}
}

func TestCache_UpdateExisting(t *testing.T) {
	c := New(100, nil)
	key := "key"
	oldVal := testValue{10}
	newVal := testValue{20}

	c.Add(key, oldVal)
	c.Add(key, newVal) // 更新值

	// 验证新值生效
	if val, _ := c.Get(key); val != newVal {
		t.Errorf("expected %v, got %v", newVal, val)
	}

	// 验证大小更新正确
	expectedSize := int64(len(key)) + int64(newVal.Len())
	if c.nBytes != expectedSize {
		t.Errorf("expected size %d, got %d", expectedSize, c.nBytes)
	}
}

func TestCache_EvictionCallback(t *testing.T) {
	var evictedKey string
	var evictedValue Value
	c := New(10, func(key string, value Value) {
		evictedKey = key
		evictedValue = value
	})

	// 添加会触发淘汰的条目
	c.Add("longkey", testValue{8}) // 大小: 6 + 8 = 14 > 10 → 立即淘汰
	if evictedKey != "longkey" || evictedValue == nil {
		t.Errorf("eviction callback not triggered properly")
	}
}

func TestCache_ZeroCapacity(t *testing.T) {
	c := New(0, nil)
	c.Add("k1", testValue{5})

	if c.Len() != 0 {
		t.Error("cache should never store items with 0 capacity")
	}
}

func TestCache_RemoveOldest(t *testing.T) {
	c := New(30, nil)
	c.Add("k1", testValue{10}) // 大小2+10=12
	c.Add("k2", testValue{15}) // 大小2+15=17 → 总29

	c.RemoveOldest()
	if _, ok := c.Get("k1"); ok {
		t.Error("k1 should be removed")
	}
	if c.Len() != 1 || c.nBytes != 17 {
		t.Error("size tracking incorrect after removal")
	}
}

func TestCache_LRUOrder(t *testing.T) {
	c := New(30, nil)

	// 初始添加顺序
	c.Add("k1", testValue{5}) // 7
	c.Add("k2", testValue{5}) // 7 → 总14

	// 访问k1使其成为最新
	c.Get("k1")

	// 添加新条目触发淘汰（总容量30）
	c.Add("k3", testValue{15}) // 2+15=17 → 总14+17=31 → 淘汰k2（当前最旧）

	if _, ok := c.Get("k2"); ok {
		t.Error("k2 should be evicted")
	}
}
