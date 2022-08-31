package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/zhufuyi/goctl/templates/http_server/config"
	"github.com/zhufuyi/goctl/templates/http_server/internal/routers"

	"github.com/gin-gonic/gin"
)

var _ IServer = (*httpServer)(nil)

type httpServer struct {
	addr   string
	server *http.Server
}

// Start http service
func (h *httpServer) Start() error {
	if err := h.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("listen server error: %v", err)
	}
	return nil
}

// Stop http service
func (h *httpServer) Stop() error {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second) //nolint
	return h.server.Shutdown(ctx)
}

// String comment
func (h *httpServer) String() string {
	return "http service, addr = " + h.addr
}

// NewHTTPServer creates a new web server
func NewHTTPServer(addr string, readTimeout time.Duration, writeTimeout time.Duration) IServer {
	if config.Get().ServiceEnv == "prod" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := routers.NewRouter()
	server := &http.Server{
		Addr:           addr,
		Handler:        router,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	return &httpServer{
		addr:   addr,
		server: server,
	}
}
