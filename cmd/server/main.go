package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/haranton/go-graphql-blog/internals/app"
	"github.com/haranton/go-graphql-blog/internals/config"
	"github.com/haranton/go-graphql-blog/internals/logger"
)

func main() {

	cfg := config.MustLoad()
	logger := logger.GetLogger(cfg.Env)

	application := app.New(cfg, logger)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go application.MustStart()

	<-stop
	application.Close()
	logger.Info("app succesfully stop")

}
