package teams_test

import (
	"context"
	"github.com/go-kit/log"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rchauhan9/sportech/database"
	"github.com/rchauhan9/sportech/teams"
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
	repository teams.Repository
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

	repository := teams.NewRepository(dbPool)
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
	suite.cleaner.Acquire("teams")
}

func (suite *RepositoryTestSuite) TearDownTest() {
	suite.cleaner.Clean("teams")
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

func (suite *RepositoryTestSuite) TestListTeams() {
	liverpoolCity := "Liverpool"
	liverpool := createTeam(suite, "Liverpool Football Club", "Liverpool", "LFC", nil, 1892, &liverpoolCity, uuid.NewString(), uuid.NewString(), uuid.NewString())
	arsenalCity := "London"
	arsenalNickname := "The Gunners"
	arsenal := createTeam(suite, "Arsenal Football Club", "Arsenal", "AFC", &arsenalNickname, 1882, &arsenalCity, uuid.NewString(), uuid.NewString(), uuid.NewString())

	tms, err := suite.repository.ListTeams(suite.ctx)
	require.NoError(suite.T(), err)

	require.Equal(suite.T(), 2, len(tms))

	expecteds := []teams.Team{liverpool, arsenal}
	for idx, exp := range expecteds {
		require.Equal(suite.T(), exp.ID, tms[idx].ID)
		require.Equal(suite.T(), exp.FullName, tms[idx].FullName)
		require.Equal(suite.T(), exp.MediumName, tms[idx].MediumName)
		require.Equal(suite.T(), exp.Acronym, tms[idx].Acronym)
		require.Equal(suite.T(), exp.Nickname, tms[idx].Nickname)
		require.Equal(suite.T(), exp.YearFounded, tms[idx].YearFounded)
		require.Equal(suite.T(), exp.City, tms[idx].City)
		require.Equal(suite.T(), exp.Country, tms[idx].Country)
		require.Equal(suite.T(), exp.Stadium, tms[idx].Stadium)
		require.Equal(suite.T(), exp.League, tms[idx].League)
	}
}

func (suite *RepositoryTestSuite) TestGetTeam() {
	liverpoolCity := "Liverpool"
	liverpool := createTeam(suite, "Liverpool Football Club", "Liverpool", "LFC", nil, 1892, &liverpoolCity, uuid.NewString(), uuid.NewString(), uuid.NewString())

	result, err := suite.repository.GetTeam(suite.ctx, liverpool.ID)
	require.NoError(suite.T(), err)

	require.Equal(suite.T(), liverpool.ID, result.ID)
	require.Equal(suite.T(), liverpool.FullName, result.FullName)
	require.Equal(suite.T(), liverpool.MediumName, result.MediumName)
	require.Equal(suite.T(), liverpool.Acronym, result.Acronym)
	require.Equal(suite.T(), liverpool.Nickname, result.Nickname)
	require.Equal(suite.T(), liverpool.YearFounded, result.YearFounded)
	require.Equal(suite.T(), liverpool.City, result.City)
	require.Equal(suite.T(), liverpool.Country, result.Country)
	require.Equal(suite.T(), liverpool.Stadium, result.Stadium)
	require.Equal(suite.T(), liverpool.League, result.League)
}

func createTeam(suite *RepositoryTestSuite, fullName string, mediumName string, acronym string, nickname *string, yearFounded int, city *string, country string, stadium string, league string) teams.Team {
	query := `
	    INSERT INTO teams (full_name, medium_name, acronym, nickname, year_founded, city, country_id, stadium_id, league_id)
	    VALUES
	    ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	    RETURNING id, full_name, medium_name, acronym, nickname, year_founded, city, country_id, stadium_id, league_id
	`
	row := suite.dbPool.QueryRow(suite.ctx, query, fullName, mediumName, acronym, nickname, yearFounded, city, country, stadium, league)
	var team teams.Team
	err := row.Scan(&team.ID, &team.FullName, &team.MediumName, &team.Acronym, &team.Nickname, &team.YearFounded, &team.City, &team.Country, &team.Stadium, &team.League)
	require.NoError(suite.T(), err)
	return team
}
