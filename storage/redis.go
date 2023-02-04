package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/BlsaCn/short-url/db"
	"github.com/BlsaCn/short-url/tools"
	"github.com/go-redis/redis"
	"github.com/pilu/go-base62"
	"time"
)

const (
	UrlIdKey           string = "next.url.id"         // 全局自增器
	ShortLinkKey       string = "shortLink:url:%s"    // 短地址和长地址的映射
	UrlHashKey         string = "urlHash:url:%s"      // 地址hash和短地址的映射
	ShortLinkDetailKey string = "shortLink:detail:%s" // 短地址和详情的映射
	// 步骤：自增id转sha1，用sha1绑定短地址，短地址分别绑定长地址和详情
)

type RedisCli struct {
	Cli *redis.Client
}

type URLDetail struct {
	// 长链接地址
	Url string `json:"url"`
	// 创建时间
	CreatedAt string `json:"created_at"`
	// 过期时间(分钟)
	ExpByMinute time.Duration `json:"exp"`
}

func NewRedis() *RedisCli {
	r := &RedisCli{}
	r.Cli = db.NewRedisCli()
	return r
}

// Shorten 长连接转短链接
func (r *RedisCli) Shorten(url string, exp int64) (string, error) {
	// url 转换成 sha1
	sha := tools.ToSha1(url)
	// 判断该url是否已经转换过
	shortUrl, err := r.Cli.Get(fmt.Sprintf(UrlHashKey, sha)).Result()
	if err != redis.Nil {
		if err != nil {
			return "", err
		} else {
			return shortUrl, nil
		}
	}

	// 自增计数器
	IncrId, err := r.Cli.Incr(UrlIdKey).Result()
	if err != nil {
		return "", err
	}

	// IncrId 使用base62编码生成短地址
	shortUrlNew := base62.Encode(int(IncrId))
	// 过期时间
	duration := time.Minute * time.Duration(exp)
	err = r.Cli.Set(fmt.Sprintf(ShortLinkKey, shortUrlNew), url, duration).Err()
	if err != nil {
		return "", err
	}
	err = r.Cli.Set(fmt.Sprintf(UrlHashKey, sha), shortUrlNew, duration).Err()
	if err != nil {
		return "", err
	}

	detail := URLDetail{
		Url:         url,
		CreatedAt:   time.Now().Format("2006-01-02 15:04:05"),
		ExpByMinute: time.Duration(exp),
	}
	d, err := json.Marshal(&detail)
	if err != nil {
		return "", err
	}
	err = r.Cli.Set(fmt.Sprintf(ShortLinkDetailKey, shortUrlNew), d, duration).Err()
	if err != nil {
		return "", err
	}

	return shortUrlNew, nil
}

// ShortLinkInfo 短连接对应的信息
func (r *RedisCli) ShortLinkInfo(shortUrl string) (interface{}, error) {
	url, err := r.Cli.Get(fmt.Sprintf(ShortLinkDetailKey, shortUrl)).Result()
	if err == redis.Nil {
		return nil, errors.New("此短链接不存在")
	}
	if err != nil {
		return nil, err
	}
	return url, nil
}

// UnShorten 短链接还原长连接
func (r *RedisCli) UnShorten(shortUrl string) (string, error) {
	url, err := r.Cli.Get(fmt.Sprintf(ShortLinkKey, shortUrl)).Result()
	if err == redis.Nil {
		return "", errors.New("此短链接不存在")
	}
	if err != nil {
		return "", err
	}
	return url, nil
}
