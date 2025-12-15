package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

const (
	ServiceTypeLabel = "SERVICE_TYPE"
	FileDirLabel     = "CONFIG_DIR"
	FilePostfixLabel = "CONFIG_POSTFIX"
)

const (
	ServiceTypeDev  = "dev"
	ServiceTypeProd = "prod"
)

const (
	FileNameLayout     = "config.%s"
	FileDirDefault     = "configs/"
	FilePostfixDefault = "yaml"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"service"`
	Log      LogConfig      `mapstructure:"log"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Session  SessionConfig  `mapstructure:"session"`
}

type ServerConfig struct {
	// GinBaseConfig
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
	// ServerConfig
	ReadTimeout  int `mapstructure:"read_timeout"`
	WriteTimeout int `mapstructure:"write_timeout"`
	CloseTimeout int `mapstructure:"close_timeout"`

	//TODO: 待添加Gin中间件
	LimitNumber int `mapstructure:"limit_number"`
}

type DatabaseConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DBName       string `mapstructure:"dbname"`
	SSLMode      string `mapstructure:"sslmode"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxLifeTime  int    `mapstructure:"max_life_time"`
	LogLevel     string `mapstructure:"log_level"`
	AutoMigrate  bool   `mapstructure:"auto_migrate"`
}

type RedisConfig struct {
	Addr         string `mapstructure:"addr"`
	Password     string `mapstructure:"password"`
	DB           int    `mapstructure:"db"`
	PoolSize     int    `mapstructure:"pool_size"`
	MinIdleConns int    `mapstructure:"min_idle_conns"`
	MaxRetries   int    `mapstructure:"max_retries"`
	DialTimeout  int    `mapstructure:"dial_timeout"`
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout"`
}

type SessionConfig struct {
	Secret string `mapstructure:"secret"`
}

type LogConfig struct {
	Level    string `mapstructure:"level"`
	FilePath string `mapstructure:"file_path"`
}

var globalConfig *Config

func Load(serviceType string, configPath string, configType string) (*Config, error) {
	// 设置 viper 的文件名、文件类型以及文件地址
	viper.SetConfigName(fmt.Sprintf(FileNameLayout, serviceType))
	viper.SetConfigType(configType)
	viper.AddConfigPath(configPath)

	// 自动扫描环境
	viper.AutomaticEnv()

	// 将配置文件读取到环境
	err := viper.ReadInConfig()

	if err != nil {
		log.Println("Read Config file error: ", err)
		return nil, err
	}

	// 反序列化配置文件
	var cfg *Config

	err = viper.Unmarshal(&cfg)

	if err != nil {
		log.Println("Unmarshal Config file error: ", err)
		return nil, err
	}

	globalConfig = cfg

	return globalConfig, nil
}

func GetConfig() *Config {
	return globalConfig
}
