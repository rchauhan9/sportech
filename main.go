package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/pkg/errors"
	"github.com/rchauhan9/sportech/commons/go/configutil"
	"github.com/rchauhan9/sportech/config"
	"github.com/rchauhan9/sportech/database"
	"github.com/rchauhan9/sportech/middleware"
	"github.com/rchauhan9/sportech/team"
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
	migrationPath := flag.String("migration-dir", "./migrations", "Directory containing migrations")
	flag.Parse()

	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))

	var conf *config.Config
	err := configutil.LoadConfig(*configPath, logger, &conf)
	if err != nil {
		return 1
	}

	migrator, err := database.NewMigrator(conf.Database.URL, *migrationPath, logger)
	if err != nil {
		panic(errors.Wrap(err, "unable to create migrator"))
	}
	if err = migrator.MigrateDb(); err != nil {
		panic(errors.Wrap(err, "unable to migrate database"))
	}
	if err = migrator.Close(); err != nil {
		level.Error(logger).Log("err", errors.Wrap(err, "error closing migrator"))
	}

	db, err := database.NewDatabasePool(ctx, conf.Database.URL)
	if err != nil {
		level.Error(logger).Log("err", err)
		return 1
	}
	defer db.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/health", health)

	teamRepository := team.NewRepository(db)
	teamService := team.NewService(teamRepository)
	listTeamsEndpoint := team.MakeListTeamsEndpoint(teamService)
	listTeamsEndpoint = middleware.AddLogging(listTeamsEndpoint, logger)
	getTeamEndpoint := team.MakeGetTeamEndpoint(teamService)
	getTeamEndpoint = middleware.AddLogging(getTeamEndpoint, logger)
	teamHandler := team.MakeHandler(listTeamsEndpoint, getTeamEndpoint)
	mux.Handle("/api/v1/teams/", teamHandler)

	baseHTTPServer := http.Server{
		Addr:    ":" + conf.Port,
		Handler: accessControl(mux),
	}

	defer func() {
		if err := baseHTTPServer.Shutdown(ctx); err != nil {
			level.Error(logger).Log("err", errors.Wrap(err, "error shutting down http server"))
		}
	}()

	errs := make(chan error, 1)
	go func() {
		level.Info(logger).Log("transport", "http", "address", conf.Port, "msg", "listening")
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

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}

func main() {
	os.Exit(realMain())
}
