package service

import (
	"BinLog/server/global"
	"BinLog/server/model/database"
	"BinLog/server/model/elasticsearch"
	"BinLog/server/utils"
	"context"
	"encoding/json"
	"errors"

	"github.com/elastic/go-elasticsearch/v8/typedapi/core/bulk"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/update"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/refresh"
	"gorm.io/gorm"
)

//Create 用于将文章创建到 Elasticsearch
func (articleService *ArticleService) Create (a *elasticsearch.Article) error {
	//将文章索引到 Elasticsearch 中，并设置刷新操作为 true
	_, err := global.ESClient.Index(elasticsearch.ArticleIndex()).Request(a).Refresh(refresh.True).Do(context.TODO())
	return err
}

//Delete 用于删除 Elasticsearch 中的文章
func (articleService *ArticleService) Delete(ids []string) error {
	var request bulk.Request

	//遍历文章，构建批量删除请求
	for _, id := range ids {
		request = append(request, types.OperationContainer{Delete: &types.DeleteOperation{Id_: &id}})
	}

	//执行批量删除请求，并设置刷新操作为 true
	_, err := global.ESClient.Bulk().			//进入 ES 的批量操作系统
	Request(&request).							//绑定刚组装好的删除请求
	Index(elasticsearch.ArticleIndex()).		//指定要删除的索引
	Refresh(refresh.True).						//删除操作完成后立即刷新，使删除结果立即可见
	Do(context.TODO())							//执行请求
	return err
}

//Get 用于通过 ID 从 Elasticsearch 获取文章
func (articleService *ArticleService) Get(id string) (elasticsearch.Article, error) {
	var a elasticsearch.Article
	//从 Elasticsearch 获取文章
	res, err := global.ESClient.Get(elasticsearch.ArticleIndex(), id).Do(context.TODO())
	if err != nil {
		return elasticsearch.Article{}, err
	}
	//如果找不到文档，则返回错误
	if !res.Found {
		return elasticsearch.Article{}, errors.New("document not found")
	}
	//将返回的源数据反序列化为 Article 对象
	err = json.Unmarshal(res.Source_, &a)
	return a, err
}

//Update 用于更新文章数据
func (articleService *ArticleService) Update(articleID string, v any) error {
	//将待更新的值转换为 JSON
	bytes, err := json.Marshal(v)
	if err != nil {
		return err
	}
	//执行更新请求，并设置刷新操作为 true
	_, err = global.ESClient.Update(elasticsearch.ArticleIndex(), articleID).Request(&update.Request{Doc: bytes}).Refresh(refresh.True).Do(context.TODO())
	return err
}

func (articleService *ArticleService) Exists(title string) (bool, error) {
	//创建查询请求，匹配标题字段
	req := &search.Request{
		Query: &types.Query{
			Match: map[string]types.MatchQuery{"keyword": {Query: title}},
		},
	}

	//执行搜索查询，查找是否存在该标题的文章
	res, err := global.ESClient.Search().Index(elasticsearch.ArticleIndex()).Request(req).Size(1).Do(context.TODO())
	if err != nil {
		return false, err
	}
	//如果存在该标题，返回true
	return res.Hits.Total.Value > 0, nil
}

