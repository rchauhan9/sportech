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
	"github.com/rchauhan9/sportech/leagues"
	"github.com/rchauhan9/sportech/managers"
	"github.com/rchauhan9/sportech/middleware"
	"github.com/rchauhan9/sportech/persons"
	"github.com/rchauhan9/sportech/players"
	"github.com/rchauhan9/sportech/stadiums"
	"github.com/rchauhan9/sportech/teams"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok\n")
}

func docs(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "GET /teams\nGET /stadiums\n")
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
	mux.HandleFunc("/", docs)

	teamRepository := teams.NewRepository(db)
	teamService := teams.NewService(teamRepository)
	listTeamsEndpoint := teams.MakeListTeamsEndpoint(teamService)
	listTeamsEndpoint = middleware.AddLogging(listTeamsEndpoint, logger)
	getTeamEndpoint := teams.MakeGetTeamEndpoint(teamService)
	getTeamEndpoint = middleware.AddLogging(getTeamEndpoint, logger)
	teamHandler := teams.MakeHandler(listTeamsEndpoint, getTeamEndpoint)
	mux.Handle("/teams/", teamHandler)

	stadiumRepository := stadiums.NewRepository(db)
	stadiumService := stadiums.NewService(stadiumRepository)
	listStadiumsEndpoint := stadiums.MakeListStadiumsEndpoint(stadiumService)
	listStadiumsEndpoint = middleware.AddLogging(listStadiumsEndpoint, logger)
	getStadiumEndpoint := stadiums.MakeGetStadiumEndpoint(stadiumService)
	getStadiumEndpoint = middleware.AddLogging(getStadiumEndpoint, logger)
	stadiumHandler := stadiums.MakeHandler(listStadiumsEndpoint, getStadiumEndpoint)
	mux.Handle("/stadiums/", stadiumHandler)

	leagueRepository := leagues.NewRepository(db)
	leagueService := leagues.NewService(leagueRepository)
	listLeaguesEndpoint := leagues.MakeListLeaguesEndpoint(leagueService)
	listLeaguesEndpoint = middleware.AddLogging(listLeaguesEndpoint, logger)
	getLeagueEndpoint := leagues.MakeGetLeagueEndpoint(leagueService)
	getLeagueEndpoint = middleware.AddLogging(getLeagueEndpoint, logger)
	leagueHandler := leagues.MakeHandler(listLeaguesEndpoint, getLeagueEndpoint)
	mux.Handle("/leagues/", leagueHandler)

	personsRepository := persons.NewRepository(db)
	personsService := persons.NewService(personsRepository)

	managerRepository := managers.NewRepository(db)
	managerService := managers.NewService(managerRepository, personsService)
	listManagersEndpoint := managers.MakeListManagersEndpoint(managerService)
	listManagersEndpoint = middleware.AddLogging(listManagersEndpoint, logger)
	getManagerEndpoint := managers.MakeGetManagerEndpoint(managerService)
	getManagerEndpoint = middleware.AddLogging(getManagerEndpoint, logger)
	managerHandler := managers.MakeHandler(listManagersEndpoint, getManagerEndpoint)
	mux.Handle("/managers/", managerHandler)

	playerRepository := players.NewRepository(db)
	playerService := players.NewService(playerRepository)
	listPlayersEndpoint := players.MakeListPlayersEndpoint(playerService)
	listPlayersEndpoint = middleware.AddLogging(listPlayersEndpoint, logger)
	getPlayerEndpoint := players.MakeGetPlayerEndpoint(playerService)
	getPlayerEndpoint = middleware.AddLogging(getPlayerEndpoint, logger)
	playerHandler := players.MakeHandler(listPlayersEndpoint, getPlayerEndpoint)
	mux.Handle("/players/", playerHandler)

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
