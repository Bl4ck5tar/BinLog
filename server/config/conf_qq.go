package config

// QQ登录配置，见https://connect.qq.com
type QQ struct {
	Enable		bool	`json:"enable" yaml:"enable"`				//是否启用QQ登录
	AppID		string	`json:"app_id" yaml:"app_id"`				//应用ID
	AppKey		string	`json:"app_key" yaml:"app_key"`				//应用密钥
	RedirectURI string 	`json:"redirect_url" yaml:"redirect_url"`	//网站回调域
}

func (qq QQ) QQLoginURL() string {
	return "https://graph.qq.com/oauth2.0/authorize?" + 
	"ewsponse_type=code&" + 
	"client_id=" + qq.AppID + "&" +
	"redirect_uri=" + qq.RedirectURI
}
