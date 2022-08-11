package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/pkg/errors"
	"github.com/rchauhan9/sportech/commons/go/configutil"
	"github.com/rchauhan9/sportech/config"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok\n")
}

func realMain() int {
	ctx := context.Background()

	configPath := flag.String("config-dir", "./config", "Directory containing config.yml")
	flag.Parse()

	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))

	var conf *config.Config
	err := configutil.LoadConfig(*configPath, logger, &conf)
	if err != nil {
		return 1
	}

	baseHTTPServer := http.Server{
		Addr:    conf.Server.HTTPAddress,
		Handler: nil,
	}
	defer func() {
		if err := baseHTTPServer.Shutdown(ctx); err != nil {
			level.Error(logger).Log("err", errors.Wrap(err, "error shutting down http server"))
		}
	}()

	http.HandleFunc("/health", health)

	errs := make(chan error, 1)
	go func() {
		level.Info(logger).Log("transport", "http", "address", conf.Server.HTTPAddress, "msg", "listening")
		if err := baseHTTPServer.ListenAndServe(); err != nil {
			errs <- err
		}
	}()
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- errors.New((<-c).String())
	}()
	logger.Log("terminated", <-errs)

	return 0
}

func main() {
	os.Exit(realMain())
}
