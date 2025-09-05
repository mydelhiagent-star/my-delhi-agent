package validate

import (
	"errors"
	"regexp"
)

func ValidatePhone(phone string) error {
	re := regexp.MustCompile(`^[0-9]{10}$`)
	if !re.MatchString(phone) {
		return errors.New("invalid phone number")
	}
	return nil
}

func ValidateEmail(email string) error {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !re.MatchString(email) {
		return errors.New("invalid email")
	}
	return nil
}

