// draft context

package draft

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type H map[string]interface{}

// Writer origin http.ResponseWriter
// Req    *http.Request
// Path   request path
// Method request method
// Params request params
// StatusCode response code
// handlers middleware
// index middleware index
type Context struct {
	// origin objects
	Writer http.ResponseWriter
	Req    *http.Request
	// request info
	Path   string
	Method string
	Params map[string]string
	// response info
	StatusCode int
	// middleware
	handlers []HandlerFunc
	index    int
	// engine pointer
	engine *Engine
}

func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

func (c *Context) File(filepath string) {
	http.ServeFile(c.Writer, c.Req, filepath)
}

func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.JSON(code, H{"message": err})
}

func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) QueryParam(key string) string {
	return c.queryTypeParam(key)
}

func (c *Context) QueryType(key string) string {
	return c.queryTypeParam(key)
}

func (c *Context) queryTypeParam(t string) string {
	urlarry := strings.Split(c.Req.RequestURI, "/")
	if t == "type" {
		return urlarry[len(urlarry)-2]
	} else if t == "param" {
		return urlarry[len(urlarry)-1]
	}
	return ""
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// return string
func (c *Context) String(code int, format string, values ...interface{}) {
	c.Status(code)
	c.SetHeader("Content-Type", "text/plain")
	_, _ = c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// return json format data
func (c *Context) JSON(code int, obj interface{}) {
	c.Status(code)
	c.SetHeader("Content-Type", "application/json")
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

// return bytes
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	_, _ = c.Writer.Write(data)
}

// HTML template render
// refer https://golang.org/pkg/html/template/
func (c *Context) HTML(code int, name string, data interface{}) {
	c.Writer.WriteHeader(code)
	c.SetHeader("Content-Type", "text/html")

	//if _, ok := CacheMap.Get(c.Req.RequestURI); !ok && strings.Contains(c.Req.RequestURI, "detail") && strings.Contains(c.Req.RequestURI, "id") {
	//	cache := &Cache{}
	//	mwriter := io.MultiWriter(cache, c.Writer)
	//	if err := c.engine.htmlTemplates.ExecuteTemplate(mwriter, name, data); err != nil {
	//		c.Fail(500, err.Error())
	//	}
	//	CacheMap.Lock()
	//	CacheMap.Caches[c.Req.RequestURI] = cache.Value
	//	fmt.Print("##################################################")
	//	CacheMap.Unlock()
	//	cache = nil
	//} else {
	if err := c.engine.htmlTemplates.ExecuteTemplate(c.Writer, name, data); err != nil {
		c.Fail(500, err.Error())
	}
	//}
}

func (c *Context) reset() {
	c.handlers = nil
	c.index = -1
	c.Path = ""
}
