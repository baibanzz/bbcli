package config

import "github.com/baibanzz/jdk/model"

type Config struct {
	Gin        Gin              `yaml:"gin"`
	Mysql      model.Mysql      `yaml:"mysql"`
	Redis      model.Redis      `yaml:"redis"`
	Sqlite     model.Sqlite3    `yaml:"sqlite"`
	PostgreSql model.PostgreSql `yaml:"postgresql"`
	Etcd       model.Etcd       `yaml:"etcd"`
}

type Gin struct {
	Addr    string `yaml:"addr"`
	Port    string `yaml:"port"`
	TlsCert string `yaml:"tlsCert"`
	TlsKey  string `yaml:"tlsKey"`
}
