package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HTTPServer struct {
	httpHandlers *any
}

func NewHTTPServer(httpHandlers *any) *HTTPServer {
	return &HTTPServer{
		httpHandlers: httpHandlers,
	}
}

func (s *HTTPServer) Start() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "ok")
	})

	r.Run(":9091")
}
