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
	suite.cleaner.Acquire("persons", "team_managers")
}

func (suite *RepositoryTestSuite) TearDownTest() {
	suite.cleaner.Clean("persons", "team_managers")
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

func (suite *RepositoryTestSuite) TestListManagers() {
	guardiola := createTeamManager(suite, uuid.New().String(), uuid.New().String(), time.Date(2017, time.July, 1, 0, 0, 0, 0, time.UTC), nil)
	klopp := createTeamManager(suite, uuid.New().String(), uuid.New().String(), time.Date(2016, time.July, 1, 0, 0, 0, 0, time.UTC), nil)

	managerDBs, err := suite.repository.ListManagers(suite.ctx)
	require.NoError(suite.T(), err)

	require.Equal(suite.T(), 2, len(managerDBs))

	expecteds := []managers.ManagerDB{klopp, guardiola}
	for idx, exp := range expecteds {
		require.Equal(suite.T(), exp.ID, managerDBs[idx].ID)
		require.Equal(suite.T(), exp.TeamID, managerDBs[idx].TeamID)
		require.Equal(suite.T(), exp.PersonID, managerDBs[idx].PersonID)
		require.Equal(suite.T(), exp.Started, managerDBs[idx].Started)
		require.Equal(suite.T(), exp.Ended, managerDBs[idx].Ended)
	}
}

func (suite *RepositoryTestSuite) TestGetManager() {
	expected := createTeamManager(suite, uuid.New().String(), uuid.New().String(), time.Date(2016, time.July, 1, 0, 0, 0, 0, time.UTC), nil)

	result, err := suite.repository.GetManager(suite.ctx, expected.ID)
	require.NoError(suite.T(), err)

	require.Equal(suite.T(), expected.ID, result.ID)
	require.Equal(suite.T(), expected.TeamID, result.TeamID)
	require.Equal(suite.T(), expected.PersonID, result.PersonID)
	require.Equal(suite.T(), expected.Started, result.Started)
	require.Equal(suite.T(), expected.Ended, result.Ended)
}

func createTeamManager(suite *RepositoryTestSuite, personID string, teamID string, started time.Time, ended *time.Time) managers.ManagerDB {
	query := `
	    INSERT INTO team_managers (person_id, team_id, started, ended)
	    VALUES
	    ($1, $2, $3, $4)
	    RETURNING id, person_id, team_id, started, ended
	`
	row := suite.dbPool.QueryRow(suite.ctx, query, personID, teamID, started, ended)
	var mDB managers.ManagerDB
	err := row.Scan(&mDB.ID, &mDB.PersonID, &mDB.TeamID, &mDB.Started, &mDB.Ended)
	require.NoError(suite.T(), err)
	return mDB
}
