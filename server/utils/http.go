package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
)

//HttpRequest 函数用于发送 HTTP 请求
func HttpRequest(
	urlStr string,								//请求的 URL 字符串
	method string,								//请求方法（GET，POST等）
	headers map[string]string,					//请求头（如 Content-Type 等）
	params map[string]string,					//查询参数（如 ?key=value&key2=value2 )
	data any) (*http.Response, error) {			//请求体的内容（如果有的话）
		//创建 URL 对象
		u, err := url.Parse(urlStr)				//将 urlStr 解析为 URL 对象
		if err != nil {
			return nil, err						//如果 URL 解析失败，返回错误
		}

		//向 URL 添加查询参数
		query := u.Query()						//获取 URL 中的查询部分（如果有）
		for k, v := range params {				//遍历参数并添加到 URL 中
			query.Set(k, v)						//使用 Set 方法保证参数键值对唯一
		}
		u.RawQuery = query.Encode()				//更新 URL 的查询部分

		//将请求体数据（如果有）编码成 JSON 格式
		buf := new(bytes.Buffer)				//创建一个缓冲区用于存储请求体
		if data != nil {
			b, err := json.Marshal(data)		//将 data 编码为 JSON 字节数组
			if err != nil {
				return nil, err					//如果编码失败，返回错误
			}
			buf = bytes.NewBuffer(b)			//将编码后的字节数组转换为缓冲区
		}

		//创建 HTTP 请求对象
		req, err := http.NewRequest(method, u.String(), buf)	//使用指定的 URL 和方法创建请求
		if err != nil {
			return nil, err						//如果请求创建失败，返回错误
		}

		//设置请求头
		for k, v := range headers {				//遍历传入的头部信息并设置到请求中
			req.Header.Set(k, v)				//设置头部
		}

		//如果请求体存在，将 Content-Type 设置为 application/json
		if data != nil {
			req.Header.Set("Content-Type", "application/json")	//设置请求头为 JSON 类型
		}

		//发送 HTTP 请求并获取响应
		resp, err := http.DefaultClient.Do(req)	//使用默认的 HTTP 客户端发送请求
		if err != nil {
			return nil, err						//如果请求失败，返回错误
		}
		return resp, nil						//返回响应，调用者可根据需要处理响应数据

	}