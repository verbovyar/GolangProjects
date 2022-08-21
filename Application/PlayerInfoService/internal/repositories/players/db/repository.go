package db

import (
	"context"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/internal/domain"
	"github.com/go-telegram-bot-api/telegram-bot-api/internal/repositories/interfaces"
	"github.com/go-telegram-bot-api/telegram-bot-api/pkg/postgres"
)

func New(client postgres.PoolInterface) *PlayersRepository {
	return &PlayersRepository{
		DbClient: client,
	}
}

type PlayersRepository struct {
	DbClient         postgres.PoolInterface
	playersInterface interfaces.Repository
}

//--------

func (r *PlayersRepository) Add(player *domain.Player) error {
	ctx := context.Background()

	var clubId int
	query := `SELECT clubId FROM clubs WHERE clubName = $1`
	err := r.DbClient.QueryRow(ctx, query, player.Club).Scan(&clubId)
	if err != nil {
		fmt.Printf(" error %v\n", err)
		return err
	}

	var nationalityId int
	query = `SELECT nationalityId FROM nationalities WHERE nationalityName = $1`
	err = r.DbClient.QueryRow(ctx, query, player.Nationality).Scan(&nationalityId)
	if err != nil {
		fmt.Printf("list error %v\n", err)
		return err
	}

	query = `INSERT INTO Players (name, club_id, nationality_id) VALUES ($1, $2, $3)`
	_, err = r.DbClient.Query(ctx, query, player.Name, clubId, nationalityId)

	var id int
	query = `SELECT id FROM Players WHERE name = $1`
	err = r.DbClient.QueryRow(ctx, query, player.Name).Scan(&id)
	if err != nil {
		fmt.Printf("list error %v\n", err)
		return err
	}
	player.Id = uint(id)

	return nil
}
func (r *PlayersRepository) List() []*domain.Player {
	ctx := context.Background()

	query := "SELECT id, name, clubName, nationalityName FROM Players JOIN Clubs ON Players.club_id = Clubs.clubId JOIN Nationalities ON Nationalities.nationalityId = Players.nationality_id"

	rows, err := r.DbClient.Query(ctx, query)
	if err != nil {
		fmt.Printf("list error %v\n", err)
	}

	players := make([]*domain.Player, 0)

	for rows.Next() {
		tempPlayer := domain.Player{}
		err = rows.Scan(&tempPlayer.Id, &tempPlayer.Name, &tempPlayer.Club, &tempPlayer.Nationality)

		players = append(players, &tempPlayer)
	}

	rows.Close()

	return players
}

func (r *PlayersRepository) Update(player *domain.Player, id uint) error {
	ctx := context.Background()

	var clubId int
	query := `SELECT clubId FROM clubs WHERE clubName = $1`
	err := r.DbClient.QueryRow(ctx, query, player.Club).Scan(&clubId)
	if err != nil {
		fmt.Printf("update error %v\n", err)
		return err
	}

	var nationalityId int
	query = `SELECT nationalityId FROM nationalities WHERE nationalityName = $1`
	err = r.DbClient.QueryRow(ctx, query, player.Nationality).Scan(&nationalityId)
	if err != nil {
		fmt.Printf("update error %v\n", err)
		return err
	}

	query = `UPDATE Players SET name = $1, club_id = $2, nationality_id = $3 WHERE id = $4`
	_, err = r.DbClient.Query(ctx, query, player.Name, clubId, nationalityId, id)
	if err != nil {
		fmt.Printf("update error %v\n", err)
		return err
	}

	return nil
}
func (r *PlayersRepository) Delete(id uint) error {

	ctx := context.Background()

	query := `DELETE FROM Players WHERE id = $1`
	_, err := r.DbClient.Query(ctx, query, id)
	if err != nil {
		fmt.Printf("delete error %v\n", err)
		return err
	}

	return nil
}
