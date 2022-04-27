package oss

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"time"
)

import (
	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
)

import (
	"github.com/glory-go/glory/config"
)

type OssService interface {
	loadConfig(conf *config.OssConfig)
	setup()
	SaveFile(file io.Reader, fileName string) error
	SaveByte([]byte, string) error
	DeleteFile(filePath string) error
	GetFileUrl(filePath string) string
}

type qiniuOssService struct {
	config           qiniuOssConfig
	mac              *qbox.Mac
	defaultUploader  *storage.FormUploader
	defaultPutPolicy *storage.PutPolicy
}

type qiniuOssConfig struct {
	buckname      string
	ossAccessKey  string
	ossSecretKey  string
	ossDomainName string
}

func (qos *qiniuOssService) loadConfig(conf *config.OssConfig) {
	qos.config.buckname = conf.Buckname
	qos.config.ossAccessKey = conf.OssAccessKey
	qos.config.ossSecretKey = conf.OssSecretKey
	qos.config.ossDomainName = conf.OssDomainName
}

func (qos *qiniuOssService) setup() {
	qos.defaultPutPolicy = &storage.PutPolicy{
		Scope: qos.config.buckname,
	}

	qos.mac = qbox.NewMac(qos.config.ossAccessKey, qos.config.ossSecretKey)

	cfg := storage.Config{}
	cfg.Zone = &storage.ZoneHuanan
	cfg.UseHTTPS = false
	cfg.UseCdnDomains = false
	qos.defaultUploader = storage.NewFormUploader(&cfg)
}

func (qos *qiniuOssService) SaveFile(file io.Reader, id string) error {
	upToken := qos.defaultPutPolicy.UploadToken(qos.mac)
	ret := storage.PutRet{}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("SaveProjectFile error: Can't read file content, %v", err)
	}
	if err := qos.defaultUploader.Put(context.Background(), &ret, upToken,
		id, bytes.NewBuffer(data), int64(len(data)), &storage.PutExtra{}); err != nil {
		return fmt.Errorf("SaveProjectFile error: Upload file with %s", err)
	}
	return nil
}

func (qos *qiniuOssService) SaveByte(data []byte, filePath string) error {
	return qos.SaveFile(bytes.NewReader(data), filePath)
}

func (qos *qiniuOssService) DeleteFile(filePath string) error {
	cfg := storage.Config{
		UseHTTPS: false,
	}
	bucketManager := storage.NewBucketManager(qos.mac, &cfg)
	if err := bucketManager.Delete(qos.config.buckname, filePath); err != nil {
		return fmt.Errorf("DeleteFileFromOSS delete file error: %v", err)
	}
	return nil
}

// getFileUrl 返回远端文件的地址
func (qos *qiniuOssService) GetFileUrl(filePath string) string {
	return storage.MakePrivateURL(
		qos.mac,
		qos.config.ossDomainName,
		filePath,
		time.Now().Add(time.Second*3600).Unix(),
	)
}

func newQiniuService() *qiniuOssService {
	return &qiniuOssService{}
}
