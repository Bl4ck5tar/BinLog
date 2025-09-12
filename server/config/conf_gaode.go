package config

//高德服务配置详情：https://lbs.amap.com
type Gaode struct {
	Enable bool		`json:"enable" yaml:"enable"`	//是否开启高德服务
	Key	   string 	`json:"key" yaml:"key"`			//高德服务的应用密钥，用于身份验证和服务访问
}