package model

type Config struct {
	Name string //项目名称
	//数据库一块
	UseMySQL      bool
	UsePostgreSQL bool
	UseSQLite3    bool
	UseRedis      bool

	//中间件一块
	UseJwt  bool
	UseCore bool
	//扩展组件一块
	UseEtcd bool
}
