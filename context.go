package go_webs

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	// 元数据
	Write http.ResponseWriter
	Req   *http.Request

	// 请求数据
	Path   string
	Method string
	Params map[string]string

	// 响应数据
	StatusCode int
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Write:  w,
		Req:    r,
		Path:   r.URL.Path,
		Method: r.Method,
	}
}

func (c *Context) Param(key string) string {
	val, _ := c.Params[key]
	return val
}

func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Write.WriteHeader(code)
}

func (c *Context) SetHeader(key, val string) {
	c.Write.Header().Set(key, val)
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Write.Write([]byte(fmt.Sprintf(format, values)))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)

	encoder := json.NewEncoder(c.Write)
	err := encoder.Encode(obj)
	if err != nil {
		http.Error(c.Write, err.Error(), 500)
		return
	}
}

func (c *Context) Data(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Write.Write([]byte(html))
}