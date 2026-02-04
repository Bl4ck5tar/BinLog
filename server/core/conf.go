package core

import (
	"gopkg.in/yaml.v3"
	"BinLog/server/config"
	"log"
	"BinLog/server/utils"
	"os"
)

//InitConf 从YAML 文件加载配置
func InitConf() *config.Config {
	c := &config.Config{}
	yamlConf, err := utils.LoadYAML()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	err = yaml.Unmarshal(yamlConf, c)
	if err != nil {
		log.Fatalf("Failed to unmarshal YAML configuration: %v", err)
	}
	// Allow local environment overrides for sensitive fields.
	if secret := os.Getenv("BINLOG_EMAIL_SECRET"); secret != "" {
		c.Email.Secret = secret
	}
	return c
}
