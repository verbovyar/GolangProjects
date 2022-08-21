package db

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/internal/repositories/interfaces"
	"github.com/go-telegram-bot-api/telegram-bot-api/pkg/postgres"
	"github.com/magiconair/properties/assert"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
)

type Fixture struct {
	repository interfaces.Repository
	pool       postgres.PoolInterface
	mock       pgxmock.PgxPoolIface
}

func setUp(t *testing.T) Fixture {
	pool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal("pool error")
	}

	playerRepo := PlayersRepository{
		Pool: pool,
	}
	return Fixture{
		repository: &playerRepo,
		pool:       pool,
		mock:       pool,
	}
}

func (f *Fixture) tearDown() {
	f.pool.Close()
}

func TestPlayersRepository_Delete(t *testing.T) {

	testTable := []struct {
		name            string
		id              uint
		expectedRequest error
	}{
		{
			name:            "OK",
			id:              6,
			expectedRequest: nil,
		},
	}

	f := setUp(t)
	defer f.tearDown()

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			f.mock.ExpectQuery(regexp.QuoteMeta("DELETE FROM Players WHERE id = $1")).WithArgs(testCase.id).WillReturnRows()
			err := f.repository.Delete(testCase.id)
			require.NoError(t, err)
		})
	}
}

func TestPlayersRepository_List(t *testing.T) {
	f := setUp(t)
	defer f.tearDown()

	query := "SELECT id, name, clubName, nationalityName FROM Players JOIN Clubs ON Players.club_id = Clubs.clubId JOIN Nationalities ON Nationalities.nationalityId = Players.nationality_id"

	row := pgxmock.NewRows([]string{"id", "name", "clubName", "nationalityName"}).AddRow(1, "Verbov", "Barcelona", "Russian")

	f.mock.ExpectQuery(query).WillReturnRows(row)

	if err := f.mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	resp := f.repository.List()
	assert.Equal(t, resp[0].Name, "Verbov")
	assert.Equal(t, resp[0].Club, "Barcelona")
	assert.Equal(t, resp[0].Nationality, "Russian")
	assert.Equal(t, resp[0].Id, 1)
}
