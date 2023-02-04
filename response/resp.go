package response

import (
	"encoding/json"
	"github.com/BlsaCn/short-url/tools"
	"log"
	"net/http"
)

// WithErr 返回错误
func WithErr(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case tools.Error:
		log.Printf("HTTP Error %d %s\n", e.Status(), e.Error())
		Success(w, e.Status(), e.Error(), nil)
	default:
		Success(w, 1, err.Error(), nil)
	}
}

// Success 成功返回
func Success(w http.ResponseWriter, code int, msg string, data interface{}) {
	if code > 0 {
		log.Printf("[Error] %s \n", msg)
	}
	resp(w, http.StatusOK, code, msg, data)
}

// Fail 失败返回
func Fail(w http.ResponseWriter, code int, msg string, data interface{}) {
	resp(w, http.StatusInternalServerError, code, msg, data)
}

func resp(w http.ResponseWriter, httpStatus, code int, msg string, data interface{}) {
	m := make(map[string]interface{})
	m["msg"] = msg
	m["code"] = code
	m["data"] = data
	b, _ := json.Marshal(&m)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	_, _ = w.Write(b)
}
