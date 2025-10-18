package hotSearch

import (
	"BinLog/server/model/other"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/tidwall/gjson"
)

type Toutiao struct {

}

func (*Toutiao) GetHotSearchData(maxNum int) (other.HotSearchData, error) {
	client := &http.Client{Timeout: 5*time.Second}
	req, err := http.NewRequest("GET", "https://www.toutiao.com/hot-event/hot-board/?origin=toutiao_pc", nil)
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

	jsonStr := string(body)

	time := gjson.Get(jsonStr, "impr_id").Str[:14]
	updateTime := time[:4] + "-" + time[4:6] + "-" + time[6:8] + " " + time[8:10] + ":" + time[10:12] + ":" + time[12:14]

	hotList := make([]other.HotItem, 0, maxNum)

	for i := 0 ; i<maxNum; i++ {
		if index := gjson.Get(jsonStr, "data."+strconv.Itoa(i)+".ClusterId"); !index.Exists(){
			break
		}

		hotList = append(hotList, other.HotItem{
			Index: 		 i+1,
			Title:		 gjson.Get(jsonStr, "data."+strconv.Itoa(i)+".Title").Str,
			Description: "",
			Image:		 gjson.Get(jsonStr, "data."+strconv.Itoa(i)+".Image.url").Str,
			Popularity:  gjson.Get(jsonStr, "data."+strconv.Itoa(i)+".HotValue").Str,
			URL:		 gjson.Get(jsonStr, "data."+strconv.Itoa(i)+".Url").Str,
		})
	}

	return other.HotSearchData{
		Source: 	"头条热榜",
		UpdateTime: updateTime,
		HotList: 	hotList,
	}, nil
}

// func (*Toutiao) GetHotSearchData(maxNum int) (other.HotSearchData, error) {
// 	client := &http.Client{Timeout: 5 * time.Second}
// 	req, err := http.NewRequest("GET", "https://www.toutiao.com/hot-event/hot-board/?origin=toutiao_pc", nil)
// 	if err != nil {
// 		return other.HotSearchData{}, err
// 	}

// 	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return other.HotSearchData{}, err
// 	}
// 	defer resp.Body.Close()

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return other.HotSearchData{}, err
// 	}

// 	jsonStr := string(body)
// 	imprID := gjson.Get(jsonStr, "impr_id").Str
// 	var updateTime string
// 	if len(imprID) >= 14 {
// 		updateTime = fmt.Sprintf("%s-%s-%s %s:%s:%s",
// 			imprID[0:4], imprID[4:6], imprID[6:8],
// 			imprID[8:10], imprID[10:12], imprID[12:14])
// 	} else {
// 		updateTime = time.Now().Format("2006-01-02 15:04:05")
// 	}

// 	dataArray := gjson.Get(jsonStr, "data").Array()
// 	hotList := make([]other.HotItem, 0, len(dataArray))

// 	for i, item := range dataArray {
// 		if i >= maxNum {
// 			break
// 		}
// 		hotList = append(hotList, other.HotItem{
// 			Index:       i + 1,
// 			Title:       item.Get("Title").Str,
// 			Description: "",
// 			Image:       item.Get("Image.url").Str,
// 			Popularity:  item.Get("HotValue").Str,
// 			URL:         item.Get("Url").Str,
// 		})
// 	}

// 	return other.HotSearchData{
// 		Source:     "头条热榜",
// 		UpdateTime: updateTime,
// 		HotList:    hotList,
// 	}, nil
// }
