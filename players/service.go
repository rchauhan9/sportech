package players

import (
	"context"
	"github.com/pkg/errors"
	"github.com/rchauhan9/sportech/persons"
	"github.com/samber/lo"
)

type Service interface {
	ListPlayers(ctx context.Context) ([]Player, error)
	GetPlayer(ctx context.Context, id string) (Player, error)
}

func NewService(repository Repository, personsService persons.Service) Service {
	return &service{repository: repository, personsService: personsService}
}

type service struct {
	repository     Repository
	personsService persons.Service
}

func (s *service) ListPlayers(ctx context.Context) ([]Player, error) {
	playersDB, err := s.repository.ListPlayers(ctx)
	if err != nil {
		return nil, err
	}

	people, err := s.personsService.ListPersons(ctx)

	personsMap := lo.KeyBy[string, persons.Person](people, func(person persons.Person) string {
		return person.ID
	})

	players := make([]Player, len(playersDB))
	for i := range playersDB {
		players[i] = Player{
			ID:               playersDB[i].ID,
			FirstName:        personsMap[playersDB[i].PersonID].FirstName,
			MiddleNames:      personsMap[playersDB[i].PersonID].MiddleNames,
			LastName:         personsMap[playersDB[i].PersonID].LastName,
			DateOfBirth:      personsMap[playersDB[i].PersonID].DateOfBirth,
			Nationality:      personsMap[playersDB[i].PersonID].Nationality,
			Team:             playersDB[i].TeamID,
			SquadNumber:      playersDB[i].SquadNumber,
			GeneralPosition:  playersDB[i].GeneralPosition,
			SpecificPosition: playersDB[i].SpecificPosition,
			Started:          playersDB[i].Started,
			Ended:            playersDB[i].Ended,
		}
	}

	return players, nil
}

func (s *service) GetPlayer(ctx context.Context, id string) (Player, error) {
	player, err := s.repository.GetPlayer(ctx, id)
	if err != nil {
		return Player{}, err
	}

	person, err := s.personsService.GetPerson(ctx, player.PersonID)
	if err != nil {
		return Player{}, errors.Wrapf(err, "error getting player with id %s", id)
	}

	return Player{
		ID:               player.ID,
		FirstName:        person.FirstName,
		MiddleNames:      person.MiddleNames,
		LastName:         person.LastName,
		DateOfBirth:      person.DateOfBirth,
		Nationality:      person.Nationality,
		Team:             player.TeamID,
		SquadNumber:      player.SquadNumber,
		GeneralPosition:  player.GeneralPosition,
		SpecificPosition: player.SpecificPosition,
		Started:          player.Started,
		Ended:            player.Ended,
	}, nil
}
