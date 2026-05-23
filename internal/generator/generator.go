package generator

import (
	"bbcli/internal/config"
	"bytes"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed *.tmpl
var templateFiles embed.FS

// Generator 代码生成器
type Generator struct {
	cfg       *config.BBConfig
	dir       string
	templates map[string]*template.Template
}

// NewGenerator 创建生成器实例
func NewGenerator(dir string, cfg *config.BBConfig) *Generator {
	g := &Generator{
		cfg:       cfg,
		dir:       dir,
		templates: make(map[string]*template.Template),
	}
	g.loadTemplates()
	return g
}

// loadTemplates 加载所有模板文件
func (g *Generator) loadTemplates() {
	funcMap := template.FuncMap{
		"lower":     strings.ToLower,
		"firstChar": firstChar,
		"restChars": restChars,
	}

	templates := []string{
		"gomod.tmpl",
		"main.tmpl",
		"router.tmpl",
		"handler.tmpl",
		"middleware.tmpl",
		"config.tmpl",
	}

	for _, name := range templates {
		tmpl, err := template.New(name).Funcs(funcMap).ParseFS(templateFiles, name)
		if err != nil {
			fmt.Printf("加载模板 %s 失败: %v\n", name, err)
			continue
		}
		g.templates[name] = tmpl
	}
}

// firstChar 返回首字符小写
func firstChar(s string) string {
	if len(s) == 0 {
		return ""
	}
	return strings.ToLower(s[:1])
}

// restChars 返回剩余字符
func restChars(s string) string {
	if len(s) <= 1 {
		return ""
	}
	return s[1:]
}

// fileStrategy 定义文件处理策略
type fileStrategy int

const (
	strategyOverwrite fileStrategy = iota
	strategyIgnore
)

// shouldHandle 决定是否处理文件
func (g *Generator) shouldHandle(relPath string, strategy fileStrategy) bool {
	switch strategy {
	case strategyOverwrite:
		return true
	case strategyIgnore:
		_, err := os.Stat(filepath.Join(g.dir, relPath))
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

// ensureDir 确保目录存在
func (g *Generator) ensureDir(path string) error {
	fullPath := filepath.Join(g.dir, path)
	return os.MkdirAll(fullPath, 0755)
}

// writeFile 写入文件
func (g *Generator) writeFile(relPath string, content []byte) error {
	fullPath := filepath.Join(g.dir, relPath)
	return os.WriteFile(fullPath, content, 0644)
}

// executeTemplate 执行模板
func (g *Generator) executeTemplate(tmplName string, data interface{}) ([]byte, error) {
	tmpl, ok := g.templates[tmplName]
	if !ok {
		return nil, fmt.Errorf("模板 %s 不存在", tmplName)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// generateGoMod 生成 go.mod
func (g *Generator) generateGoMod() error {
	data := struct {
		Project config.ProjectConfig
	}{
		Project: g.cfg.Project,
	}

	content, err := g.executeTemplate("gomod.tmpl", data)
	if err != nil {
		return err
	}

	return g.writeFile("go.mod", content)
}

// generateMainGo 生成 main.go
func (g *Generator) generateMainGo() error {
	data := struct {
		Project config.ProjectConfig
	}{
		Project: g.cfg.Project,
	}

	content, err := g.executeTemplate("main.tmpl", data)
	if err != nil {
		return err
	}

	return g.writeFile("main.go", content)
}

// generateRouter 生成 router.go
func (g *Generator) generateRouter() error {
	if err := g.ensureDir("router"); err != nil {
		return err
	}

	data := struct {
		Project     config.ProjectConfig
		Routes      []config.RouteConfig
		Middlewares []string
	}{
		Project:     g.cfg.Project,
		Routes:      g.cfg.Routes,
		Middlewares: g.cfg.Middlewares,
	}

	content, err := g.executeTemplate("router.tmpl", data)
	if err != nil {
		return err
	}

	return g.writeFile("router/router.go", content)
}

// generateHandler 生成 handler
func (g *Generator) generateHandler(handlerName string) error {
	if err := g.ensureDir("handler"); err != nil {
		return err
	}

	data := struct {
		Handler string
	}{
		Handler: handlerName,
	}

	content, err := g.executeTemplate("handler.tmpl", data)
	if err != nil {
		return err
	}

	// 小写处理文件名
	fileName := firstChar(handlerName) + restChars(handlerName) + ".go"
	return g.writeFile("handler/"+fileName, content)
}

// generateMiddleware 生成 middleware
func (g *Generator) generateMiddleware(mwName string) error {
	if err := g.ensureDir("middleware"); err != nil {
		return err
	}

	data := struct {
		Name string
	}{
		Name: mwName,
	}

	content, err := g.executeTemplate("middleware.tmpl", data)
	if err != nil {
		return err
	}

	// 小写处理文件名
	fileName := firstChar(mwName) + restChars(mwName) + ".go"
	return g.writeFile("middleware/"+fileName, content)
}

// generateDatabase 生成数据库配置
func (g *Generator) generateDatabase() error {
	if err := g.ensureDir("config"); err != nil {
		return err
	}

	data := g.cfg.Project

	content, err := g.executeTemplate("config.tmpl", data)
	if err != nil {
		return err
	}

	return g.writeFile("config/config.go", content)
}
