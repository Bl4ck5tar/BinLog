package config

type Upload struct {
	Size		int		`json:"size" yaml:"size"`	//图片上传大小，单位MB
	Path		string	`json:"path" yaml:"path"`	//图片上传路径
}