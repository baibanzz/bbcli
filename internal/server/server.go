package server

import (
	"bbcli/internal/config"
	"bbcli/internal/generator"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// Server Gin 服务器
type Server struct {
	cfg  *config.BBConfig
	dir  string
	addr string
}

// NewServer 创建服务器实例
func NewServer(dir string, addr string) (*Server, error) {
	cfg, err := config.LoadConfig(dir)
	if err != nil {
		return nil, fmt.Errorf("加载配置失败: %w", err)
	}

	return &Server{
		cfg:  cfg,
		dir:  dir,
		addr: addr,
	}, nil
}

// Start 启动服务器
func (s *Server) Start() error {
	r := gin.Default()

	// API 路由
	r.GET("/api/config", s.getConfig)
	r.POST("/api/config", s.saveConfig)
	r.POST("/api/generate", s.generate)

	// 提供前端静态文件
	r.Static("/static", "./web")
	r.GET("/", s.serveIndex)

	return r.Run(s.addr)
}

// getConfig 获取配置
func (s *Server) getConfig(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    s.cfg,
		"exists":  config.ConfigExists(s.dir),
	})
}

// saveConfig 保存配置
func (s *Server) saveConfig(c *gin.Context) {
	var cfg config.BBConfig
	if err := c.ShouldBindJSON(&cfg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := config.SaveConfig(s.dir, &cfg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	s.cfg = &cfg
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "配置已保存",
	})
}

// generate 生成代码
func (s *Server) generate(c *gin.Context) {
	g := generator.NewGenerator(s.dir, s.cfg)

	if err := g.Generate(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "代码已生成",
	})
}

// serveIndex 提供前端入口
func (s *Server) serveIndex(c *gin.Context) {
	indexPath := filepath.Join("web", "index.html")
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		c.String(http.StatusNotFound, "前端文件未找到，请先构建前端")
		return
	}
	c.File(indexPath)
}
