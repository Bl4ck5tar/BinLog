package appTypes

import (
	"encoding/json"
)

type Storage int

const (
	Local Storage = iota
	Qiniu
)
//实现json.Marshal接口
func (s Storage) MarshaJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *Storage) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	*s = ToStorage(str)
	return nil
}
//String 方法返回 Storage 的字符串表示
func (s Storage) String() string {
	var str string
	switch s {
		case Local:
			str = "本地"
		case Qiniu:
			str = "七牛"
		default:
			str = "未知存储"
	}
	return str
}
//ToStorage 方法将字符串转换为Storage
func ToStorage(str string) Storage {
	switch str {
		case "本地":
			return Local
		case "七牛":
			return Qiniu
		default:
			return -1
	}
}