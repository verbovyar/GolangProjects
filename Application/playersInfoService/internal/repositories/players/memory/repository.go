package memory

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/domain"
	"github.com/go-telegram-bot-api/telegram-bot-api/internal/repositories/interfaces"
	"strconv"
	"sync"

	"github.com/pkg/errors"
)

const poolSize = 10

var PlayersMap *InMemoryPlayerRepository

func New() *InMemoryPlayerRepository {
	PlayersMap = &InMemoryPlayerRepository{
		mutex:       sync.RWMutex{},
		mapData:     make(map[uint]*domain.Player),
		poolChannel: make(chan struct{}, poolSize),
	}
	return PlayersMap
}

//--------Methods

type InMemoryPlayerRepository struct { // IRepository
	mutex       sync.RWMutex
	mapData     map[uint]*domain.Player
	poolChannel chan struct{}

	playersInterface interfaces.Repository
}

func (r *InMemoryPlayerRepository) List() []*domain.Player {
	PlayersMap.poolChannel <- struct{}{}
	PlayersMap.mutex.Lock()
	defer func() {
		PlayersMap.mutex.Unlock()
		<-PlayersMap.poolChannel
	}()

	res := make([]*domain.Player, len(PlayersMap.mapData))

	for i, value := range PlayersMap.mapData {
		res[i] = &domain.Player{
			Name:        value.Name,
			Club:        value.Club,
			Nationality: value.Nationality,
			Id:          value.Id}
	}

	return res
}

func (r *InMemoryPlayerRepository) Add(player *domain.Player) error {
	PlayersMap.poolChannel <- struct{}{}
	PlayersMap.mutex.Lock()
	defer func() {
		PlayersMap.mutex.Unlock()
		<-PlayersMap.poolChannel
	}()

	if _, ok := PlayersMap.mapData[player.GetId()]; ok {
		return errors.Wrap(errors.New("user exists"), strconv.FormatUint(uint64(player.GetId()), 10))
	}
	PlayersMap.mapData[player.GetId()] = player

	return nil
}

func (r *InMemoryPlayerRepository) Update(player *domain.Player, id uint) error {
	PlayersMap.poolChannel <- struct{}{}
	PlayersMap.mutex.Lock()
	defer func() {
		PlayersMap.mutex.Unlock()
		<-PlayersMap.poolChannel
	}()

	if _, ok := PlayersMap.mapData[id]; !ok {
		return errors.Wrap(errors.New("user does not exist"), strconv.FormatUint(uint64(player.GetId()), 10))
	}
	PlayersMap.mapData[id] = player
	player.Id = id

	return nil
}

func (r *InMemoryPlayerRepository) Delete(id uint) error {
	PlayersMap.poolChannel <- struct{}{}
	PlayersMap.mutex.Lock()
	defer func() {
		PlayersMap.mutex.Unlock()
		<-PlayersMap.poolChannel
	}()

	if _, ok := PlayersMap.mapData[id]; ok {
		delete(PlayersMap.mapData, id)

		return nil
	}

	return errors.Wrap(errors.New("user does not exist"), strconv.FormatUint(uint64(id), 10))
}
