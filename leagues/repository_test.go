package leagues_test

import (
	"context"
	"github.com/go-kit/log"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rchauhan9/sportech/database"
	"github.com/rchauhan9/sportech/leagues"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gopkg.in/khaiql/dbcleaner.v2"
	"gopkg.in/khaiql/dbcleaner.v2/engine"
	"testing"
)

const (
	postgresURL   = "postgres://postgres:password@localhost:5432/sportech_testing?sslmode=disable"
	migrationPath = "../migrations"
)

type RepositoryTestSuite struct {
	suite.Suite
	ctx        context.Context
	migrator   *database.Migrator
	cleaner    dbcleaner.DbCleaner
	repository leagues.Repository
	dbPool     *pgxpool.Pool
}

func (suite *RepositoryTestSuite) SetupSuite() {
	ctx := context.Background()
	suite.ctx = ctx

	migrator, err := database.NewMigrator(postgresURL, migrationPath, log.NewNopLogger())
	require.NoError(suite.T(), err)
	suite.migrator = migrator

	err = migrator.PurgeDB()
	require.NoError(suite.T(), err)

	err = migrator.MigrateDb()
	require.NoError(suite.T(), err)

	dbPool, err := database.NewDatabasePool(ctx, postgresURL)
	require.NoError(suite.T(), err)
	suite.dbPool = dbPool

	repository := leagues.NewRepository(dbPool)
	suite.repository = repository

	cleaner := dbcleaner.New()
	cleaner.SetEngine(engine.NewPostgresEngine(postgresURL))
	suite.cleaner = cleaner
}

func (suite *RepositoryTestSuite) TearDownSuite() {
	suite.dbPool.Close()
	suite.cleaner.Close()
	suite.migrator.Close()
}

func (suite *RepositoryTestSuite) SetupTest() {
	suite.cleaner.Acquire("countries", "leagues")
}

func (suite *RepositoryTestSuite) TearDownTest() {
	suite.cleaner.Clean("countries", "leagues")
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

func (suite *RepositoryTestSuite) TestListLeagues() {
	england := uuid.New().String()
	spain := uuid.New().String()
	_ = createLeague(suite, "Premier League", 20, england)
	_ = createLeague(suite, "La Liga", 20, spain)

	leagues, err := suite.repository.ListLeagues(suite.ctx)
	require.NoError(suite.T(), err)

	require.Equal(suite.T(), 2, len(leagues))

	require.NotNil(suite.T(), leagues[0].ID)
	require.Equal(suite.T(), "La Liga", leagues[0].Name)
	require.Equal(suite.T(), int32(20), leagues[0].NumberOfTeams)
	require.Equal(suite.T(), "/countries/"+spain, leagues[0].Country)

	require.NotNil(suite.T(), leagues[1].ID)
	require.Equal(suite.T(), "Premier League", leagues[1].Name)
	require.Equal(suite.T(), int32(20), leagues[1].NumberOfTeams)
	require.Equal(suite.T(), "/countries/"+england, leagues[1].Country)
}

func (suite *RepositoryTestSuite) TestGetLeague() {
	england := uuid.New().String()
	premierLeagueID := createLeague(suite, "Premier League", 20, england)

	league, err := suite.repository.GetLeague(suite.ctx, premierLeagueID)
	require.NoError(suite.T(), err)

	require.Equal(suite.T(), premierLeagueID, league.ID)
	require.Equal(suite.T(), "Premier League", league.Name)
	require.Equal(suite.T(), int32(20), league.NumberOfTeams)
	require.Equal(suite.T(), "/countries/"+england, league.Country)
}

func createLeague(suite *RepositoryTestSuite, name string, numberOfTeams int32, countryID string) string {
	query := `
	    INSERT INTO leagues (name, number_of_teams, country)
	    VALUES
	    ($1, $2, $3)
	    RETURNING id
	`
	row := suite.dbPool.QueryRow(suite.ctx, query, name, numberOfTeams, countryID)
	var leagueID string
	err := row.Scan(&leagueID)
	require.NoError(suite.T(), err)
	return leagueID
}
