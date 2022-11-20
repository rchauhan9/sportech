package players_test

import (
	"context"
	"github.com/go-kit/log"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rchauhan9/sportech/database"
	"github.com/rchauhan9/sportech/players"
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
	repository players.Repository
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

	repository := players.NewRepository(dbPool)
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
	suite.cleaner.Acquire("persons", "team_players")
}

func (suite *RepositoryTestSuite) TearDownTest() {
	suite.cleaner.Clean("persons", "team_players")
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

func (suite *RepositoryTestSuite) TestListPlayers() {
	salah := createTeamPlayer(suite, uuid.New().String(), uuid.New().String(), 11, "FWD", "RW", time.Date(2017, time.July, 1, 0, 0, 0, 0, time.UTC), nil)
	alisson := createTeamPlayer(suite, uuid.New().String(), uuid.New().String(), 1, "GK", "GK", time.Date(2016, time.July, 1, 0, 0, 0, 0, time.UTC), nil)

	playerDBs, err := suite.repository.ListPlayers(suite.ctx)
	require.NoError(suite.T(), err)

	require.Equal(suite.T(), 2, len(playerDBs))

	expecteds := []players.PlayerDB{alisson, salah}
	for idx, exp := range expecteds {
		require.Equal(suite.T(), exp.ID, playerDBs[idx].ID)
		require.Equal(suite.T(), exp.TeamID, playerDBs[idx].TeamID)
		require.Equal(suite.T(), exp.PersonID, playerDBs[idx].PersonID)
		require.Equal(suite.T(), exp.SquadNumber, playerDBs[idx].SquadNumber)
		require.Equal(suite.T(), exp.GeneralPosition, playerDBs[idx].GeneralPosition)
		require.Equal(suite.T(), exp.SpecificPosition, playerDBs[idx].SpecificPosition)
		require.Equal(suite.T(), exp.Started, playerDBs[idx].Started)
		require.Equal(suite.T(), exp.Ended, playerDBs[idx].Ended)
	}
}

func (suite *RepositoryTestSuite) TestGetPlayer() {
	expected := createTeamPlayer(suite, uuid.New().String(), uuid.New().String(), 11, "FWD", "RW", time.Date(2016, time.July, 1, 0, 0, 0, 0, time.UTC), nil)

	result, err := suite.repository.GetPlayer(suite.ctx, expected.ID)
	require.NoError(suite.T(), err)

	require.Equal(suite.T(), expected.ID, result.ID)
	require.Equal(suite.T(), expected.TeamID, result.TeamID)
	require.Equal(suite.T(), expected.PersonID, result.PersonID)
	require.Equal(suite.T(), expected.SquadNumber, result.SquadNumber)
	require.Equal(suite.T(), expected.GeneralPosition, result.GeneralPosition)
	require.Equal(suite.T(), expected.SpecificPosition, result.SpecificPosition)
	require.Equal(suite.T(), expected.Started, result.Started)
	require.Equal(suite.T(), expected.Ended, result.Ended)
}

func createTeamPlayer(suite *RepositoryTestSuite, personID string, teamID string, squadNumber int, generalPosition string, specificPosition string, started time.Time, ended *time.Time) players.PlayerDB {
	query := `
	    INSERT INTO team_players (person_id, team_id, squad_number, general_position, specific_position, started, ended)
	    VALUES
	    ($1, $2, $3, $4, $5, $6, $7)
	    RETURNING id, person_id, team_id, squad_number, general_position, specific_position, started, ended
	`
	row := suite.dbPool.QueryRow(suite.ctx, query, personID, teamID, squadNumber, generalPosition, specificPosition, started, ended)
	var pDB players.PlayerDB
	err := row.Scan(&pDB.ID, &pDB.PersonID, &pDB.TeamID, &pDB.SquadNumber, &pDB.GeneralPosition, &pDB.SpecificPosition, &pDB.Started, &pDB.Ended)
	require.NoError(suite.T(), err)
	return pDB
}
