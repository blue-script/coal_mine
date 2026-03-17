package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HTTPServer struct {
	httpHandlers *HTTPHandlers
	server       *http.Server
}

func NewHTTPServer(httpHandlers *HTTPHandlers) *HTTPServer {
	router := gin.Default()

	router.POST("/miners", httpHandlers.HireMiner)
	router.POST("/equipment/:name/purchase", httpHandlers.BuyEquipment)
	router.GET("/miners/cost/:class", httpHandlers.HireCost)

	server := &HTTPServer{
		httpHandlers: httpHandlers,
		server: &http.Server{
			Addr:    ":9091",
			Handler: router,
		},
	}

	return server
}

func (s *HTTPServer) Start() error {
	return s.server.ListenAndServe()
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
