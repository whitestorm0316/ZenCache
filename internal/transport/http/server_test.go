package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	v1 "zencache/internal/transport/api/v1"

	"github.com/gin-gonic/gin"
)

func BenchmarkServer(b *testing.B) {
	gin.SetMode("release")
	s := New(":8080")
	// 准备测试数据
	group := "test_group"
	key := "test_key"
	value := []byte("test_value")

	// 测试存储操作
	b.Run("Store", func(b *testing.B) {
		req := &v1.StoreRequest{
			Group: group,
			Key:   key,
			Value: value,
		}
		body, _ := json.Marshal(req)
		b.ResetTimer()
		b.SetParallelism(100) // 设置并发数
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Request, _ = http.NewRequest("POST", "/store", bytes.NewReader(body))
				c.Request.Header.Set("Content-Type", "application/json")
				s.handleStoreKey(c)
			}
		})
	})

	// 测试获取操作
	b.Run("Get", func(b *testing.B) {
		req := &v1.GetRequest{
			Group: group,
			Key:   key,
		}
		body, _ := json.Marshal(req)
		b.ResetTimer()
		b.SetParallelism(100) // 设置并发数
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Request, _ = http.NewRequest("POST", "/get", bytes.NewReader(body))
				c.Request.Header.Set("Content-Type", "application/json")
				s.handleGetKey(c)
			}
		})
	})

	// 测试删除操作
	b.Run("Delete", func(b *testing.B) {
		req := &v1.DeleteRequest{
			Group: group,
			Key:   key,
		}
		body, _ := json.Marshal(req)
		b.ResetTimer()
		b.SetParallelism(100) // 设置并发数
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Request, _ = http.NewRequest("POST", "/delete", bytes.NewReader(body))
				c.Request.Header.Set("Content-Type", "application/json")
				s.handleDeleteKey(c)
			}
		})
	})
}
