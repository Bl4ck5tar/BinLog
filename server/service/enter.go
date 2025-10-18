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
	WebsiteService
	HotSearchService
	CalendarService
}

var ServiceGroupApp = new(ServiceGroup)