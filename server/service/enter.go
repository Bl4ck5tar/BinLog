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
	AdvertisementService
	FriendLinkService
	FeedbackService
}

var ServiceGroupApp = new(ServiceGroup)