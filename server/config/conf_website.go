package config

type Website struct {
	Logo					string	`json:"logo" yaml:"logo"`										//网站logo
	FullLogo				string	`json:"full_logo" yaml:"full_logo"`								//完整logo
	Title					string	`json:"title" yaml:"title"`										//网站标题
	Slogan					string	`json:"slogan" yaml:"slogan"`									//网站标语
	SloganEn				string	`json:"slogan_en" yaml:"slogan_en"`								//英文标语
	Description				string	`json:"description" yaml:"description"`							//网站描述
	Version 				string	`json:"version" yaml:"version"`									//网站版本
	CreateAt				string 	`json:"create_at" yaml:"create_at"`								//创建时间
	IcpFiling				string	`json:"icp_filling" yaml:"icp_filling"`							//ICP 备案
	PublicSecurityFilling	string	`json:"public_security_filling" yaml:"public_security_filling"`	//公安备案
	BilibiliURL				string	`json:"bilibili_url" yaml:"bilibili_url"`						//Bilibili链接
	GiteeURL				string	`json:"gitee_url" yaml:"gitee_url"`								//Gitee链接
	GithubURL				string	`json:"github_url" yaml:"github_url"`							//Github链接
	Name					string	`json:"name" yaml:"name"`										//昵称
	Job						string	`json:"job" yaml:"job"`											//职业
	Address					string	`json:"address" yaml:"address"`									//地址
	Email					string	`json:"email" yaml:"email"`										//邮箱
	QQImage					string	`json:"qq_image" yaml:"qq_image"`								//QQ图片链接
	WechatImage				string	`json:"wechat_image" yaml:"wechat_image"`						//微信图片链接
}