package go_webs

import "net/http"

type HandleFunc func(c *Context)

type Engine struct {
	router *router
}

func New() *Engine {
	return &Engine{
		router: newRouter(),
	}
}

func (e *Engine) addRouter(method string, pattern string, handle HandleFunc) {
	e.router.addRouter(method, pattern, handle)
}

func (e *Engine) POST(pattern string, handle HandleFunc) {
	e.addRouter("POST", pattern, handle)
}

func (e *Engine) GET(pattern string, handle HandleFunc) {
	e.addRouter("GET", pattern, handle)
}

func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := NewContext(w, r)
	e.router.handler(c)
}
