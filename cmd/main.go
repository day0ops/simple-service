package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"

	"github.com/day0ops/simple-service/pkg/config"
	"github.com/day0ops/simple-service/pkg/server"
	"github.com/day0ops/simple-service/pkg/version"
)

var (
	httpPort = flag.String("port", "8081", "port used for the server")
)

func main() {
	os.Exit(start())
}

func start() int {
	log, err := createLogger()
	if err != nil {
		fmt.Println("error setting up the logger:", err)
		return 1
	}
	log = log.With(zap.String("release", version.HumanVersion))
	defer func() {
		// If we cannot sync, there's probably something wrong with outputting logs,
		// so we probably cannot write using fmt.Println either.
		// Hence, ignoring the error for now.
		_ = log.Sync()
	}()

	flag.Parse()

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(ginzap.Ginzap(log, time.RFC3339, true))
	r.Use(gin.Recovery())

	s, err := server.New(context.Background(), log, r, server.WithServerAddress(*httpPort))
	if err != nil {
		log.Error("error creating server", zap.Error(err))
		return 1
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		log.Info("starting server on port", zap.String("port", *httpPort))
		if err := s.Serve(); err != nil {
			log.Info("error starting server", zap.Error(err))
			return err
		}
		return nil
	})

	<-ctx.Done()

	eg.Go(func() error {
		log.Info("gracefully stopping server")
		err = s.Stop()
		if err != nil {
			return err
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		return 1
	}
	return 0
}

func createLogger() (*zap.Logger, error) {
	encoder := zap.NewProductionEncoderConfig()
	level := zap.NewAtomicLevelAt(getLevelLogger(config.LogLevel))

	zapConfig := zap.NewProductionConfig()
	zapConfig.EncoderConfig = encoder
	zapConfig.Level = level
	zapConfig.OutputPaths = []string{"stdout"}
	zapConfig.ErrorOutputPaths = []string{"stderr"}
	return zapConfig.Build()
}

func getLevelLogger(level string) zapcore.Level {
	if level == "debug" {
		return zap.DebugLevel
	}
	return zap.InfoLevel
}
