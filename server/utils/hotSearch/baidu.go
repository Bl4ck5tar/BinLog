package hotSearch

import (
	"BinLog/server/model/other"
	"errors"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/tidwall/gjson"
)

type Baidu struct {

}
func (*Baidu) GetHotSearchData(maxNum int) (HotSearchData other.HotSearchData, err error) {
	client := &http.Client{Timeout: 5*time.Second}	//添加超时控制，防止阻塞
	req, _ := http.NewRequest("GET", "https://top.baidu.com/board?tab=realtime", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; HotSearchBot/1.0; +https://example.com)")
	resp, err := client.Do(req)
	if err != nil {
		return other.HotSearchData{}, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return other.HotSearchData{}, err
	}

	var jsonStr string
	reg := regexp.MustCompile(`<!--s-data:({.*?})-->`)
	result := reg.FindAllStringSubmatch(string(body), -1)
	if len(result) > 0 && len(result[0]) > 1 {
		jsonStr = result[0][1]
	}else {
		return other.HotSearchData{}, errors.New("failed to get data")
	}

	updateTime := time.Unix(gjson.Get(jsonStr, "data.cards.0.updateTime").Int(), 0).Format("2006-01-02 15:04:05")

	var hotList []other.HotItem
	for i := 0; i < maxNum; i++ {
		if index := gjson.Get(jsonStr, "data.cards.0.content."+strconv.Itoa(i)+".index"); !index.Exists() {
			break
		}
		hotList = append(hotList, other.HotItem{
			Index:	 i+1,
			Title:       gjson.Get(jsonStr, "data.cards.0.content."+strconv.Itoa(i)+".word").Str,
			Description: gjson.Get(jsonStr, "data.cards.0.content."+strconv.Itoa(i)+".desc").Str,
			Image: 		 gjson.Get(jsonStr, "data.cards.0.content."+strconv.Itoa(i)+".img").Str,
			Popularity:  gjson.Get(jsonStr, "data.cards.0.content."+strconv.Itoa(i)+".hotScore").Str,
			URL: 		 gjson.Get(jsonStr, "data.cards.0.content."+strconv.Itoa(i)+".rawUrl").Str,
		})

	}

	//待测试
	// reg := regexp.MustCompile(`<!--s-data:(\{[\s\S]*?\})-->`)
	// result := reg.FindStringSubmatch(string(body))
	// if len(result) < 2 {
	// 	return other.HotSearchData{}, errors.New("failed to extract JSON data from HTML")
	// }
	// jsonStr := result[1]
	// updateTime := time.Unix(gjson.Get(jsonStr, "data.cards.0.updateTime").Int(), 0).Format("2006-01-02 15:04:05")

	// content := gjson.Get(jsonStr, "data.cards.0.content")
	// hotList := make([]other.HotItem, 0, maxNum)

	// for i , item := range content.Array() {
	// 	if i >= maxNum {
	// 		break
	// 	}
	// 	hotList = append(hotList, other.HotItem{
	// 		Index: 			i+1,
	// 		Title:			item.Get("word").Str,
	// 		Description:	item.Get("desc").Str,
	// 		Image:		 	item.Get("img").Str,
	// 		Popularity:	 	item.Get("hotScore").Str,
	// 		URL:	 		item.Get("rawUrl").Str,
	// 	})
	// }
	return other.HotSearchData{
		Source: "百度热搜", 
		UpdateTime: updateTime, 
		HotList: hotList,
		}, nil
}