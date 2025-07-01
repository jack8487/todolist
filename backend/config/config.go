package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

/*
mapstructure它是一个结构体标签（struct tag），用于指定配置文件中的键名如何映射到 Go 结构体的字段
当 viper 读取配置文件时，会使用这些标签来进行自动映射
*/
// Config 配置结构体
type Config struct {
	Server ServerConfig `mapstructure:"server"`
	MySQL  MySQLConfig  `mapstructure:"mysql"`
	Redis  RedisConfig  `mapstructure:"redis"`
	JWT    JWTConfig    `mapstructure:"jwt"`
	Log    LogConfig    `mapstructure:"log"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

// MySQLConfig MySQL配置
type MySQLConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	Username        string        `mapstructure:"username"`
	Password        string        `mapstructure:"password"`
	Database        string        `mapstructure:"database"`
	Charset         string        `mapstructure:"charset"`
	ParseTime       bool          `mapstructure:"parse_time"`
	Loc             string        `mapstructure:"loc"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

// DSN 返回MySQL连接字符串
func (c MySQLConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%v&loc=%s",
		c.Username, c.Password, c.Host, c.Port, c.Database,
		c.Charset, c.ParseTime, c.Loc)
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	Password        string        `mapstructure:"password"`
	DB              int           `mapstructure:"db"`
	PoolSize        int           `mapstructure:"pool_size"`
	MinIdleConns    int           `mapstructure:"min_idle_conns"`
	MaxConnLifetime time.Duration `mapstructure:"max_conn_lifetime"`
}

// Addr 返回Redis地址
func (c RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// JWTConfig JWT配置
type JWTConfig struct {
	SecretKey   string        `mapstructure:"secret_key"`
	ExpireHours time.Duration `mapstructure:"expire_hours"`
	Issuer      string        `mapstructure:"issuer"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
	Compress   bool   `mapstructure:"compress"`
}

var GlobalConfig Config

// LoadConfig 加载配置
func LoadConfig(configFile string) error {
	viper.SetConfigFile(configFile)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := viper.Unmarshal(&GlobalConfig); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 转换时间单位
	GlobalConfig.MySQL.ConnMaxLifetime *= time.Second
	GlobalConfig.Redis.MaxConnLifetime *= time.Second
	GlobalConfig.JWT.ExpireHours *= time.Hour

	return nil
}
