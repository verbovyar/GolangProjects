package grpcHandlers

import (
	"errors"
	"unicode"
)

const (
	nameError        string = "bad name argument: "
	clubError        string = "bad club argument: "
	nationalityError string = "bad nationality name argument: "
	idError          string = "bad id argument: "
)

const maxSize = 17

func validateId(id int32) error {
	if id <= 0 {
		return errors.New(idError + "id must be greater than zero")
	}

	return nil
}

func validateName(name string) error {
	nameSize := len(name)
	for i := 0; i < nameSize; i++ {
		if !unicode.IsLetter(rune(name[i])) {
			return errors.New(clubError + "the string must contain only Latin letters")
		}
	}
	if nameSize == 0 || nameSize > maxSize {
		return errors.New(nameError + "name length must be less 17 symbols or not equal 0 symbols")
	}
	if !unicode.IsUpper(rune(name[0])) {
		return errors.New(nameError + "the name must begin with a capital letter")
	}

	return nil
}

func validateClub(club string) error {
	clubSize := len(club)
	for i := 0; i < clubSize; i++ {
		if !unicode.IsLetter(rune(club[i])) {
			return errors.New(clubError + "the string must contain only Latin letters")
		}
	}
	if clubSize == 0 || clubSize > maxSize {
		return errors.New(clubError + "club length must be less 17 symbols or not equal 0 symbols")
	}
	if !unicode.IsUpper(rune(club[0])) {
		return errors.New(clubError + "the club must begin with a capital letter")
	}

	return nil
}

func validateNationality(nationality string) error {
	nationalitySize := len(nationality)
	for i := 0; i < nationalitySize; i++ {
		if !unicode.IsLetter(rune(nationality[i])) {
			return errors.New(clubError + "the string must contain only Latin letters")
		}
	}
	if nationalitySize == 0 || nationalitySize > maxSize {
		return errors.New(clubError + "nationality length must be less 17 symbols or not equal 0 symbols")
	}
	if !unicode.IsUpper(rune(nationality[0])) {
		return errors.New(nationalityError + "the nationality must begin with a capital letter")
	}

	return nil
}

func validateUpdateRequest(name, club, nationality string, id int32) error {
	err := validateId(id)
	if err != nil {
		return err
	}

	err = validateName(name)
	if err != nil {
		return err
	}

	err = validateClub(club)
	if err != nil {
		return err
	}

	err = validateNationality(nationality)
	if err != nil {
		return err
	}

	return nil
}

func validateAddRequest(name, club, nationality string) error {
	err := validateName(name)
	if err != nil {
		return err
	}

	err = validateClub(club)
	if err != nil {
		return err
	}

	err = validateNationality(nationality)
	if err != nil {
		return err
	}

	return nil
}

func validateDeleteRequest(id int32) error {
	err := validateId(id)
	if err != nil {
		return err
	}

	return nil
}
