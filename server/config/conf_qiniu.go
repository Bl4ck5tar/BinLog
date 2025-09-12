package config

type Qiniu struct {
	Zone 		  string 	`json:"zone" yaml:"zone"`						//存储区域
	Bucket		  string	`json:"bucket" yaml:"bucket"`					//空间名称
	ImgPath 	  string 	`json:"img_path" yaml:"img_path"`				//CDN加速域名
	AccessKey 	  string	`json:"acess_key" yaml:"access_key"`			//密钥AK
	SecretKey 	  string	`json:"secret_key" yaml:"secret_key"`			//密钥SK
	UseHTTPS	  bool		`json:"use_https" yaml:"use_https"`				//是否使用https
	UseCdnDomains bool		`json:"use_cdn_domains" yaml:"use_cdn_domains"`	//上传是否使用CDN上传加速
}