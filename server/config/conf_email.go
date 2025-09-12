package config

//Email邮箱配置， 可登录QQ邮箱，https://mail.qq.com
type Email struct {
	Host	  string	`json:"host" yaml:"host"`			//邮件服务器地址，例如smtp.qq.com
	Port	  int 		`json:"port" yaml:"port"`			//邮件服务器端口，常见的如587（TLS）或465（SSL）
	From 	  string 	`json:"from" yaml:"from"`			//发件人邮箱地址
	Nickname  string	`json:"nickname" yaml:"nickname"`	//发件人昵称，用于显示在邮件中发件人信息
	Secret 	  string	`json:"secret" yaml:"secret"`		//发件人邮箱的密码或应用专用密码，用于身份验证
	IsSSL 	  bool		`json:"is_ssl" yaml:"is_ssl"`		//是否使用SSL加密连接，true表示使用，false表示不使用
}