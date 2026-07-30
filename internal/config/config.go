package config

type Config struct {
	Port string
}

func Load() Config {
	return Config{
		Port: "8080",
	}
}

func (c Config) Addr() string {
	return ":" + c.Port
}
