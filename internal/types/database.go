package types

import "time"

type DatabaseConfig struct {
	Driver                string        `yaml:"driver"`
	Source                string        `yaml:"source"`
	MaxConnections        int64         `yaml:"max_connections"`
	MaxConnectionLifetime time.Duration `yaml:"max_connection_lifetime"`
	ConnectTimeout        time.Duration `yaml:"connect_timeout"`
	ReadTimeout           time.Duration `yaml:"read_timeout"`
	WriteTimeout          time.Duration `yaml:"write_timeout"`
}
