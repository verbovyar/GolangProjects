package domain

import (
	"github.com/stretchr/testify/require"
	"testing"
)

type playerDataTest struct {
	testName string

	name        string
	club        string
	nationality string
}

func TestCreateNewPlayer(t *testing.T) {
	player := &playerDataTest{
		name:        "Verbov",
		club:        "Barcelona",
		nationality: "Russian",
	}

	_, err := NewPlayer(player.name, player.club, player.nationality)
	require.NoError(t, err)
}

func TestCreateNewPlayerError(t *testing.T) {
	var cases [3]playerDataTest

	cases[0] = playerDataTest{
		testName: "bad name",

		name:        "Verbovvvvvvvvvvvvvvvvvvvvvvvvvvvvv",
		club:        "Barcelona",
		nationality: "Russian",
	}
	cases[1] = playerDataTest{
		testName: "bad club",

		name:        "Verbov",
		club:        "Barcelonaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		nationality: "Russian",
	}
	cases[2] = playerDataTest{
		testName: "bad nationality",

		name:        "Verbov",
		club:        "Barcelona",
		nationality: "",
	}

	for _, testCases := range cases {
		t.Run(testCases.testName, func(t *testing.T) {
			_, err := NewPlayer(testCases.name, testCases.club, testCases.nationality)
			require.ErrorContains(t, err, testCases.testName)
		})
	}
}

func TestPlayer_GetName(t *testing.T) {
	player := Player{
		Name:        "Verbov",
		Club:        "Barcelona",
		Nationality: "Russian",
		Id:          12,
	}

	name := player.GetName()
	require.Equal(t, player.Name, name)
}

func TestPlayer_GetClub(t *testing.T) {
	player := Player{
		Name:        "Verbov",
		Club:        "Barcelona",
		Nationality: "Russian",
		Id:          12,
	}

	club := player.GetClub()
	require.Equal(t, player.Club, club)
}

func TestPlayer_GetNationality(t *testing.T) {
	player := Player{
		Name:        "Verbov",
		Club:        "Barcelona",
		Nationality: "Russian",
		Id:          12,
	}

	nationality := player.GetNationality()
	require.Equal(t, player.Nationality, nationality)
}

func TestPlayer_GetId(t *testing.T) {
	player := Player{
		Name:        "Verbov",
		Club:        "Barcelona",
		Nationality: "Russian",
		Id:          12,
	}

	id := player.GetId()
	require.Equal(t, player.Id, id)
}
