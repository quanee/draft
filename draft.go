package draft

import (
	"html/template"
	"net/http"
	"strings"
	"sync"
)

type HandlerFunc func(*Context)

type Engine struct {
	*RouterGroup
	router        *router
	groups        []*RouterGroup     // 存储所有group
	htmlTemplates *template.Template // html渲染
	funcMap       template.FuncMap   // html渲染
	pool          sync.Pool
}

// 创建Engine
func New() *Engine {
	engine := &Engine{
		router:  newRouter(),
		funcMap: template.FuncMap{},
	}
	engine.pool.New = func() interface{} {
		return engine.allocateContext()
	}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// 自定义渲染函数
func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}

// http服务器开始运行的方法
func (engine *Engine) Run(address string) (err error) {
	//engine.log.Debug("Listening and serving HTTP on %s\n", address)
	err = http.ListenAndServe(address, engine)

	return
}

func (engine *Engine) RunTLS(address, certFile, keyFile string) (err error) {
	err = http.ListenAndServeTLS(address, certFile, keyFile, engine)
	return
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.handlers...)
		}
	}
	c := engine.pool.Get().(*Context)
	c.reset()
	c.Path = req.URL.Path
	c.Method = req.Method
	c.Req = req
	c.Writer = w
	c.handlers = middlewares
	c.engine = engine
	engine.router.handle(c)
	engine.pool.Put(c)
}

func (engine *Engine) allocateContext() *Context {
	return &Context{}
}
