package generator

import (
	"bbcli/internal/config"
	"bbcli/templates"
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// TemplatePath 模板路径映射
type TemplatePath struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

// PathConfig 路径配置文件
type PathConfig struct {
	Templates []TemplatePath `json:"templates"`
}

// Generator 代码生成器
type Generator struct {
	cfg         *config.BBConfig
	dir         string
	templates   map[string]*template.Template
	templateDir string
}

// NewGenerator 创建生成器实例
func NewGenerator(dir string, cfg *config.BBConfig) *Generator {
	g := &Generator{
		cfg:         cfg,
		dir:         dir,
		templates:   make(map[string]*template.Template),
		templateDir: getTemplateDir(),
	}
	g.loadTemplates()
	return g
}

// getTemplateDir 获取模板目录
func getTemplateDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	return filepath.Join(home, ".bbcli", "templates")
}

// ensureTemplates 确保模板存在
func (g *Generator) ensureTemplates() error {
	// 检查本地是否有模板
	if _, err := os.Stat(g.templateDir); os.IsNotExist(err) {
		// 下载模板
		if err := g.downloadTemplates(); err != nil {
			return fmt.Errorf("下载模板失败: %w", err)
		}
	}
	return nil
}

// downloadTemplates 从 GitHub 下载模板
func (g *Generator) downloadTemplates() error {
	fmt.Println("正在下载模板...")

	// 创建模板目录
	if err := os.MkdirAll(g.templateDir, 0755); err != nil {
		return err
	}

	// 复制内置模板到用户目录
	files, err := templates.Templates.ReadDir(".")
	if err != nil {
		return err
	}
	for _, f := range files {
		data, err := templates.Templates.ReadFile(f.Name())
		if err != nil {
			continue
		}
		path := filepath.Join(g.templateDir, f.Name())
		os.WriteFile(path, data, 0644)
	}
	fmt.Println("模板已复制到:", g.templateDir)

	// 尝试从 GitHub 下载（备用方案）
	githuURL := "https://raw.githubusercontent.com/baibanzz/bbcli/main/templates"
	remoteFiles := []string{"path.tmpl", "main.tmpl", "router.tmpl", "handler.tmpl", "middleware.tmpl", "config.tmpl", "gomod.tmpl"}

	for _, file := range remoteFiles {
		url := githuURL + "/" + file
		resp, err := http.Get(url)
		if err != nil {
			continue
		}
		if resp.StatusCode == 200 {
			data := new(bytes.Buffer)
			data.ReadFrom(resp.Body)
			resp.Body.Close()
			path := filepath.Join(g.templateDir, file)
			os.WriteFile(path, data.Bytes(), 0644)
		}
	}

	return nil
}

// loadTemplates 加载所有模板文件
func (g *Generator) loadTemplates() {
	// 确保模板存在
	if err := g.ensureTemplates(); err != nil {
		fmt.Printf("警告: %v\n", err)
	}

	funcMap := template.FuncMap{
		"lower":     strings.ToLower,
		"firstChar": firstChar,
		"restChars": restChars,
	}

	// 读取 path.tmpl
	pathConfig, err := g.loadPathConfig()
	if err != nil {
		fmt.Printf("加载 path.tmpl 失败: %v\n", err)
		return
	}

	for _, tp := range pathConfig.Templates {
		tmplPath := filepath.Join(g.templateDir, tp.Name)
		tmpl, err := template.New(tp.Name).Funcs(funcMap).ParseFiles(tmplPath)
		if err != nil {
			fmt.Printf("加载模板 %s 失败: %v\n", tp.Name, err)
			continue
		}
		g.templates[tp.Name] = tmpl
	}
}

// loadPathConfig 加载路径配置
func (g *Generator) loadPathConfig() (*PathConfig, error) {
	pathFile := filepath.Join(g.templateDir, "path.tmpl")
	data, err := os.ReadFile(pathFile)
	if err != nil {
		return nil, err
	}

	var cfg PathConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
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
	pathConfig, err := g.loadPathConfig()
	if err != nil {
		return fmt.Errorf("加载路径配置失败: %w", err)
	}

	for _, tp := range pathConfig.Templates {
		if err := g.generateFromPath(tp); err != nil {
			return fmt.Errorf("生成 %s 失败: %w", tp.Name, err)
		}
	}

	return nil
}

// generateFromPath 根据路径配置生成文件
func (g *Generator) generateFromPath(tp TemplatePath) error {
	tmpl, ok := g.templates[tp.Name]
	if !ok {
		return fmt.Errorf("模板 %s 未加载", tp.Name)
	}

	var buf bytes.Buffer
	var err error

	// 根据模板类型准备数据
	switch tp.Name {
	case "handler.tmpl":
		if len(g.cfg.Routes) == 0 {
			return nil
		}
		data := struct{ Handler string }{Handler: g.cfg.Routes[0].Handler}
		err = tmpl.Execute(&buf, data)
	case "middleware.tmpl":
		if len(g.cfg.Middlewares) == 0 {
			return nil
		}
		data := struct{ Name string }{Name: g.cfg.Middlewares[0]}
		err = tmpl.Execute(&buf, data)
	case "router.tmpl":
		data := struct {
			Project     config.ProjectConfig
			Routes      []config.RouteConfig
			Middlewares []string
		}{
			Project:     g.cfg.Project,
			Routes:      g.cfg.Routes,
			Middlewares: g.cfg.Middlewares,
		}
		err = tmpl.Execute(&buf, data)
	case "config.tmpl":
		err = tmpl.Execute(&buf, g.cfg.Project)
	default:
		data := struct{ Project config.ProjectConfig }{Project: g.cfg.Project}
		err = tmpl.Execute(&buf, data)
	}

	if err != nil {
		return err
	}

	// 解析输出路径（支持模板变量）
	outPath, err := template.New("").Parse(tp.Path)
	if err != nil {
		return err
	}

	var pathBuf bytes.Buffer
	switch tp.Name {
	case "handler.tmpl":
		if len(g.cfg.Routes) == 0 {
			return nil
		}
		data := struct{ Handler string }{Handler: g.cfg.Routes[0].Handler}
		outPath.Execute(&pathBuf, data)
	case "middleware.tmpl":
		if len(g.cfg.Middlewares) == 0 {
			return nil
		}
		data := struct{ Name string }{Name: g.cfg.Middlewares[0]}
		outPath.Execute(&pathBuf, data)
	default:
		data := struct{ Project config.ProjectConfig }{Project: g.cfg.Project}
		outPath.Execute(&pathBuf, data)
	}

	relPath := pathBuf.String()
	if err := g.ensureDir(filepath.Dir(relPath)); err != nil {
		return err
	}

	return g.writeFile(relPath, buf.Bytes())
}

// ensureDir 确保目录存在
func (g *Generator) ensureDir(path string) error {
	if path == "" {
		return nil
	}
	fullPath := filepath.Join(g.dir, path)
	return os.MkdirAll(fullPath, 0755)
}

// writeFile 写入文件
func (g *Generator) writeFile(relPath string, content []byte) error {
	fullPath := filepath.Join(g.dir, relPath)
	return os.WriteFile(fullPath, content, 0644)
}
