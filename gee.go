package go_webs

import (
	"log"
	"net/http"
	"strings"
)

type HandleFunc func(c *Context)

type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup
}

type RouterGroup struct {
	prefix      string
	middlewares []HandleFunc // 中间件
	parent      *RouterGroup
	engine      *Engine
}

func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (g *RouterGroup) Group(prefix string) *RouterGroup {
	engine := g.engine
	group := &RouterGroup{
		prefix: g.prefix + prefix,
		parent: g,
		engine: engine,
	}
	engine.groups = append(engine.groups, group)
	return group
}

func (g *RouterGroup) addRouter(method string, comp string, handler HandleFunc) {
	pattern := g.prefix + comp
	log.Printf("Route %s - %s", method, pattern)
	g.engine.router.addRouter(method, pattern, handler)
}
func (g *RouterGroup) GET(pattern string, handler HandleFunc) {
	g.addRouter("GET", pattern, handler)
}

func (g *RouterGroup) POST(pattern string, handler HandleFunc) {
	g.addRouter("POST", pattern, handler)
}

func (g *RouterGroup) Use(middleware ...HandleFunc) {
	g.middlewares = append(g.middlewares, middleware...)
}

func (e *Engine) addRouter(method string, pattern string, handle HandleFunc) {
	log.Printf("Route %s - %s", method, pattern)
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
	var middleware []HandleFunc

	// 路由组中间件
	for _, group := range e.groups {
		if strings.HasPrefix(r.URL.Path, group.prefix) {
			middleware = append(middleware, e.middlewares...)
		}
	}

	c := newContext(w, r)
	c.handlers = middleware
	e.router.handler(c)
}
