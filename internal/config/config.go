package config

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Type     string `json:"type"` // mysql, postgres, sqlite
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

// ProjectConfig 项目配置
type ProjectConfig struct {
	Name     string         `json:"name"`
	Port     int            `json:"port"`
	Database DatabaseConfig `json:"database"`
}

// RouteConfig 路由配置
type RouteConfig struct {
	Method     string   `json:"method"`
	Path       string   `json:"path"`
	Handler    string   `json:"handler"`
	Middleware []string `json:"middleware"`
}

// BBConfig 主配置结构
type BBConfig struct {
	Project     ProjectConfig `json:"project"`
	Routes      []RouteConfig `json:"routes"`
	Middlewares []string      `json:"middlewares"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *BBConfig {
	return &BBConfig{
		Project: ProjectConfig{
			Name: "myapp",
			Port: 8080,
			Database: DatabaseConfig{
				Type:     "mysql",
				Host:     "localhost",
				Port:     3306,
				User:     "root",
				Password: "",
				Database: "myapp",
			},
		},
		Routes: []RouteConfig{
			{
				Method:     "GET",
				Path:       "/",
				Handler:    "HomeHandler",
				Middleware: []string{"Logger", "Recovery"},
			},
		},
		Middlewares: []string{"Logger", "Recovery", "CORS"},
	}
}
