package config

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

const defaultPath = "./config/config.yaml"

var (
	errFileIsNotExist = errors.New("file is not exist")
	errParseToStruct  = errors.New("failed to parse")
)

type Config struct {
	App     App     `yaml:"app"`
	HTTP    HTTP    `yaml:"http"`
	Cache   Cache   `yaml:"cache"`
	SS      SS      `yaml:"ss"`
	Cleaner Cleaner `yaml:"cleaner"`
}

// App is application config
type App struct {
	Name             string `yaml:"name"`
	Version          string `yaml:"version"`
	Env              string `yaml:"env"`
	TmpDirectoryPath string `yaml:"tmp_directory_path"`
	TmpFilePath      string `yaml:"tmp_file_path"`
}

// HTTP is webserver's config
type HTTP struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

// Cache is cache dsn config
type Cache struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

// SS is shared storage config
type SS struct {
	Host               string `yaml:"host"`
	Port               string `yaml:"port"`
	User               string `yaml:"user"`
	Password           string `yaml:"password"`
	ShareName          string `yaml:"share_name"`
	ConnectionPoolSize int    `yaml:"connection_pool_size"`
}

type Cleaner struct {
	TimeOffset time.Duration `yaml:"time_offset"`
}

func MustLoad() (*Config, error) {
	cfg := &Config{}

	f, err := os.Open(defaultPath)
	if err != nil {
		return cfg, fmt.Errorf("%w in %s", errFileIsNotExist, defaultPath)
	}
	defer f.Close()

	err = yaml.NewDecoder(f).Decode(cfg)
	if err != nil {
		return cfg, errParseToStruct
	}

	return cfg, nil
}
