package task

import (
	"BinLog/server/global"
	"BinLog/server/model/elasticsearch"
	"BinLog/server/service"
	"context"
	"strconv"
	"time"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/scriptlanguage"
	"go.uber.org/zap"
)

// func UpdateArticleViewsSyncTask() error {
// 	// 1. 从 Redis 获取文章浏览量缓存（map[article_id]view_count）
// 	articleView := service.ServiceGroupApp.ArticleService.NewArticleView()
// 	viewsInfo := articleView.GetInfo()

// 	for id, num := range viewsInfo {
// 		// 2. 跳过无变化的记录
// 		if num == 0 {
// 			continue
// 		}

// 		// 3. 构造 Elasticsearch 更新脚本
// 		source := "ctx._source.views += " + strconv.Itoa(num)
// 		script := types.Script{Source: &source, Lang: &scriptlanguage.Painless}

// 		// 4. 执行更新
// 		_, err := global.ESClient.Update(elasticsearch.ArticleIndex(), id).
// 			Script(&script).
// 			Do(context.TODO())
		//更新第一篇文章后就返回，无法更新其他数据
// 		return err
// 	}

// 	// 5. 清空 Redis 缓存
	//如果中途ES更新失败，Redis中的缓存也会被清空，导致浏览量丢失
// 	articleView.Clear()
// 	return nil
// }

//改进版
func UpdateArticleViewsSyncTask() error {
	//防止长时间阻塞或卡死任务
	ctx, cancel := context.WithTimeout(context.Background(), 30 * time.Second)
	defer cancel()

	articleView := service.ServiceGroupApp.ArticleService.NewArticleView()
	viewsInfo := articleView.GetInfo()
	if len(viewsInfo) == 0 {
		return nil
	}

	for id, num := range viewsInfo {
		if num == 0 {
			continue
		}

		source := "ctx._source.views +=" + strconv.Itoa(num)
		script := types.Script{Source: &source, Lang: &scriptlanguage.Painless}

		_, err := global.ESClient.Update(elasticsearch.ArticleIndex(), id).Script(&script).Do(ctx)
		if err != nil {
			global.Log.Error("Failed to sync article view", zap.String("id", id), zap.Error(err))
			continue	//不中断整个任务
		}
	}
	//仅在全部完成后再清空 Redis
	articleView.Clear()
	return nil
}
