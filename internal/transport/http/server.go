package http

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"io"
	"net/http"
	"sync"
	"zencache/internal/cache"
	"zencache/internal/config"
	"zencache/internal/consistenthash"
	"zencache/internal/peers"
	v1 "zencache/internal/transport/api/v1"

	"github.com/gin-gonic/gin"
)

type httpGetter struct {
	baseURL string
}

// 从远程获取
func (h *httpGetter) Get(group string, key string) ([]byte, error) {
	requestUrl := fmt.Sprint(h.baseURL + v1.GET_KEY)
	response, err := http.Post(requestUrl, "application/json", bytes.NewBuffer([]byte(fmt.Sprintf(`{"group":%s,"key":%s}`, group, key))))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status :%s", response.Status)
	}
	return io.ReadAll(response.Body)
}

type Server struct {
	ginEngine       *gin.Engine
	cacheEngine     *cache.Engine
	addr            string
	self            string // 192.168.1.134:8080
	baseUrl         string // https://
	mutex           sync.Mutex
	peersHttpGetter map[string]*httpGetter
	peers           *consistenthash.Map
}

// NewWithConfig 使用配置创建新的Server实例
func NewWithConfig(conf *config.Config) *Server {
	// 创建缓存引擎
	cacheEngine := cache.NewEngine()

	// 创建一致性哈希实例
	peers := consistenthash.New(conf, func(b []byte) uint32 {
		hash := sha1.Sum(b)
		return binary.LittleEndian.Uint32(hash[:4])
	})

	// 创建HTTP服务器
	ginEngine := gin.Default()

	// 创建并初始化服务器实例
	s := &Server{
		ginEngine:       ginEngine,
		cacheEngine:     cacheEngine,
		addr:            fmt.Sprintf("%s:%d", conf.HTTP.Address, conf.HTTP.Port),
		peers:           peers,
		peersHttpGetter: make(map[string]*httpGetter),
	}

	// 注册路由
	s.ginEngine.POST(v1.STORE_KEY, s.handleStoreKey)
	s.ginEngine.POST(v1.GET_KEY, s.handleGetKey)
	s.ginEngine.POST(v1.DELETE_KEY, s.handleDeleteKey)

	return s
}

func (s *Server) PickPeer(key string) (peers.PeerGetter, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	peer := s.peers.Get(key)
	if peer == s.self {
		return nil, false
	}
	getter, ok := s.peersHttpGetter[peer]
	return getter, ok
}

var _ peers.PeersPicker = (*Server)(nil)

func New(addr string) *Server {
	cacheEngine := cache.NewEngine()
	ginEngine := gin.Default()
	s := &Server{
		ginEngine:   ginEngine,
		cacheEngine: cacheEngine,
		addr:        addr,
	}

	s.ginEngine.POST(v1.STORE_KEY, s.handleStoreKey)
	s.ginEngine.POST(v1.GET_KEY, s.handleGetKey)
	s.ginEngine.POST(v1.DELETE_KEY, s.handleDeleteKey)

	return s
}

func (s *Server) SetNodes(nodes ...string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.peers.Add(nodes...)
	for _, node := range nodes {
		s.peersHttpGetter[node] = &httpGetter{
			baseURL: s.baseUrl + node,
		}
	}
}

func (s *Server) handleStoreKey(c *gin.Context) {
	var req v1.StoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, v1.Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	group := s.cacheEngine.GetGroup(req.Group)
	if group == nil {
		s.cacheEngine.AddGroup(req.Group, nil, 1<<20) // 默认1MB
		group = s.cacheEngine.GetGroup(req.Group)
		group.RegisterPicker(s)
	}

	if err := group.Add(req.Key, cache.NewByteView(req.Value)); err != nil {
		c.JSON(http.StatusInternalServerError, v1.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, v1.Response{
		Code:    http.StatusOK,
		Message: "success",
	})
}

func (s *Server) handleGetKey(c *gin.Context) {
	var req v1.GetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, v1.Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	group := s.cacheEngine.GetGroup(req.Group)
	if group == nil {
		c.JSON(http.StatusNotFound, v1.Response{
			Code:    http.StatusNotFound,
			Message: "group not found",
		})
		return
	}

	value, err := group.Get(req.Key)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == cache.ErrKeyNotFound {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, v1.Response{
			Code:    int32(statusCode),
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, v1.Response{
		Code:    http.StatusOK,
		Message: "success",
		Data:    value.ByteSlices(),
	})
}

func (s *Server) handleDeleteKey(c *gin.Context) {
	var req v1.DeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, v1.Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	group := s.cacheEngine.GetGroup(req.Group)
	if group == nil {
		c.JSON(http.StatusNotFound, v1.Response{
			Code:    http.StatusNotFound,
			Message: "group not found",
		})
		return
	}

	// TODO: 实现删除键的功能
	c.JSON(http.StatusOK, v1.Response{
		Code:    http.StatusOK,
		Message: "success",
	})
}

func (s *Server) Run() {
	s.ginEngine.Run(s.addr)
}
