package config

type DbConfig struct {
	DsName   string `env:"NAME" envDefault:"datasource"`
	Driver   string `env:"DRIVER" envDefault:"pgx"`
	Host     string `env:"HOST" envDefault:"localhost"`
	Port     string `env:"PORT" envDefault:"5432"`
	User     string `env:"USER" envDefault:"user"`
	Password string `env:"PASSWORD" envDefault:"pass"`
	Database string `env:"DATABASE" envDefault:"test"`
	MaxConn  string `env:"CONMAX" envDefault:"10"`
}
