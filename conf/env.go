package conf

import (
	"gopkg.in/ini.v1"
)

// 读取环境变量相关
// 可以改成读取配置文件

var Cfg *ini.File

type RedisCfg struct {
	Host string
	Pwd  string
	DB   int
}

func init() {
	var err error
	if Cfg == nil {
		Cfg, err = ini.Load("conf/cfg.ini")
	}
	if err != nil {
		panic(err)
	}
}

func (r *RedisCfg) Redis() *RedisCfg {
	r.Host = Cfg.Section("redis").Key("HOST").String()
	r.Pwd = Cfg.Section("redis").Key("PWD").String()
	r.DB, _ = Cfg.Section("redis").Key("DB").Int()
	return r
}
