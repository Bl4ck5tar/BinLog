package upload

import (
	"BinLog/server/global"
	"BinLog/server/utils"
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

type Qiniu struct{

}
func (*Qiniu) UploadImage(file *multipart.FileHeader) (string, string, error) {
	size := float64(file.Size) / float64(1024*1024)
	if size >= float64(global.Config.Upload.Size) {
		return "", "", fmt.Errorf("the image size exceeds the set size, the current size is: %.2f MB, the set size is: %d MB", size, global.Config.Upload.Size)

	}

	ext := filepath.Ext(file.Filename)
	name := strings.TrimSuffix(file.Filename, ext)
	if _, exists := WhiteImageList[ext]; !exists {
		return "", "", errors.New("don't upload files that aren't image types")
	}

	putPolicy := storage.PutPolicy{Scope: global.Config.Qiniu.Bucket}						//定义上传策略，指定上传到哪个存储空间
	mac := qbox.NewMac(global.Config.Qiniu.AccessKey, global.Config.Qiniu.SecretKey)		//生成鉴权对象，签发上传 token
	upToken := putPolicy.UploadToken(mac)													//上传文件权限校验 token
	cfg := qiniuConfig()																	//获取当前存储区域配置
	formUploader := storage.NewFormUploader(cfg)											//初始化表单上传器
	putRet := storage.PutRet{}																//接收上传返回结果
	putExtra := storage.PutExtra{Params :map[string]string{}}								//可选额外参数

	fileKey := utils.MD5V([]byte(name)) + "-" + time.Now().Format("20060101150405") + ext

	data, err := file.Open()
	if err != nil {
		return "", "", err
	}
	defer data.Close()

	err = formUploader.Put(context.Background(), &putRet, upToken, fileKey, data, file.Size, &putExtra)
	if err != nil {
		return "", "", err
	}
	return global.Config.Qiniu.ImgPath + putRet.Key, putRet.Key, nil
}

func (*Qiniu) DeleteImage(key string) error {
	mac := qbox.NewMac(global.Config.Qiniu.AccessKey, global.Config.Qiniu.SecretKey)
	cfg := qiniuConfig()
	bucketManager := storage.NewBucketManager(mac, cfg)
	return bucketManager.Delete(global.Config.Qiniu.Bucket, key)
}

func qiniuConfig() *storage.Config {
	cfg := storage.Config{
		UseHTTPS: 		global.Config.Qiniu.UseHTTPS,
		UseCdnDomains: 	global.Config.Qiniu.UseCdnDomains,
	}
	switch global.Config.Qiniu.Zone {
	case "z0", "ZoneHuadong":
		cfg.Zone = &storage.ZoneHuadong
	case "z1", "ZoneHuabei":
		cfg.Zone = &storage.ZoneHuabei
	case "z2", "ZoneHuanan":
		cfg.Zone = &storage.ZoneHuanan
	case "na0", "ZoneBeimei":
		cfg.Zone = &storage.ZoneBeimei
	case "as0", "ZoneXinjiapo":
		cfg.Zone = &storage.ZoneXinjiapo
	case "ZoneHuadongZheJiang2":
		cfg.Zone = &storage.ZoneHuadongZheJiang2
	}
	return &cfg
}