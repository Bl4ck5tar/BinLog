package config

//Redis缓存数据库配置
type Redis struct {
	Address 	string	`json:"address" yaml:"address"`		//Redis服务器地址
	Password	string	`json:"password" yaml:"password"`	//连接Redis的密码，没有则留空
	DB			int 	`json:"db" yaml:"db"`				//指定使用的数据库索引，单实例模式下可选择的数据库，默认值为0
}