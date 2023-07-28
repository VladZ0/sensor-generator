package config

import (
	"log"
	"sensors-generator/pkg/client/postgresql"
	"sensors-generator/pkg/client/redis"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	IsDebug      bool `yaml:"is_debug" env-default:"false"`
	IsProduction bool `yaml:"is_production" env-default:"true"`

	Listen struct {
		Type       string `yaml:"type" env-default:"port" env-description:"port or sock"`
		BindIP     string `yaml:"bind_ip" env-default:"127.0.0.1"`
		Port       string `yaml:"port" env-default:"8080"`
		SocketFile string `yaml:"socket_file" env-default:"app.sock"`
	} `yaml:"listen"`

	AppConfig struct {
		LogLevel string `yaml:"log_level" env-default:"trace"`
	} `yaml:"app_config"`

	CorsConfig struct {
		AllowedMethods     []string `yaml:"allowed_methods"`
		AllowedOrigins     []string `yaml:"allowed_origins"`
		AllowedHeaders     []string `yaml:"allowed_headers"`
		ExposedHeaders     []string `yaml:"exposed_headers"`
		AllowCredentials   bool     `yaml:"allow_credentials"`
		OptionsPassthrough bool     `yaml:"options_passthrough"`
	} `yaml:"cors_config"`

	PgConfig    postgresql.PgConfig `yaml:"pg_config"`
	RedisConfig redis.RedisConfig   `yaml:"redis_config"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{}

		if err := cleanenv.ReadConfig("config.yml", instance); err != nil {
			helpText := "Messenger project"
			help, _ := cleanenv.GetDescription(instance, &helpText)
			log.Print(help)
			log.Fatal(err)
		}
	})

	return instance
}
