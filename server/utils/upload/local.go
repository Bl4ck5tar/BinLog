package upload

import (
	"BinLog/server/global"
	"BinLog/server/utils"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Local struct{

}

func (*Local) UploadImage(file *multipart.FileHeader) (string, string, error) {
	//校验文件大小
	size := float64(file.Size) / float64(1024*1024)
	if size >= float64(global.Config.Upload.Size) {
		return "", "", fmt.Errorf("the image size exceeds the set size, the current size is: %.2f MB, the set size is: %d MB", size, global.Config.Upload.Size)

	}

	//取文件扩展名进行校验
	ext := filepath.Ext(file.Filename)
	name := strings.TrimSuffix(file.Filename, ext)
	if _, exists := WhiteImageList[ext]; !exists {
		return "", "", errors.New("don't upload files that aren't image types")
	}

	//生成md5文件名（避免冲突+安全化）
	filename := utils.MD5V([]byte(name)) + "-" + time.Now().Format("20060101150405") + ext
	path := global.Config.Upload.Path + "/image/"

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return "", "", err
	}

	filepath := path + filename

	out, err := os.Create(filepath)
	if err != nil {
		return "", "", err
	}
	defer out.Close()

	f, err := file.Open()
	if err != nil {
		return "", "", err
	}
	defer f.Close()

	if _, err = io.Copy(out, f); err != nil {
		return "", "", err
	}
	return "/" + filepath, filename, nil
}

func (*Local) DeleteImage(key string) error {
	path := global.Config.Upload.Path + "/image/" + key
	return os.Remove(path)
}