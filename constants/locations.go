package constants

var Locations = []string{
	"Noida Electronic City", "Noida Sector 62", "Noida Sector 59", "Noida Sector 61",
	"Noida Sector 52", "Noida Sector 34", "Noida City Centre", "Golf Course",
	"Botanical Garden", "Noida Sector 18", "Noida Sector 16", "Noida Sector 15",
	"New Ashok Nagar", "Mayur Vihar Extension", "Mayur Vihar-I", "Mayur Vihar Pocket-I",
	"Akshardham", "Yamuna Bank", "Indraprastha", "Laxmi Nagar", "Preet Vihar", "Nirman Vihar",
	"Karkarduma", "Karkarduma Court", "Anand Vihar", "Kaushambi", "Vaishali",
	"Mandawali – West Vinod Nagar", "IP Extension", "East Vinod Nagar – Mayur Vihar-II",
	"Trilokpuri – Sanjay Lake", "Krishna Nagar", "Noida Sector 51", "Noida Sector 50 (Rainbow)",
	"Noida Sector 76", "Noida Sector 101", "Noida Sector 81", "NSEZ", "Noida Sector 83",
	"Noida Sector 137", "Noida Sector 142", "Noida Sector 143", "Noida Sector 144",
	"Noida Sector 145", "Noida Sector 146", "Noida Sector 147", "Noida Sector 148",
	"Knowledge Park II", "Pari Chowk", "Alpha 1", "Delta 1", "GNIDA Office", "Depot Station",
	"Extra",
}

func IsValidLocation(location string) bool {
	for _, allowed := range Locations {
		if location == allowed {
			return true
		}
	}
	return false
}
