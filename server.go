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
	h "github.com/gorilla/handlers"
)

func main() {
	cfg := config.LoadConfig()
	client := database.ConnectMongo(cfg.MongoURI)
	dealerCollection := client.Database(cfg.MongoDB).Collection("dealers")
	leadCollection := client.Database(cfg.MongoDB).Collection("leads")
	propertyCollection := client.Database(cfg.MongoDB).Collection("property")

	ctx := context.Background()
	services.InitializeB2Service(ctx)

	dealerService := &services.DealerService{
		DealerCollection: dealerCollection,
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
	routes.RegisterAuthRoutes(r, dealerHandler)
	routes.RegisterLeadRoutes(r,leadHandler,cfg.JWTSecret)
	routes.RegisterPropertyRoutes(r,propertyHandler,cfg.JWTSecret)

	 corsHandler := h.CORS(
		h.AllowedOrigins([]string{"http://localhost:5173"}),
		h.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		h.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)(r)

	log.Printf("🚀 Server running on port %s\n", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, corsHandler))
}
