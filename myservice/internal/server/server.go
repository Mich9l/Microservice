package server

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"myservice/internal/application/processor"
)

type HTTPServer struct {
	addr string
	srv  *http.Server
	log  *zap.Logger
}

func New(addr string, _ *processor.RecipeProcessor, log *zap.Logger) (*HTTPServer, error) {
	engine := gin.New()
	engine.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })

	return &HTTPServer{
		addr: addr,
		srv:  &http.Server{Addr: addr, Handler: engine},
		log:  log,
	}, nil
}

func (s *HTTPServer) Start() error {
	return s.srv.ListenAndServe()
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}