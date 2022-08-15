package domain

import "fmt"

type Player struct {
	Name        string
	Club        string
	Nationality string
	Id          uint
}

var LastId = uint(0)

func NewPlayer(name, club, nationality string) (*Player, error) {
	user := Player{}

	if err := user.SetName(name); err != nil {
		return nil, err
	}

	if err := user.SetClub(club); err != nil {
		return nil, err
	}

	if err := user.SetNationality(nationality); err != nil {
		return nil, err
	}

	LastId++
	user.Id = LastId

	return &user, nil
}

func (user *Player) SetName(name string) error {
	nameSize := len(name)
	if nameSize == 0 || nameSize > 10 {
		return fmt.Errorf("bad name <%v>", name)
	}

	user.Name = name

	return nil
}

func (user *Player) SetClub(club string) error {
	nameSize := len(club)
	if nameSize == 0 || nameSize > 10 {
		return fmt.Errorf("bad club naming <%v>", club)
	}

	user.Club = club

	return nil
}

func (user *Player) SetNationality(nationality string) error {
	nationalitySize := len(nationality)
	if nationalitySize == 0 || nationalitySize > 20 {
		return fmt.Errorf("bad nationality <%v>", nationality)
	}

	user.Nationality = nationality

	return nil
}

func (user Player) GetData() string {
	return fmt.Sprintf("{%d}%s: %s / %s", user.Id, user.Name, user.Club, user.Nationality)
}

func (user Player) GetName() string {
	return user.Name
}

func (user Player) GetClub() string {
	return user.Club
}

func (user Player) GetNationality() string {
	return user.Nationality
}

func (user Player) GetId() uint {
	return user.Id
}
