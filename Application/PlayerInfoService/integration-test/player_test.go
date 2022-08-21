//go:build integration
// +build integration

package test

import (
	"context"
	"github.com/go-telegram-bot-api/telegram-bot-api/api/apiPb"
	"github.com/magiconair/properties/assert"
	"testing"
)

var (
	Id = 25
)

func TestAdd(t *testing.T) {
	in := &apiPb.AddRequest{
		Name:        "Jane",
		Club:        "Psg",
		Nationality: "Russian",
	}

	client := NewClient()

	t.Run("name", func(t *testing.T) {

		response, err := client.Add(context.Background(), in)

		if err != nil {
			assert.Equal(t, response.Id, Id)
		}
	})
}

func TestUpdate(t *testing.T) {
	in := &apiPb.UpdateRequest{
		Name:        "Jane",
		Club:        "Barcelona",
		Nationality: "Russian",
		Id:          int32(Id),
	}

	client := NewClient()

	t.Run("name", func(t *testing.T) {

		response, err := client.Update(context.Background(), in)

		if err != nil {
			assert.Equal(t, response.Id, Id)
		}
	})
}

func TestList(t *testing.T) {
	in := &apiPb.ListRequest{}

	client := NewClient()

	t.Run("name", func(t *testing.T) {

		response, err := client.List(context.Background(), in)

		if err != nil {
			assert.Equal(t, response.Players[0].Name, "Mitrovich")
		}
	})
}

func TestDelete(t *testing.T) {
	in := &apiPb.DeleteRequest{Id: int32(Id)}

	client := NewClient()

	t.Run("name", func(t *testing.T) {

		response, err := client.Delete(context.Background(), in)

		if err != nil {
			assert.Equal(t, response.Result, true)
		}
	})
}
