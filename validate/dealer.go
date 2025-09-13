package validate

import (
	"errors"
	"myapp/mongo_models"
)

func ValidateDealer(dealer models.Dealer) error {
	if dealer.Name == "" || len(dealer.Name) < 2 {
		return errors.New("invalid name")
	}
	if err := ValidatePhone(dealer.Phone); err != nil {
		return errors.New("invalid phone number")
	}
	if dealer.SubLocation == "" || len(dealer.SubLocation) < 2 {
		return errors.New("invalid sublocation")
	}
	if dealer.Location == "" || len(dealer.Location) < 2 {
		return errors.New("invalid location")
	}
	if dealer.Password == "" || len(dealer.Password) < 6 {
		return errors.New("invalid password")
	}
	if dealer.OfficeAddress == "" || len(dealer.OfficeAddress) < 2 {
		return errors.New("invalid office address")
	}
	if dealer.ShopName == "" || len(dealer.ShopName) < 2 {
		return errors.New("invalid shop name")
	}
	if err := ValidateEmail(dealer.Email); err != nil {
		return err
	}

	return nil
}