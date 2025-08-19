package main

import (
	"log"
	"myapp/config"
	"myapp/database"
	"myapp/handlers"
	"myapp/models"
	"myapp/response"
	"myapp/routes"
	"myapp/services"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	h "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	cfg := config.LoadConfig()
	client := database.ConnectMongo(cfg.MongoURI)
	
	dealerCollection := client.Database(cfg.MongoDB).Collection("dealers")
	leadCollection := client.Database(cfg.MongoDB).Collection("leads")
	propertyCollection := client.Database(cfg.MongoDB).Collection("property")
	tokenCollection := client.Database(cfg.MongoDB).Collection("token")

	
	

	dealerService := &services.DealerService{
		DealerCollection: dealerCollection,
		TokenCollection: tokenCollection,
		JWTSecret:        cfg.JWTSecret,
	}
	dealerHandler := &handlers.DealerHandler{Service: dealerService}

	leadService := &services.LeadService{
		LeadCollection: leadCollection,
	}
	leadHandler := &handlers.LeadHandler{Service: leadService}

	propertyService := &services.PropertyService{
		PropertyCollection: propertyCollection,
	}

	propertyHandler := &handlers.PropertyHandler{Service: propertyService}

	r := mux.NewRouter()

	r.HandleFunc("/admin/login", func(w http.ResponseWriter, r *http.Request) {
		email := r.FormValue("email")
		password := r.FormValue("password")

		adminEmail := cfg.AdminEmail
		adminPassword := cfg.AdminPassword

		if email != adminEmail || password != adminPassword {
			response.Error(w, http.StatusUnauthorized, "invalid credentials")
			return
		}
		claims := &models.Claims{
			ID:    "1",
			Phone: "9873462385",
			Role:  "admin",
			RegisteredClaims: jwt.RegisteredClaims{
				IssuedAt:  jwt.NewNumericDate(time.Now()), // unique timestamp
				ID:        uuid.New().String(),
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
		if err != nil {
			response.Error(w, http.StatusUnauthorized, err.Error())
			return
		}
		response.JSON(w, http.StatusOK, map[string]string{"token": tokenString})

	}).Methods("POST")

	routes.RegisterDealerRoutes(r, dealerHandler, cfg.JWTSecret)
	routes.RegisterLeadRoutes(r, leadHandler, cfg.JWTSecret)
	routes.RegisterPropertyRoutes(r, propertyHandler, cfg.JWTSecret)

	corsHandler := h.CORS(
		h.AllowedOrigins([]string{"*"}),
		h.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		h.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)(r)

	

	log.Printf("🚀 Server running on port %s\n", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, corsHandler))
}
