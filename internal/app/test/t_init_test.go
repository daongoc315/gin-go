package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/chunganhbk/gin-go/internal/app"
	"github.com/chunganhbk/gin-go/internal/app/router"
	"io"
	"net/http"
	"net/url"

	"github.com/chunganhbk/gin-go/internal/app/config"
	"github.com/gin-gonic/gin"
)

const (
	configFile = "../../../configs/config.toml"
	modelFile  = "../../../configs/model.conf"
	apiPrefix  = "/api/"
	token      = ""
)

var engine *gin.Engine

func init() {

	config.MustLoad(configFile)

	config.C.RunMode = "test"
	config.C.Log.Level = 2
	config.C.Casbin.Model = modelFile
	config.C.Gorm.Debug = true
	config.C.Gorm.DBType = "sqlite3"

	container, _ := app.BuildContainer()
	engine = router.InitGinEngine(container)
}

// ResID
type ResID struct {
	ID string `json:"id,omitempty"`
}

func toReader(v interface{}) io.Reader {
	buf := new(bytes.Buffer)
	_ = json.NewEncoder(buf).Encode(v)
	return buf
}

func parseReader(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

func parseSuccess(r io.Reader) error {
	var status struct {
		Status string `json:"msg"`
	}
	err := parseReader(r, &status)
	if err != nil {
		return err
	}
	if status.Status != "Success" {
		return errors.New("not OK")
	}
	return nil
}

func newPageParam(extra ...map[string]string) map[string]string {
	data := map[string]string{
		"current":  "1",
		"pageSize": "1",
	}

	if len(extra) > 0 {
		for k, v := range extra[0] {
			data[k] = v
		}
	}

	return data
}

type PaginationResult struct {
	Total    int64 `json:"total"`
	Current  int   `json:"current"`
	PageSize int   `json:"pageSize"`
}

type PageResult struct {
	List       interface{}       `json:"list"`
	Pagination *PaginationResult `json:"pagination"`
}

func parsePageReader(r io.Reader, v interface{}) error {
	result := &PageResult{List: v}
	return parseReader(r, result)
}

func newPostRequest(formatRouter string, v interface{}, args ...interface{}) *http.Request {
	req, _ := http.NewRequest("POST", fmt.Sprintf(formatRouter, args...), toReader(v))
	req.Header.Add("Authorization", token)
	return req
}

func newPutRequest(formatRouter string, v interface{}, args ...interface{}) *http.Request {
	req, _ := http.NewRequest("PUT", fmt.Sprintf(formatRouter, args...), toReader(v))
	req.Header.Add("Authorization", token)
	return req
}

func newDeleteRequest(formatRouter string, args ...interface{}) *http.Request {
	req, _ := http.NewRequest("DELETE", fmt.Sprintf(formatRouter, args...), nil)
	req.Header.Add("Authorization", token)
	return req
}

func newGetRequest(formatRouter string, params map[string]string, args ...interface{}) *http.Request {
	values := make(url.Values)
	for k, v := range params {
		values.Set(k, v)
	}

	urlStr := fmt.Sprintf(formatRouter, args...)
	if len(values) > 0 {
		urlStr += "?" + values.Encode()
	}

	req, _ := http.NewRequest("GET", urlStr, nil)
	return req
}
