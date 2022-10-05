package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Router struct {
	*gin.Engine
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

// func (r *Router) POST(path string, handler HandlerFunc) {
// 	r.Engine.POST(path, func(c *gin.Context) {
// 		handler(&Context{c})
// 	})
// }
