package config

// ES ElasticSearch 配置
type ES struct {
	URL				string 	`json:"url" yaml:"url"`								//ES服务的URL
	Username 		string	`json:"username" yaml:"username"`					//用于连接ES的用户名
	Password 		string	`json:"password" yaml:"pasword"`					//用于连接ES的密码
	IsConsolePrint	bool 	`json:"is_console_print" yaml:"is_console_print"`	//是否在控制台打印ES语句，true表示打印，false表示不打印
}