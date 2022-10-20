package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Context struct {
	*gin.Context
}

func (c *Context) Bind(v interface{}) error {
	return c.Context.ShouldBindJSON(v)
}

func (c *Context) Param(param string) string {
	return c.Context.Param(param)
}

func (c *Context) Query(query string) string {
	return c.Context.Request.URL.Query().Get(query)
}

func (c *Context) JSON(statusCode int, v interface{}) {
	c.Context.JSON(statusCode, v)
}

type Router struct {
	*gin.Engine
}

func (r *Router) POST(path string, handler func(Context)) {
	r.Engine.POST(path, NewGinHandler(handler))
}

func (r *Router) GET(path string, handler func(Context)) {
	r.Engine.GET(path, NewGinHandler(handler))
}

func (r *Router) PUT(path string, handler func(Context)) {
	r.Engine.PUT(path, NewGinHandler(handler))
}

func (r *Router) DELETE(path string, handler func(Context)) {
	r.Engine.DELETE(path, NewGinHandler(handler))
}

func NewMyContext(c *gin.Context) *Context {
	return &Context{Context: c}
}

func NewGinHandler(handler func(Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler(*NewMyContext(c))
	}
}

func NewRouter() *Router {
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{
		"*",
	}

	config.AllowHeaders = []string{
		"Content-Type",
		"Content-Length",
		"Accept-Encoding",
		"X-CSRF-Token",
		"Authorization",
		"Access-Token",
		"accept",
		"origin",
		"Cache-Control",
		"X-Requested-With",
		// "TransactionID",
	}

	config.AllowMethods = []string{
		"POST",
		"OPTIONS",
		"GET",
		"PUT",
		"DELETE",
	}
	r.Use(cors.New(config))
	return &Router{r}
}
