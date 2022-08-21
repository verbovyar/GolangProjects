package interfaces

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/internal/domain"
)

//go:generate mockgen -source=repository.go -destination=mock/mock.go

type Repository interface {
	Add(player *domain.Player) error
	List() []*domain.Player
	Update(user *domain.Player, id uint) error
	Delete(id uint) error
}
