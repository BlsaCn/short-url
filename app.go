package main

import (
	"encoding/json"
	"github.com/BlsaCn/short-url/response"
	"github.com/BlsaCn/short-url/storage"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"log"
	"net/http"
)

// App 包含 Router，Middleware，Storage
type App struct {
	Router      *mux.Router
	Middlewares *Middleware
	S           storage.Storage
}

// 转短链接请求结构体
type shortenReq struct {
	Url           string `json:"url"`            // url
	ExpirationMin int64  `json:"expiration_min"` // 过期时间(分钟)
}

// NewApp 初始化App
func NewApp() *App {
	a := &App{}
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	a.Router = mux.NewRouter()
	a.S = storage.NewRedis()
	a.Middlewares = &Middleware{}
	a.initRouters()
	return a
}

func (a *App) initRouters() {
	// a.Router.HandleFunc("/api/shorten", a.createShorten).Methods("POST")
	// a.Router.HandleFunc("/api/info", a.shortenInfo).Methods("GET")
	// a.Router.HandleFunc("/api/{shortLink:[a-zA-Z0-9]{1,11}}", a.redirect).Methods("GET")
	// Alice 包使用中间件
	c := alice.New(a.Middlewares.LoggingHandler, a.Middlewares.RecoverHandler)
	a.Router.Handle("/api/shorten", c.ThenFunc(a.createShorten)).Methods("POST")
	a.Router.Handle("/api/info", c.ThenFunc(a.shortenInfo)).Methods("GET")
	a.Router.Handle("/api/{shortLink:[a-zA-Z0-9]{1,11}}", c.ThenFunc(a.redirect)).Methods("GET")
}

// 生成短链接
func (a *App) createShorten(w http.ResponseWriter, r *http.Request) {
	var req shortenReq
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WithErr(w, err)
		return
	}
	// @todo 需校验req.Url 和 req.ExpirationMin

	shorten, err := a.S.Shorten(req.Url, req.ExpirationMin)
	if err != nil {
		response.WithErr(w, err)
		return
	}
	response.Success(w, 0, "", shorten)
}

// 短链接信息
func (a *App) shortenInfo(w http.ResponseWriter, r *http.Request) {
	val := r.URL.Query()
	s := val.Get("shortLink")
	info, err := a.S.ShortLinkInfo(s)
	if err != nil {
		response.WithErr(w, err)
		return
	}
	response.Success(w, 0, "", info)
}

// 重定向接口，跳转到长连接
func (a *App) redirect(w http.ResponseWriter, r *http.Request) {
	// 获取参数
	val := mux.Vars(r)
	// 根据短链接获取长连接
	url, err := a.S.UnShorten(val["shortLink"])
	if err != nil {
		response.WithErr(w, err)
		return
	}

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// Run 运行
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}
