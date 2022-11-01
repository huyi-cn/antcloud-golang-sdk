package sdk

import (
	"time"
)

// SDK使用配置
type Config struct {
	AutoRetry     bool          // 是否重试
	MaxRetryTimes int           // 重试次数
	Debug         bool          // 是否输出DEBUG信息
	Timeout       time.Duration // 超时时间
}

func NewConfig() (config *Config) {
	config = &Config{
		AutoRetry:     true,
		MaxRetryTimes: 3,
		Debug:         false,
		Timeout:       10000000000,
	}
	return
}
