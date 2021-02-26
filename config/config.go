package config

import "time"

type Config struct {
	DbConfig DbConfig  `toml:"db"`
	JWT      JWTConfig `toml:"jwt"`
	Log      LogConfig `toml:"log"`
	API      APIConfig `toml:"api"`
}

type LogConfig struct {
	Path string `toml:"path"`
}

type APIConfig struct {
	Address string
}

type DbConfig struct {
	Type   string       `toml:"type"`
	MySql  MySqlConfig  `toml:"mysql"`
	Sqlite SqliteConfig `toml:"sqlite"`
}

type SqliteConfig struct {
	Path string `toml:"path"`
}

type MySqlConfig struct {
	Addr            string        `toml:"addr"`
	User            string        `toml:"user"`
	Pass            string        `toml:"pass"`
	Name            string        `toml:"name"`
	MaxOpenConn     int           `toml:"maxOpenConn"`
	MaxIdleConn     int           `toml:"maxIdleConn"`
	ConnMaxLifeTime time.Duration `toml:"connMaxLifeTime"`
}

type JWTConfig struct {
	Secret         string        `toml:"secret"`
	ExpireDuration time.Duration `toml:"expireDuration"`
}
