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
	CommentService
}

var ServiceGroupApp = new(ServiceGroup)