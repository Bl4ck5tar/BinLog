package service

type ServiceGroup struct {
	BaseService
	EsService
	JwtService
	GaodeService
	UserService
	QQService
	ImageService
	ArticleService
}

var ServiceGroupApp = new(ServiceGroup)