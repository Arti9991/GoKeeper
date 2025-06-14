package config

import (
	"flag"
	"io"
	"os"

	"github.com/Arti9991/GoKeeper/server/internal/logger"
	"github.com/caarlos0/env"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

var FileConfig = true
var ConfigPath = "Config.yaml"

type Config struct {
	HostAddr   string `env:"HOST_ADDRESS" yaml:"host_address"`
	DBAdr      string `env:"DATABASE_URI" yaml:"database_info"`
	InFileLog  bool   `yaml:"save_log_to_file"`
	StorageDir string `env:"STORAGE_DIR" yaml:"storage_dir"`
}

// инициализация конфигурации
// (внутри есть флаг FileConfig для чтения конфигурации из файла)
func InitConf() Config {
	var conf Config

	if FileConfig {
		conf = ReadConfig(ConfigPath)
	}

	err := env.Parse(&conf)
	if err != nil {
		logger.Log.Error("Wrong config file!", zap.Error(err))
	}

	flag.StringVar(&conf.HostAddr, "a", conf.HostAddr, "server host adress")
	flag.StringVar(&conf.DBAdr, "d", conf.DBAdr, "database connetion data") //"host=localhost user=myuser password=123456 dbname=Keeper sslmode=disable"
	//flag.StringVar(&conf.AccurAddr, "r", conf.AccurAddr, "another api address")
	flag.BoolVar(&conf.InFileLog, "l", conf.InFileLog, "enable logging in file")
	flag.StringVar(&conf.StorageDir, "s", conf.StorageDir, "direktory with binary storage")
	flag.Parse()

	//CreateConfig(ConfigPath, conf)

	return conf
}

func InitConfTest() Config {

	var conf Config

	conf.HostAddr = ":8081"
	conf.InFileLog = false

	return conf
}

// чтение конфигурации из файла
func ReadConfig(cfgFilePath string) Config {
	var config Config
	file, err := os.OpenFile(cfgFilePath, os.O_RDONLY, 0644)
	if err != nil {
		logger.Log.Error("Wrong config file!", zap.Error(err))
		return config
	}
	defer file.Close()
	buff, err := io.ReadAll(file)
	if err != nil {
		logger.Log.Error("Bad read config file!", zap.Error(err))
		return config
	}
	err = yaml.Unmarshal(buff, &config)
	if err != nil {
		logger.Log.Error("Bad unmarshall config file!", zap.Error(err))
		return config
	}
	return config
}

// создание файла конфигурации с данными переданными через флаги или переменными окружения
func CreateConfig(cfgFilePath string, config Config) {
	//var config Config

	file, err := os.OpenFile(cfgFilePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		logger.Log.Error("Wrong config file!", zap.Error(err))
	}
	defer file.Close()

	config.InFileLog = true

	enc := yaml.NewEncoder(file)
	err = enc.Encode(&config)
	if err != nil {
		logger.Log.Error("Bad unmarshall config file!", zap.Error(err))
	}
	logger.Log.Info("Config file created!", zap.Error(err))
}
