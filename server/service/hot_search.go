package service

import (
	"BinLog/server/global"
	"BinLog/server/model/other"
	"BinLog/server/utils/hotSearch"
	"encoding/json"
	"time"
)

type HotSearchService struct {

}

func (hotSearchService *HotSearchService) GetHotSearchDataBySource(sourceStr string) (other.HotSearchData, error) {
	//查询热榜缓存，如果命中则反序列化后直接返回；未命中则进入回溯逻辑
	result, err := global.Redis.Get(sourceStr).Result()
	if err != nil {
		source := hotSearch.NewSource(sourceStr)
		//抓取热榜数据
		hotSearchData, err := source.GetHotSearchData(30)
		if err != nil {
			return other.HotSearchData{}, err
		}
		bytes, err := json.Marshal(hotSearchData)
		if err != nil {
			return other.HotSearchData{}, err
		}

		//设置缓存有效期为1小时
		if err := global.Redis.Set(sourceStr, bytes, time.Hour).Err(); err != nil {
			return other.HotSearchData{}, err
		}
		return hotSearchData, nil
	}

	var hotSearchData other.HotSearchData
	if err := json.Unmarshal([]byte(result), &hotSearchData); err != nil {
		return other.HotSearchData{}, err
	}
	return hotSearchData, nil
}