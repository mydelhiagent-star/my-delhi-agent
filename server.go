package main

import (
	"context"
	"log"
	"myapp/config"
	"myapp/database"
	"myapp/handlers"
	"myapp/routes"
	"myapp/services"
	"net/http"
	

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.LoadConfig()
	client := database.ConnectMongo(cfg.MongoURI)
	userCollection := client.Database(cfg.MongoDB).Collection("users")
	leadCollection := client.Database(cfg.MongoDB).Collection("leads")
	propertyCollection := client.Database(cfg.MongoDB).Collection("property")

	ctx := context.Background()
	services.InitializeB2Service(ctx)

	authService := &services.AuthService{
		UserCollection: userCollection,
		JWTSecret:        cfg.JWTSecret,
	}
	authHandler := &handlers.AuthHandler{Service: authService}

	leadService := &services.LeadService{
		LeadCollection: leadCollection,
	}
	leadHandler := &handlers.LeadHandler{Service: leadService}

	propertyService := &services.PropertyService{
		PropertyCollection: propertyCollection,
	}
	

	propertyHandler := &handlers.PropertyHandler{Service: propertyService}



	r := mux.NewRouter()
	routes.RegisterAuthRoutes(r, authHandler)
	routes.RegisterLeadRoutes(r,leadHandler,cfg.JWTSecret)
	routes.RegisterPropertyRoutes(r,propertyHandler,cfg.JWTSecret)

	log.Printf("🚀 Server running on port %s\n", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
