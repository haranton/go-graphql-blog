package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env string `yaml:"env" env-required:"true"`

	App struct {
		Port int `yaml:"port" env-required:"true"`
	} `yaml:"app"`

	Database struct {
		Host     string `yaml:"host" env-required:"true"`
		Port     int    `yaml:"port" env-required:"true"`
		User     string `yaml:"user" env-required:"true"`
		Password string `yaml:"password" env-required:"true"`
		Name     string `yaml:"name" env-required:"true"`
	} `yaml:"database"`

	Migrations struct {
		Path string `yaml:"path" env-required:"true"`
	} `yaml:"migrations"`

	Storage struct {
		Type string `yaml:"type" env-required:"true"`
	} `yaml:"storage"`
}

type Flags struct {
	ConfigPath  string
	StorageType string
}

func MustLoad() *Config {
	flags := parseFlags()

	if flags.ConfigPath == "" {
		panic("config path is empty (use --config or CONFIG_PATH)")
	}

	if _, err := os.Stat(flags.ConfigPath); os.IsNotExist(err) {
		panic("config file not found: " + flags.ConfigPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(flags.ConfigPath, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	typeStorage := flags.StorageType
	if typeStorage == "" {
		panic("storage type is empty (use --storage)")
	}
	cfg.Storage.Type = typeStorage

	return &cfg
}

func parseFlags() *Flags {
	var f Flags
	flag.StringVar(&f.ConfigPath, "config", "", "path to config file (or use CONFIG_PATH)")
	flag.StringVar(&f.StorageType, "storage", "", "type of storage: memory or postgres")
	flag.Parse()
	return &f
}
