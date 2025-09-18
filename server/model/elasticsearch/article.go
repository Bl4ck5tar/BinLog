package elasticsearch

import "github.com/elastic/go-elasticsearch/v8/typedapi/types"

//Article 文章表
type Article struct {
	CreatedAt	string		`json:"created_at"`	//创建时间
	UpdatedAt	string		`json:"updated_at"`	//更新时间

	Cover		string		`json:"cover"`		//文章封面
	Title 		string		`json:"title"`		//文章标题
	Keyword		string		`json:"keyword"`	//文章标题-关键字
	Category	string		`json:"category"`	//文章类别
	Tags		[]string	`json:"tags"`		//文章标签
	Abstract	string		`json:"abstract"`	//文章简介
	Content		string		`json:"content"`	//文章内容

	Views		int			`json:"views"`		//浏览量
	Comments	int			`json:"comments"`	//评论量
	Likes 		int			`json:"likes"`		//收藏量
}

//ArticleIndex 文章ES索引
func ArticleIndex() string {
	return "article_index"
}

//ArticleMapping 文章 Mapping 字段映射
func ArticleMapping() *types.TypeMapping {
	return &types.TypeMapping{
		Properties: map[string]types.Property{
			"created_at":	types.DateProperty{NullValue: nil, Format: func(s string) *string { return &s}("yyyy-MM-dd HH:mm:ss")},
			"updated_at":	types.DateProperty{NullValue: nil, Format: func(s string) *string { return &s}("yyyy-MM-dd HH:mm:ss")},
			"cover":		types.TextProperty{},			//字符串字段TextProperty，该类型特点：可分词，适合做全文检索
			"title":		types.TextProperty{},
			"keyword":		types.KeywordProperty{},		//精确匹配字段KeywordProperty，特点：不分词，适合做精确匹配或聚合
			"category":		types.KeywordProperty{},
			"tags":			[]types.KeywordProperty{},
			"abstract":		types.TextProperty{},
			"content":		types.TextProperty{},
			"views":		types.IntegerNumberProperty{},	//数值字段IntegerNumberProperty，整数型，适合排序、统计、聚合
			"comments":		types.IntegerNumberProperty{},
			"likes":		types.IntegerNumberProperty{},
		},
	}
}