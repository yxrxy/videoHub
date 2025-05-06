package upyun

import (
	"bytes"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/yxrxy/videoHub/config"
	"github.com/yxrxy/videoHub/pkg/errno"
)

// UploadImg 又拍云上传文件
func UploadYun(file []byte, url string) error {
	body := bytes.NewReader(file)
	req, err := http.NewRequest("PUT", url, body)
	if err != nil {
		return errno.UpYunFileError.WithMessage(err.Error())
	}
	req.SetBasicAuth(config.Upyun.Operator, config.Upyun.Password)
	req.Header.Add("Date", time.Now().UTC().Format(http.TimeFormat))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errno.UpYunFileError.WithMessage(err.Error())
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Printf("uploadimg close request meet error: %v", err)
		}
	}()
	if res.StatusCode != http.StatusOK {
		return errno.UpYunFileError
	}
	return nil
}

// DeleteImg 又拍云删除文件
func DeleteYun(url string) error {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return errno.UpYunFileError.WithMessage(err.Error())
	}
	req.SetBasicAuth(config.Upyun.Operator, config.Upyun.Password)
	req.Header.Add("Date", time.Now().UTC().Format(http.TimeFormat))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errno.UpYunFileError.WithMessage(err.Error())
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Printf("deleteImg close request meet error: %v", err)
		}
	}()
	if res.StatusCode != http.StatusOK {
		return errno.UpYunFileError
	}
	return nil
}

// 获取图片URL
func GetImageUrl(url string) string {
	return strings.Join([]string{config.Upyun.ImageDomain, strings.TrimPrefix(url, config.Upyun.UssDomain)}, "")
}

// 获取视频URL
func GetVideoUrl(url string) string {
	return strings.Join([]string{config.Upyun.VideoDomain, strings.TrimPrefix(url, config.Upyun.UssDomain)}, "")
}
