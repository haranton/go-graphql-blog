package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env string `yaml:"env"`

	App struct {
		Port int `yaml:"port"`
	} `yaml:"app"`

	Database struct {
		Host     string `yaml:"host"` //todo добавть reqired чтобы выдавать ошибку если параметров нет
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
	} `yaml:"database"`

	Migrations struct {
		Path string `yaml:"path"`
	} `yaml:"migrations"`

	NotifyServiceAddr string `yaml:"NotifyServiceAddr"`
	Storage           struct {
		Type string `yaml:"type"`
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
