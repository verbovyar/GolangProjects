package interfaces

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/domain"
)

type Repository interface {
	Add(player *domain.Player) error
	List() []*domain.Player
	Update(user *domain.Player, id uint) error
	Delete(id uint) error
}
