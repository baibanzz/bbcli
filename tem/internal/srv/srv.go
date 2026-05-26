package srv

import (
	"tem/internal/config"

	"github.com/baibanzz/jdk/core"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Srv struct {
	Config     config.Config
	Mysql      *gorm.DB
	Sqlite     *gorm.DB
	PostgreSql *gorm.DB
	Redis      *redis.Client
	Etcd       *core.Etcd
}

func NewSrv(config config.Config) (*Srv, error) {
	var (
		err error
		srv = &Srv{
			Config: config,
		}
	)
	if srv.Mysql, err = core.NewDB(config.Mysql, nil); err != nil {
		return nil, err
	}
	if srv.Sqlite, err = core.NewDB(config.Sqlite, nil); err != nil {
		return nil, err
	}
	if srv.PostgreSql, err = core.NewDB(config.PostgreSql, nil); err != nil {
		return nil, err
	}
	if srv.Redis, err = core.NewRedis(config.Redis); err != nil {
		return nil, err
	}
	if srv.Etcd, err = core.NewEtcd(config.Etcd); err != nil {
		return nil, err
	}
	return srv, nil
}
