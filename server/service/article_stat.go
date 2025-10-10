package service

import (
	"BinLog/server/global"
	"strconv"
)
//文章浏览量、点赞、收藏数统计
type CountDB struct {
	Index string
}

func (article *ArticleService) NewArticleView() CountDB {
	return CountDB{
		//Redis 哈希表的key
		Index: 	"article_views",
	}
}

//Set 在原有基础上加一
func (c CountDB) Set(id string) error {
	// num, _ := global.Redis.HGet(c.Index, id).Int()
	// num++
	// err := global.Redis.HSet(c.Index, id, num).Err()
	// return err
	//使用以下可避免并发冲突，同时性能更好
	return global.Redis.HIncrBy(c.Index, id, 1).Err()
}

//GetInfo 取出数据
func (c CountDB) GetInfo() map[string]int {
	var Info = map[string]int{}
	maps := global.Redis.HGetAll(c.Index).Val()
	for id, val := range maps {
		num, _ := strconv.Atoi(val)
		Info[id] = num
	}
	return Info
}

//Clear 清除数据
func (c CountDB) Clear() {
	global.Redis.Del(c.Index)
}