package utils

import "regexp"

func FindIllustrations(text string) ([]string, error) {
	//正则表达式匹配 Markdown 图片语法 ![描述](图片链接)
	regex := `!\[([^\]]*)\]\(([^)]+)\)`

	//编译正则表达式
	re, err := regexp.Compile(regex)
	if err != nil {
		return nil, err
	}
	//查找所有匹配项，返回一个二维切片[][]string
	matches := re.FindAllStringSubmatch(text, -1)
	//matches[0]为整个匹配字符串，match[1]为第一个捕获组-图片描述，match[2]为第二个捕获组-图片链接

	var illustrations []string

	for _, match := range matches {
		if len(match) > 2 {
			illustrations = append(illustrations, match[2])
		}
	}
	return illustrations, nil
}