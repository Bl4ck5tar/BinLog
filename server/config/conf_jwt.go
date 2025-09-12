package config

type Jwt struct {
	AccessTokenSecret 	    string	`json:"access_token_secret" yaml:"access_token_secret"`				//用于生成和验证访问令牌的密钥
	RefreshTokenSecret	    string 	`json:"refresh_token_secret" yaml:"refresh_token_secret"`			//用于生成和验证刷新令牌的密钥
	AccessTokenExpiryTime   string	`json:"access_token_expiry_time" yaml:"access_token_expiry_time"`	//访问令牌的过期时间，例如“15m”表示15分钟
	RefreshTokenExpiryTime	string	`json:"refresh_token_expiry_time" yaml:"refresh_token_expiry_time"`	//刷新令牌的过期时间
	Issuer					string 	`json:"issuer" yaml:"issuer"`										//JWT的签发者信息，通常是应用或服务的名称
}