//UpdateCategoryCount 更新文章类别的计数（增加或减少）
func (articleService *ArticleService) UpdateCategoryCount(tx *gorm.DB, oldCategory, newCategory string) error {
	//如果新旧类别相同，直接返回，不进行更新
	if newCategory == oldCategory {
		return nil
	}

	//如果新类别不为空，更新新类别的文章技数
	if newCategory != "" {
		var newArticleCategory database.ArticleCategory
		//如果新类别不存在，则创建新类别并设置技数为1
		// if errors.Is(tx.Where("category = ?", newCategory).First(&newArticleCategory).Error, gorm.ErrRecordNotFound) {
		// 	if err := tx.Create(&database.ArticleCategory{Category: newCategory, Number: 1}).Error; err != nil {
		// 		return err
		// 	}
		// }else {
		// 	//如果类别已存在，更新该类别的计数
		// 	if err := tx.Model(&newArticleCategory).Update("number", gorm.Expr("number + ?"), 1).Error; err != nil {
		// 		return err
		// 	}
		// }
		//如果没有则创建，默认 number=0
		if err := tx.FirstOrCreate(&newArticleCategory,database.ArticleCategory{Category: newCategory},).Error; err != nil {
			return err
		}

		//number + 1（原子操作）
		if err := tx.Model(&newArticleCategory).UpdateColumn("number", gorm.Expr("number + ?", 1)).Error; err != nil {
			return err
		}
	}

	//如果旧类别不为空，更新旧类别的文章计数
	if oldCategory != "" {
		var oldArticleCategory database.ArticleCategory
		//更新旧类别的文章计数，减少1
		// if err := tx.Where("category = ?", oldCategory).First(&oldArticleCategory).Update("number", gorm.Expr("number - ?", 1)).Error; err != nil {
		// 	return err
		// }
		// //如果旧类别的计数为1（减少之前），则删除该类别
		// if oldArticleCategory.Number == 1 {
		// 	if err := tx.Delete(&oldArticleCategory).Error; err != nil {
		// 		return err
		// 	}
		// }
		//先找到记录
		if err := tx.Where("category = ?", oldCategory).First(&oldArticleCategory).Error; err != nil {
			//如果没有旧分类，直接跳过
			if errors.Is(err, gorm.ErrRecordNotFound){
				return nil
			}
			return err
		}

		//number - 1（原子操作）
		if err := tx.Model(&oldArticleCategory).UpdateColumn("number", gorm.Expr("number - ?", 1)).Error; err != nil {
			return err
		}

		//检查是否变为0，如果是0就删除类别
		if oldArticleCategory.Number == 1 {
			if err := tx.Delete(&oldArticleCategory).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func (articleService *ArticleService) UpdateTagsCount(tx *gorm.DB, oldTags, newTags []string) error {
	// 比较旧标签和新标签，获取新增和移除的标签
	addedTags, removedTags := utils.DiffArrays(oldTags, newTags)

	// 处理新增的标签
	for _, addedTag := range addedTags {
		var t database.ArticleTag
		// 如果标签不存在，则创建该标签并设置计数为1
		// if errors.Is(tx.Where("tag = ?", addedTag).First(&t).Error, gorm.ErrRecordNotFound) {
		// 	if err := tx.Create(&database.ArticleTag{Tag: addedTag, Number: 1}).Error; err != nil {
		// 		return err
		// 	}
		// } else {
		// 	// 如果标签已存在，更新标签的计数
		// 	if err := tx.Model(&t).Update("number", gorm.Expr("number + ?", 1)).Error; err != nil {
		// 		return err
		// 	}
		// }

		//如果没有标签就创建一个，number=0
		if err := tx.FirstOrCreate(&t, database.ArticleTag{Tag: addedTag}).Error; err != nil {
			return err
		}

		//number + 1
		if err := tx.Model(&t).UpdateColumn("number", gorm.Expr("number + ?", 1)).Error; err != nil {
			return err
		}
	}

	// 处理移除的标签
	for _, removedTag := range removedTags {
		var t database.ArticleTag
		// // 更新标签计数，减少 1
		// if err := tx.Where("tag = ?", removedTag).First(&t).Update("number", gorm.Expr("number - ?", 1)).Error; err != nil {
		// 	return err
		// }
		// // 如果标签的计数为 1（减少 1 之前），则删除该标签
		// if t.Number == 1 {
		// 	if err := tx.Delete(&t).Error; err != nil {
		// 		return err
		// 	}
		// }

		if err := tx.Where("tag = ?", removedTag).First(&t).Error; err != nil {
			//标签不存在则跳过
			if errors.Is(err, gorm.ErrRecordNotFound) {
				continue
			}
			return err
		}

		if err := tx.Model(&t).UpdateColumn("number", gorm.Expr("number - ?", 1)).Error; err != nil {
			return err
		}

		//原来的值为1，更新后即为0，需要删除
		if t.Number == 1 {
			if err := tx.Delete(&t).Error; err != nil {
				return err
			}
		}
	}
	return nil
}