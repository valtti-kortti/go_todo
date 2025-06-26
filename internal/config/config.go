package config

const EnvConfigPath = ".env"

type AppConfig struct {
	Loglevel string
	Rest     Rest
}

type Rest struct {
	ListenAddress string `envconfig:"PORT" required:"true"`
	WriteTimeout  string `envconfig:"WRITE_TIMEOUT" required:"true"`
	ServerName    string `envconfig:"SERVER_NAME" required:"true"`
	Token         string `envconfig:"TOKEN" required:"true"`
}
