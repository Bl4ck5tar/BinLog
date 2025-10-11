package database

import (
	"BinLog/server/global"
	"BinLog/server/model/elasticsearch"
	"context"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/scriptlanguage"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
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

//AfterCreate 钩子，创建后调用
func (c *Comment) AfterCreate(_ *gorm.DB) error {
	source := "ctx._source.comments += 1"
	script := types.Script{Source: &source, Lang: &scriptlanguage.Painless}
	_, err := global.ESClient.Update(elasticsearch.ArticleIndex(), c.ArticleID).Script(&script).Do(context.TODO())
	return err
}

//AfterDelete 钩子，删除后调用
func (c *Comment) BeforeDelete(_ *gorm.DB) error {
	source := "ctx._source.comments -= 1"
	script := types.Script{Source: &source, Lang: &scriptlanguage.Painless}
	_, err := global.ESClient.Update(elasticsearch.ArticleIndex(), c.ArticleID).Script(&script).Do(context.TODO())
	return err
}