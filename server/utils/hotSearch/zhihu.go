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

type Zhihu struct {
	
}

// func (*Zhihu) GetHotSearchData(maxNum int) (other.HotSearchData, error) {
// 	client := &http.Client{Timeout: 5*time.Second}
// 	req, err := http.NewRequest("GET", "https://www.zhihu.com/billboard", nil)
// 	if err != nil {
// 		return other.HotSearchData{}, err
// 	}

// 	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return  other.HotSearchData{}, err
// 	}

// 	defer resp.Body.Close()

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return other.HotSearchData{}, err
// 	}

// 	var jsonStr string
// 	reg := regexp.MustCompile(`(?s)<script id="js-initialData" type="text/json">(.*?)</script>`)
// 	result := reg.FindAllStringSubmatch(string(body), -1)
// 	if len(result) > 0 && len(result[0]) > 1 {
// 		jsonStr = result[0][1]
// 	}else {
// 		return other.HotSearchData{}, errors.New("failed to get data")
// 	}

// 	hotList := make([]other.HotItem, 0, maxNum)
// 	for i:= 0; i<maxNum; i++ {
// 		if index := gjson.Get(jsonStr, "initialState.topstory.hotList."+strconv.Itoa(i)+".id"); !index.Exists() {
// 			break
// 		}

// 		hotList = append(hotList, other.HotItem{
// 			Index:		 i+1,
// 			Title:		 gjson.Get(jsonStr, "initialState.topstory.hotList."+strconv.Itoa(i)+".target.titleArea.text").Str,
// 			Description: gjson.Get(jsonStr, "initialState.topstory.hotList."+strconv.Itoa(i)+".target.excerptArea.text").Str,
// 			Image: 		 gjson.Get(jsonStr, "initialState.topstory.hotList."+strconv.Itoa(i)+".target.imageArea.url").Str,
// 			Popularity:  gjson.Get(jsonStr, "initialState.topstory.hotList."+strconv.Itoa(i)+".target.metricsArea.text").Str,
// 			URL:		 gjson.Get(jsonStr, "initialState.topstory.hotList."+strconv.Itoa(i)+".target.link.url").Str,
// 		})
// 	}

// 	return other.HotSearchData{
// 		Source: 	"知乎热榜",
// 		UpdateTime: time.Now().Format("2006-01-02 15:04:05"),
// 		HotList: 	hotList,
// 	}, nil
// }

func (*Zhihu) GetHotSearchData(maxNum int) (other.HotSearchData, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", "https://www.zhihu.com/billboard", nil)
	if err != nil {
		return other.HotSearchData{}, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		return other.HotSearchData{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return other.HotSearchData{}, err
	}

	// ✅ 正则匹配 application/json
	reg := regexp.MustCompile(`(?s)<script id="js-initialData" type="application/json">(.*?)</script>`)
	result := reg.FindStringSubmatch(string(body))
	if len(result) < 2 {
		return other.HotSearchData{}, errors.New("failed to extract initialData JSON")
	}

	jsonStr := result[1]

	// ✅ 构造热榜列表
	hotList := make([]other.HotItem, 0, maxNum)
	for i := 0; i < maxNum; i++ {
		if !gjson.Get(jsonStr, "initialState.topstory.hotList."+strconv.Itoa(i)+".id").Exists() {
			break
		}
		hotList = append(hotList, other.HotItem{
			Index:       i + 1,
			Title:       gjson.Get(jsonStr, "initialState.topstory.hotList."+strconv.Itoa(i)+".target.titleArea.text").Str,
			Description: gjson.Get(jsonStr, "initialState.topstory.hotList."+strconv.Itoa(i)+".target.excerptArea.text").Str,
			Image:       gjson.Get(jsonStr, "initialState.topstory.hotList."+strconv.Itoa(i)+".target.imageArea.url").Str,
			Popularity:  gjson.Get(jsonStr, "initialState.topstory.hotList."+strconv.Itoa(i)+".target.metricsArea.text").Str,
			URL:         gjson.Get(jsonStr, "initialState.topstory.hotList."+strconv.Itoa(i)+".target.link.url").Str,
		})
	}

	return other.HotSearchData{
		Source:     "知乎热榜",
		UpdateTime: time.Now().Format("2006-01-02 15:04:05"),
		HotList:    hotList,
	}, nil
}
