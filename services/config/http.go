package config

type HTTPConfig struct {
	ListenAddr string `env:"LISTEN" envDefault:"127.0.0.1:8080"`
	Sock       string `env:"SOCK"`
	Domain     string `env:"DOMAIN"   envDefault:"localhost:8080"`
}
