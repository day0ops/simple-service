package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/day0ops/simple-service/pkg/config"
)

const (
	defaultAddress      = ":8081"
	defaultShutdownWait = 5 * time.Second
)

type PodName struct {
	Host string `json:"host"`
}

type Server struct {
	ctx     context.Context
	log     *zap.Logger
	address string
	httpSrv *http.Server
}

type Option func(*Server)

func New(ctx context.Context, log *zap.Logger, handler *gin.Engine, opts ...Option) (*Server, error) {
	srv := &Server{
		ctx: ctx,
		log: log,
	}

	for _, opt := range opts {
		opt(srv)
	}

	if srv.address == "" {
		srv.address = defaultAddress
	}

	// register the main routes
	RegisterRoutes(handler)

	srv.httpSrv = &http.Server{
		Addr:    srv.address,
		Handler: handler.Handler(),
	}

	return srv, nil
}

func (s *Server) Serve() error {
	if s.ctx == nil {
		s.ctx = context.TODO()
	}

	errCh := make(chan error, 1)

	go func() {
		s.log.Info("starting server", zap.String("address", s.address))
		if err := s.httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- fmt.Errorf("cannot listen: %w", err)
			return
		}
		errCh <- nil
	}()

	select {
	case <-s.ctx.Done():
		return s.Stop()
	case err := <-errCh:
		return err
	}
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if s.httpSrv != nil {
		s.log.Info("stopping grpc server")
		if err := s.httpSrv.Shutdown(ctx); err != nil {
			s.log.Error("server shutdown", zap.Error(err))
		}
	}
	time.Sleep(defaultShutdownWait)
	return nil
}

func WithServerAddress(address string) Option {
	return func(s *Server) {
		s.address = fmt.Sprintf(":%s", address)
	}
}

func HealthHandler(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusOK)
}

func RootHandler(c *gin.Context) {
	podName := config.GetEnv(config.PodNameEnvVar, os.Getenv(config.HostnameEnvVar))
	if podName == "" {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("pod name is empty"))
		return
	}
	c.JSON(http.StatusOK, PodName{Host: podName})
}
