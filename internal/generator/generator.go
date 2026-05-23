package generator

import (
	"bbcli/internal/config"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// Generator 代码生成器
type Generator struct {
	cfg *config.BBConfig
	dir string
}

// NewGenerator 创建生成器实例
func NewGenerator(dir string, cfg *config.BBConfig) *Generator {
	return &Generator{
		cfg: cfg,
		dir: dir,
	}
}

// fileStrategy 定义文件处理策略
type fileStrategy int

const (
	strategyOverwrite fileStrategy = iota
	strategyIgnore
	strategyCreate
)

// shouldHandle 决定是否处理文件
func (g *Generator) shouldHandle(relPath string, strategy fileStrategy) bool {
	switch strategy {
	case strategyOverwrite:
		return true
	case strategyIgnore:
		_, err := os.Stat(filepath.Join(g.dir, relPath))
		return os.IsNotExist(err)
	case strategyCreate:
		fullPath := filepath.Join(g.dir, relPath)
		_, err := os.Stat(fullPath)
		return os.IsNotExist(err)
	default:
		return false
	}
}

// Generate 生成所有代码
func (g *Generator) Generate() error {
	// 生成 go.mod
	if g.shouldHandle("go.mod", strategyOverwrite) {
		if err := g.generateGoMod(); err != nil {
			return fmt.Errorf("生成 go.mod 失败: %w", err)
		}
	}

	// 生成 main.go
	if g.shouldHandle("main.go", strategyOverwrite) {
		if err := g.generateMainGo(); err != nil {
			return fmt.Errorf("生成 main.go 失败: %w", err)
		}
	}

	// 生成 router.go
	if g.shouldHandle("router/router.go", strategyOverwrite) {
		if err := g.generateRouter(); err != nil {
			return fmt.Errorf("生成 router.go 失败: %w", err)
		}
	}

	// 生成 handlers
	for _, route := range g.cfg.Routes {
		if err := g.generateHandler(route.Handler); err != nil {
			return fmt.Errorf("生成 handler %s 失败: %w", route.Handler, err)
		}
	}

	// 生成 middleware（如不存在）
	for _, mw := range g.cfg.Middlewares {
		if err := g.generateMiddleware(mw); err != nil {
			return fmt.Errorf("生成 middleware %s 失败: %w", mw, err)
		}
	}

	// 生成数据库配置（如启用）
	if g.cfg.Project.Database.Type != "" {
		if err := g.generateDatabase(); err != nil {
			return fmt.Errorf("生成数据库配置失败: %w", err)
		}
	}

	return nil
}

// templateData 用于模板渲染的数据
type templateData struct {
	Project        config.ProjectConfig
	Routes         []config.RouteConfig
	Middleware     []string
	Handler        string
	MiddlewareName string
}

// renderTemplate 渲染模板
func (g *Generator) renderTemplate(tmplStr string, data interface{}) (string, error) {
	tmpl, err := template.New("").Parse(tmplStr)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// ensureDir 确保目录存在
func (g *Generator) ensureDir(path string) error {
	fullPath := filepath.Join(g.dir, path)
	return os.MkdirAll(fullPath, 0755)
}

// writeFile 写入文件
func (g *Generator) writeFile(relPath, content string) error {
	fullPath := filepath.Join(g.dir, relPath)
	return os.WriteFile(fullPath, []byte(content), 0644)
}

// 生成 go.mod
func (g *Generator) generateGoMod() error {
	content := fmt.Sprintf(`module %s

go 1.21

require (
	github.com/gin-gonic/gin v1.9.1
	gorm.io/gorm v1.25.2
)

require (
	github.com/bytedance/sonic v1.9.1 // indirect
	github.com/chenzhuoyu/base64x v0.0.0-20221115062448-fe3a3abad311 // indirect
	github.com/gabriel-vasile/mimetype v1.4.2 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.14.0 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.2.4 // indirect
	github.com/leodido/go-urn v1.2.4 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.0.8 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.11 // indirect
	golang.org/x/arch v0.3.0 // indirect
	golang.org/x/crypto v0.9.0 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
`, g.cfg.Project.Name)
	return g.writeFile("go.mod", content)
}

// 生成 main.go
func (g *Generator) generateMainGo() error {
	content := `package main

import (
	"fmt"
	"log"
	"%[1]s/config"
	"%[1]s/router"
)

func main() {
	// 加载配置
	cfg := config.Load()
	
	// 创建 Gin 实例
	r := router.Setup(cfg)
	
	// 启动服务器
	addr := fmt.Sprintf(":%d", cfg.Project.Port)
	log.Printf("服务器启动在 http://localhost%s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
`
	content = fmt.Sprintf(content, g.cfg.Project.Name)
	return g.writeFile("main.go", content)
}

// 生成 router.go
func (g *Generator) generateRouter() error {
	if err := g.ensureDir("router"); err != nil {
		return err
	}

	var buf bytes.Buffer
	buf.WriteString("package router\n\n")
	buf.WriteString("import (\n")
	buf.WriteString("\t\"" + g.cfg.Project.Name + "/config\"\n")
	buf.WriteString("\t\"" + g.cfg.Project.Name + "/handler\"\n")

	// 导入 middleware 包
	if len(g.cfg.Middlewares) > 0 {
		buf.WriteString("\t\"" + g.cfg.Project.Name + "/middleware\"\n")
	}

	buf.WriteString("\t\"github.com/gin-gonic/gin\"\n")
	buf.WriteString(")\n\n")

	// Setup 函数
	buf.WriteString("func Setup(cfg *config.Config) *gin.Engine {\n")
	buf.WriteString("\tr := gin.Default()\n\n")

	// 注册中间件
	for _, mw := range g.cfg.Middlewares {
		mwVar := strings.ToLower(mw[0:1]) + mw[1:]
		buf.WriteString(fmt.Sprintf("\tr.Use(middleware.%s())\n", mwVar))
	}
	buf.WriteString("\n")

	// 注册路由
	for _, route := range g.cfg.Routes {
		handlerVar := strings.ToLower(route.Handler[0:1]) + route.Handler[1:]
		buf.WriteString(fmt.Sprintf("\tr.%s(\"%s\", handler.%s)\n",
			route.Method, route.Path, handlerVar))
	}

	buf.WriteString("\n\treturn r\n}\n")

	return g.writeFile("router/router.go", buf.String())
}

// 生成 handler
func (g *Generator) generateHandler(handlerName string) error {
	if err := g.ensureDir("handler"); err != nil {
		return err
	}

	data := templateData{
		Handler: handlerName,
	}

	tmpl := `package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func {{.Handler}}(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, World!",
	})
}
`

	content, err := g.renderTemplate(tmpl, data)
	if err != nil {
		return err
	}

	// 小写处理文件名
	fileName := strings.ToLower(handlerName[0:1]) + handlerName[1:] + ".go"
	return g.writeFile("handler/"+fileName, content)
}

// 生成 middleware
func (g *Generator) generateMiddleware(mwName string) error {
	if err := g.ensureDir("middleware"); err != nil {
		return err
	}

	data := templateData{
		MiddlewareName: mwName,
	}

	tmpl := `package middleware

import (
	"github.com/gin-gonic/gin"
)

func {{.MiddlewareName}}() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: 实现中间件逻辑
		c.Next()
	}
}
`

	content, err := g.renderTemplate(tmpl, data)
	if err != nil {
		return err
	}

	// 小写处理文件名
	fileName := strings.ToLower(mwName[0:1]) + mwName[1:] + ".go"
	return g.writeFile("middleware/"+fileName, content)
}

// 生成数据库配置
func (g *Generator) generateDatabase() error {
	if err := g.ensureDir("config"); err != nil {
		return err
	}

	var buf bytes.Buffer
	buf.WriteString("package config\n\n")

	switch g.cfg.Project.Database.Type {
	case "mysql":
		buf.WriteString("var DB *gorm.DB\n\n")
		buf.WriteString("func InitDB() error {\n")
		buf.WriteString("\tdsn := \"root:@tcp(localhost:3306)/myapp?charset=utf8mb4&parseTime=True&loc=Local\"\n")
		buf.WriteString("\tdb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})\n")
		buf.WriteString("\tif err != nil {\n")
		buf.WriteString("\t\treturn err\n")
		buf.WriteString("\t}\n")
		buf.WriteString("\tDB = db\n")
		buf.WriteString("\treturn nil\n}\n")

	case "postgres":
		buf.WriteString("var DB *gorm.DB\n\n")
		buf.WriteString("func InitDB() error {\n")
		buf.WriteString("\tdsn := \"host=localhost user=postgres password= dbname=myapp port=5432 sslmode=disable\"\n")
		buf.WriteString("\tdb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})\n")
		buf.WriteString("\tif err != nil {\n")
		buf.WriteString("\t\treturn err\n")
		buf.WriteString("\t}\n")
		buf.WriteString("\tDB = db\n")
		buf.WriteString("\treturn nil\n}\n")
	}

	buf.WriteString("func Load() *Config {\n")
	buf.WriteString("\treturn &Config{\n")
	buf.WriteString("\t\tProject: ProjectConfig{\n")
	buf.WriteString(fmt.Sprintf("\t\t\tName: \"%s\",\n", g.cfg.Project.Name))
	buf.WriteString(fmt.Sprintf("\t\t\tPort: %d,\n", g.cfg.Project.Port))
	buf.WriteString("\t\t\tDatabase: DatabaseConfig{\n")
	buf.WriteString(fmt.Sprintf("\t\t\t\tType: \"%s\",\n", g.cfg.Project.Database.Type))
	buf.WriteString(fmt.Sprintf("\t\t\t\tHost: \"%s\",\n", g.cfg.Project.Database.Host))
	buf.WriteString(fmt.Sprintf("\t\t\t\tPort: %d,\n", g.cfg.Project.Database.Port))
	buf.WriteString(fmt.Sprintf("\t\t\t\tUser: \"%s\",\n", g.cfg.Project.Database.User))
	buf.WriteString(fmt.Sprintf("\t\t\t\tPassword: \"%s\",\n", g.cfg.Project.Database.Password))
	buf.WriteString(fmt.Sprintf("\t\t\t\tDatabase: \"%s\",\n", g.cfg.Project.Database.Database))
	buf.WriteString("\t\t\t},\n")
	buf.WriteString("\t\t},\n")
	buf.WriteString("\t}\n")
	buf.WriteString("}\n")

	return g.writeFile("config/config.go", buf.String())
}
