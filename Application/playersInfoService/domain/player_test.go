package domain

import (
	"github.com/stretchr/testify/require"
	"testing"
)

type playerDataTest struct {
	errorName string

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
		errorName: "name bad argument",

		name:        "Verbovvvvvvvvvvvvvvvvvvvvvvvvvvvvv",
		club:        "Barcelona",
		nationality: "Russian",
	}
	cases[1] = playerDataTest{
		errorName: "club bad argument",

		name:        "Verbov",
		club:        "Barcelonaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		nationality: "Russian",
	}
	cases[2] = playerDataTest{
		errorName: "nationality bad argument",

		name:        "Verbov",
		club:        "Barcelona",
		nationality: "",
	}

	for _, tCases := range cases {
		t.Run(tCases.errorName, func(t *testing.T) {
			_, err := NewPlayer(tCases.name, tCases.club, tCases.nationality)
			require.NoError(t, err)
		})
	}
}
