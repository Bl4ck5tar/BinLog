package hotSearch

import (
	"BinLog/server/model/other"
	"errors"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

type Kuaishou struct {

}

func (*Kuaishou) GetHotSearchData(maxNum int) (HotSearchData other.HotSearchData, err error) {
	client := &http.Client{Timeout: 5*time.Second}
	req, err := http.NewRequest("GET", "https://www.kuaishou.com/?isHome=1", nil)
	if err != nil {
		return other.HotSearchData{}, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36 Edg/131.0.0.0")
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
	reg := regexp.MustCompile(`window.__APOLLO_STATE__=({.*?});`)
	result := reg.FindAllStringSubmatch(string(body), -1)
	if len(result) > 0 && len(result[0]) > 1 {
		jsonStr = result[0][1]
	}else {
		return other.HotSearchData{}, errors.New("failed to get data")
	}

	updateTime := time.Now().Format("2006-01-02 15:04:05")
	hotList := make([]other.HotItem, 0, maxNum)
	for i := range maxNum {
		index := gjson.Get(jsonStr, `defaultClient.$ROOT_QUERY\.visionHotRank({\"page\":\"home\"}).items.`+strconv.Itoa(i)+".id")
		if !index.Exists() {
			break
		}
		result := escapeSpecialCharacters(index.Str)
		hotList = append(hotList, other.HotItem{
			Index:       int(gjson.Get(jsonStr, "defaultClient."+result+".rank").Int() + 1),
			Title: 		 gjson.Get(jsonStr, "defaultClient."+result+".name").Str,
			Description: "",
			Image: 		 gjson.Get(jsonStr, "defaultClient."+result+".poster").Str,
			Popularity:  gjson.Get(jsonStr, "defaultClient."+result+".hotValue").Str,
			URL:         "https://www.kuaishou.com/short-video/" + 
			gjson.Get(jsonStr, "defaultClient."+result+".photoIds.json.0").Str + 
			"?streamSource=hotrank&trendingId=" + 
			gjson.Get(jsonStr, "defaultClient."+result+".id").Str + 
			"&area=homexxunknown",
		})
		
	}
	return other.HotSearchData{
		Source: "快手热榜",
		UpdateTime: updateTime,
		HotList: hotList,
	}, nil
}

//待测版本
// func (*Kuaishou) GetHotsearchData(maxNum int) (other.HotSearchData, error) {
// 	client := &http.Client{Timeout: 5 * time.Second}
// 	req, err := http.NewRequest("GET", "https://www.kuaishou.com/?isHome=1", nil)
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

// 	reg := regexp.MustCompile(`window.__APOLLO_STATE__=({[\s\S]*?});`)
// 	match := reg.FindStringSubmatch(string(body))
// 	if len(match) < 2 {
// 		return other.HotSearchData{}, errors.New("failed to extract JSON from HTML")
// 	}
// 	jsonStr := match[1]

// 	hotList := make([]other.HotItem, 0, maxNum)
// 	items := gjson.Get(jsonStr, `defaultClient.$ROOT_QUERY\.visionHotRank({\"page\":\"home\"}).items`)

// 	for i, item := range items.Array() {
// 		if i >= maxNum {
// 			break
// 		}
// 		itemID := item.Get("id").Str
// 		if itemID == "" {
// 			continue
// 		}
// 		key := escapeSpecialCharacters(itemID)

// 		hotList = append(hotList, other.HotItem{
// 			Index:       i + 1,
// 			Title:       gjson.Get(jsonStr, "defaultClient."+key+".name").Str,
// 			Description: "",
// 			Image:       gjson.Get(jsonStr, "defaultClient."+key+".poster").Str,
// 			Popularity:  gjson.Get(jsonStr, "defaultClient."+key+".hotValue").Str,
// 			URL: "https://www.kuaishou.com/short-video/" +
// 				gjson.Get(jsonStr, "defaultClient."+key+".photoIds.0").Str +
// 				"?streamSource=hotrank&trendingId=" +
// 				gjson.Get(jsonStr, "defaultClient."+key+".id").Str +
// 				"&area=home",
// 		})
// 	}

// 	return other.HotSearchData{
// 		Source:     "快手热榜",
// 		UpdateTime: time.Now().Format("2006-01-02 15:04:05"),
// 		HotList:    hotList,
// 	}, nil
// }

func escapeSpecialCharacters(str string) string {
	var result strings.Builder

	//遍历字符串的每个字符
	for _, char := range str {
		if char == '.' {
			result.WriteRune('\\')	//在符号前加上反斜杠
		}
		result.WriteRune(char)		//将当前字符添加到结果中
	}
	return result.String()
}

// func escapeSpecialCharacters(s string) string {
// 	replacer := strings.NewReplacer(
// 		".", "\\.",
// 		"(", "\\(",
// 		")", "\\)",
// 		"{", "\\{",
// 		"}", "\\}",
// 		"[", "\\[",
// 		"]", "\\]",
// 		"\"", "\\\"",
// 	)
// 	return replacer.Replace(s)
// }
