package stadiums_test

import (
	"context"
	"github.com/go-kit/log"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rchauhan9/sportech/database"
	"github.com/rchauhan9/sportech/stadiums"
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
	repository stadiums.Repository
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

	repository := stadiums.NewRepository(dbPool)
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
	suite.cleaner.Acquire("stadiums")
}

func (suite *RepositoryTestSuite) TearDownTest() {
	suite.cleaner.Clean("stadiums")
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

func (suite *RepositoryTestSuite) TestListStadiums() {
	england := uuid.New().String()
	anfield := createStadium(suite, "Anfield", 54000, "Liverpool", england)
	oldTrafford := createStadium(suite, "Old Trafford", 76000, "Manchester", england)
	emirates := createStadium(suite, "Emirates", 60000, "London", england)

	stads, err := suite.repository.ListStadiums(suite.ctx)
	require.NoError(suite.T(), err)

	require.Equal(suite.T(), 3, len(stads))

	expecteds := []stadiums.Stadium{anfield, emirates, oldTrafford}
	for idx, exp := range expecteds {
		require.Equal(suite.T(), exp.ID, stads[idx].ID)
		require.Equal(suite.T(), exp.Name, stads[idx].Name)
		require.Equal(suite.T(), exp.Capacity, stads[idx].Capacity)
		require.Equal(suite.T(), exp.City, stads[idx].City)
		require.Equal(suite.T(), exp.Country, stads[idx].Country)
	}
}

func (suite *RepositoryTestSuite) TestGetStadium() {
	england := uuid.New().String()
	anfield := createStadium(suite, "Anfield", 54000, "Liverpool", england)

	stad, err := suite.repository.GetStadium(suite.ctx, anfield.ID)
	require.NoError(suite.T(), err)

	require.Equal(suite.T(), stad.ID, anfield.ID)
	require.Equal(suite.T(), stad.Name, anfield.Name)
	require.Equal(suite.T(), stad.Capacity, anfield.Capacity)
	require.Equal(suite.T(), stad.City, anfield.City)
	require.Equal(suite.T(), stad.Country, anfield.Country)
}

func createStadium(suite *RepositoryTestSuite, name string, capacity int, city string, countryID string) stadiums.Stadium {
	query := `
	    INSERT INTO stadiums (name, capacity, city, country_id)
	    VALUES
	    ($1, $2, $3, $4)
	    RETURNING id, name, capacity, city, country_id
	`
	row := suite.dbPool.QueryRow(suite.ctx, query, name, capacity, city, countryID)
	var stadium stadiums.Stadium
	err := row.Scan(&stadium.ID, &stadium.Name, &stadium.Capacity, &stadium.City, &stadium.Country)
	require.NoError(suite.T(), err)
	return stadium
}
