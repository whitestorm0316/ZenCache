package http

import (
	"net/http"
	"zencache/internal/cache"
	v1 "zencache/internal/transport/api/v1"

	"github.com/gin-gonic/gin"
)

type Server struct {
	ginEngine   *gin.Engine
	cacheEngine *cache.Engine
	addr        string
}

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
