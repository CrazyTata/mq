package qiniu

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"mq/infrastructure/svc"

	"mq/common/util"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

type Uploader struct {
	svcCtx *svc.ServiceContext
}

func NewUploader(svcCtx *svc.ServiceContext) *Uploader {
	return &Uploader{
		svcCtx: svcCtx,
	}
}

func (u *Uploader) UploadFromURL(url string) (string, error) {
	// 下载图片到本地
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 创建临时文件
	tempFile, err := os.CreateTemp("", "upload-*.png")
	if err != nil {
		return "", err
	}
	defer os.Remove(tempFile.Name())

	// 将图片写入临时文件
	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		return "", err
	}

	// 上传临时文件到七牛云
	mac := qbox.NewMac(u.svcCtx.GetConfig().Qiniu.AccessKey, u.svcCtx.GetConfig().Qiniu.SecretKey)
	cfg := storage.Config{
		Zone:          dealRegion(u.svcCtx.GetConfig().Qiniu.Region),
		UseHTTPS:      false,
		UseCdnDomains: false,
	}
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	putPolicy := storage.PutPolicy{
		Scope: u.svcCtx.GetConfig().Qiniu.Bucket,
	}
	upToken := putPolicy.UploadToken(mac)

	pathExt := "knowledge/image"

	fileKey := fmt.Sprintf("%s/%s/%s", pathExt, time.Now().Format("20060102"), util.NewSnowflake().String()+".png") // 文件名格式 自己可以改 建议保证唯一性

	err = formUploader.PutFile(context.Background(), &ret, upToken, fileKey, tempFile.Name(), &storage.PutExtra{Params: map[string]string{"x:name": "github logo"}})
	if err != nil {
		return "", err
	}
	return u.svcCtx.GetConfig().Qiniu.Domain + "/" + ret.Key, nil
}

func (u *Uploader) UploadFilesWithStructure(ctx context.Context, baseDir, pathType string) error {
	// 初始化七牛云上传配置
	mac := qbox.NewMac(u.svcCtx.GetConfig().Qiniu.AccessKey, u.svcCtx.GetConfig().Qiniu.SecretKey)
	putPolicy := storage.PutPolicy{
		Scope: u.svcCtx.GetConfig().Qiniu.Bucket,
	}
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{
		Zone:          dealRegion(u.svcCtx.GetConfig().Qiniu.Region),
		UseHTTPS:      false,
		UseCdnDomains: false,
	}
	formUploader := storage.NewFormUploader(&cfg)

	return filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			// 构建上传路径，保持目录结构
			relativePath := strings.TrimPrefix(path, baseDir+"/")
			key := fmt.Sprintf("knowledge/%s/%s", pathType, relativePath)

			// 上传文件
			err := formUploader.PutFile(ctx, &storage.PutRet{}, upToken, key, path, &storage.PutExtra{})
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func dealRegion(origin string) *storage.Region {
	// 根据配置的 Region 字符串选择合适的 *storage.Region
	switch origin {
	case "huadong":
		return &storage.ZoneHuadong
	case "huabei":
		return &storage.ZoneHuabei
	case "huanan":
		return &storage.ZoneHuanan
	case "beimei":
		return &storage.ZoneBeimei
	default:
		return nil
	}
}

func (u *Uploader) GetUploadToken(ctx context.Context) (string, error) {
	putPolicy := storage.PutPolicy{Scope: u.svcCtx.GetConfig().Qiniu.Bucket}
	mac := qbox.NewMac(u.svcCtx.GetConfig().Qiniu.AccessKey, u.svcCtx.GetConfig().Qiniu.SecretKey)
	upToken := putPolicy.UploadToken(mac)
	return upToken, nil
}
