package database

import "BinLog/server/global"

//ArticleLike 文章收藏表
type ArticleLike struct {
	global.MODEL
	ArticleID	string	`json:"article_id"`		//文章ID
	UserID		uint	`json:"user_id"`		//用户ID
	User		User	`json:"-" gorm:"foreignKey:UserID"`
}