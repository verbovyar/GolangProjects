package botService

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/internal/domain"
	"github.com/go-telegram-bot-api/telegram-bot-api/internal/repositories/interfaces"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

const (
	helpCmd   = "help"
	listCmd   = "list"
	addCmd    = "add"
	changeCmd = "change"
	deleteCmd = "delete"
)

var Service *PlayersService

func New(repository interfaces.Repository) *PlayersService {
	Service = &PlayersService{
		repository: repository,
	}
	return Service
}

type PlayersService struct {
	repository interfaces.Repository
}

//---------

func listFunc(string) string {
	data := Service.repository.List()
	res := make([]string, 0, len(data))

	for _, value := range data {
		res = append(res, value.GetData())
	}

	return strings.Join(res, "\n")
}

func helpFunc(string) string {
	return "/help - list commands\n" +
		"/list - list data\n" +
		"/add <name> <club> <nationality> - add new player \n" +
		"/change <id> <name> <club> <nationality> - change info about person\n" +
		"/delete <id> - delete person"
}

func addFunc(data string) string {
	params := strings.Split(data, " ")
	if len(params) != 3 {
		return errors.Wrapf(errors.New("bad argument"), "%d items: <%v>", len(params), params).Error()
	}

	user, err := domain.NewPlayer(params[0], params[1], params[2])
	if err != nil {
		return err.Error()
	}

	err = Service.repository.Add(user)
	if err != nil {
		return err.Error()
	}

	return fmt.Sprintf("user %v added", user)
}

func changeFunc(data string) string {
	params := strings.Split(data, " ")
	if len(params) != 4 {
		return errors.Wrapf(errors.New("bad argument"), "%d items: <%v>", len(params), params).Error()
	}

	user, err := domain.NewPlayer(params[1], params[2], params[3])
	domain.LastId--
	if err != nil {
		return err.Error()
	}

	id, _ := strconv.Atoi(params[0])
	err = Service.repository.Update(user, uint(id))
	if err != nil {
		return err.Error()
	}

	return fmt.Sprintf("user info was changed")
}

func deleteFunc(data string) string {
	params := strings.Split(data, "\n")
	if len(params) != 1 {
		return errors.Wrapf(errors.New("bad argument"), "%d items: <%v>", len(params), params).Error()
	}

	id, _ := strconv.Atoi(params[0])

	err := Service.repository.Delete(uint(id))
	if err != nil {
		return err.Error()
	}

	return fmt.Sprintf("user deleted")
}

func AddHandlers(bot *Bot) {
	bot.RegisterHandler(helpCmd, helpFunc)
	bot.RegisterHandler(listCmd, listFunc)
	bot.RegisterHandler(addCmd, addFunc)
	bot.RegisterHandler(changeCmd, changeFunc)
	bot.RegisterHandler(deleteCmd, deleteFunc)
}
