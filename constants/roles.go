package constants

const(
	Admin = "admin"
	SuperDealer = "superdealer"
	Dealer = "dealer"
)

var Roles = []string{Admin, SuperDealer, Dealer}



func IsValidRole(role string) bool {
	for _, r := range Roles {
		if r == role {
			return true
		}
	}
	return false
}

