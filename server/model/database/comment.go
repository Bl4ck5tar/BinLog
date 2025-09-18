package database

import (
	"BinLog/server/global"

	"github.com/gofrs/uuid"
)

//Comment 评论表
type Comment struct {
	global.MODEL
	ArticleID 	string		`json:"article_id"`						//文章ID
	PID 		*uint		`json:"p_id"`							//父评论ID
	PComment	*Comment	`json:"-" gorm:"foreignKey:PID"`		
	Children	[]Comment	`json:"children" gorm:"foreignKey:PID"`	//子评论
	UserUUID	uuid.UUID	`json:"user_uuid" gorm:"type:char(36)"`	//关联的用户
	Content		string		`json:"content"`						//内容
}