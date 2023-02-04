package storage

// Storage 数据存储器
type Storage interface {
	// Shorten 长连接转短链接
	Shorten(url string, exp int64) (string, error)
	// ShortLinkInfo 短连接对应的信息
	ShortLinkInfo(shortUrl string) (interface{}, error)
	// UnShorten 短链接还原长连接
	UnShorten(shortUrl string) (string, error)
}
