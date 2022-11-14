package managers_test

import (
	"context"
	"github.com/go-kit/log"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rchauhan9/sportech/database"
	"github.com/rchauhan9/sportech/managers"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gopkg.in/khaiql/dbcleaner.v2"
	"gopkg.in/khaiql/dbcleaner.v2/engine"
	"testing"
	"time"
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
	repository managers.Repository
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

	repository := managers.NewRepository(dbPool)
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
	suite.cleaner.Acquire("persons", "team_manager")
}

func (suite *RepositoryTestSuite) TearDownTest() {
	suite.cleaner.Clean("persons", "team_manager")
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

func (suite *RepositoryTestSuite) TestListManagers() {
	guardiola := createPerson(suite, "Pep", nil, "Guardiola", time.Date(1971, time.January, 18, 0, 0, 0, 0, time.UTC), uuid.New().String())
	kloppMiddleNames := "Norbert"
	klopp := createPerson(suite, "Jürgen", &kloppMiddleNames, "Guardiola", time.Date(1967, time.June, 16, 0, 0, 0, 0, time.UTC), uuid.New().String())

	guardiolaManager := createTeamManager(suite, guardiola, uuid.New().String(), time.Date(2017, time.July, 1, 0, 0, 0, 0, time.UTC), nil)
	kloppManager := createTeamManager(suite, guardiola, uuid.New().String(), time.Date(2017, time.July, 1, 0, 0, 0, 0, time.UTC), nil)

	managers, err := suite.repository.ListManagers(suite.ctx)
	require.NoError(suite.T(), err)

	require.Equal(suite.T(), 2, len(managers))

	require.NotNil(suite.T(), managers[0].ID)
	require.Equal(suite.T(), "La Liga", managers[0].FirstName)
	require.Equal(suite.T(), int32(20), managers[0].NumberOfTeams)
	require.Equal(suite.T(), "/countries/"+spain, managers[0].Country)

	require.NotNil(suite.T(), managers[1].ID)
	require.Equal(suite.T(), "Premier League", managers[1].Name)
	require.Equal(suite.T(), int32(20), managers[1].NumberOfTeams)
	require.Equal(suite.T(), "/countries/"+england, managers[1].Country)
}

func (suite *RepositoryTestSuite) TestGetLeague() {
	england := createCountry(suite, "England")
	premierLeagueID := createLeague(suite, "Premier League", 20, england)

	league, err := suite.repository.GetLeague(suite.ctx, premierLeagueID)
	require.NoError(suite.T(), err)

	require.Equal(suite.T(), premierLeagueID, league.ID)
	require.Equal(suite.T(), "Premier League", league.Name)
	require.Equal(suite.T(), int32(20), league.NumberOfTeams)
	require.Equal(suite.T(), "/countries/"+england, league.Country)
}

func createPerson(suite *RepositoryTestSuite, firstName string, middleNames *string, lastName string, dateOfBirth time.Time, nationality string) string {
	query := `
	    INSERT INTO persons (first_name, middle_names, last_name, date_of_birth, nationality)
	    VALUES
	    ($1, $2, $3, $4, $5)
	    RETURNING id
	`
	row := suite.dbPool.QueryRow(suite.ctx, query, firstName, middleNames, lastName, dateOfBirth, nationality)
	var personID string
	err := row.Scan(&personID)
	require.NoError(suite.T(), err)
	return personID
}

func createTeamManager(suite *RepositoryTestSuite, personID string, teamID string, started time.Time, ended *time.Time) string {
	query := `
	    INSERT INTO team_manager (person, team, started, ended)
	    VALUES
	    ($1, $2, $3, $4)
	    RETURNING id
	`
	row := suite.dbPool.QueryRow(suite.ctx, query, personID, teamID, started, ended)
	var teamManagerID string
	err := row.Scan(&teamManagerID)
	require.NoError(suite.T(), err)
	return teamManagerID
}
