package server

type ServerConfig struct {
	Host string `json:"host" yaml:"host" toml:"host" env:"SERVER_HOST"`
	Port string `json:"port" yaml:"port" toml:"port" env:"SERVER_PORT"`
}