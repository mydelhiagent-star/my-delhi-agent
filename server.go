package main

import (
	"log"
	"myapp/config"
	"myapp/databases"
	"myapp/handlers"
	"myapp/models"
	"myapp/mongo_repositories"
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
	client := databases.ConnectMongo(cfg.MongoURI)

	redisClient, err := databases.ConnectRedis(cfg.RedisURI, cfg.RedisUsername, cfg.RedisPassword)
	if err != nil {
		log.Printf("⚠️  Redis connection failed: %v", err)
	} else {
		log.Printf("✅ Redis connected successfully")
	}

	r2Service, err := services.NewCloudflareR2Service(cfg.CloudflareAccountID, cfg.CloudflareAccessKeyID, cfg.CloudflareAccessKeySecret, cfg.CloudflareBucketName)
	if err != nil {
		log.Printf("⚠️  R2 service failed: %v", err)

	} else {
		log.Printf("✅ R2 service connected successfully")
	}

	dealerCollection := client.Database(cfg.MongoDB).Collection("dealers")
	leadCollection := client.Database(cfg.MongoDB).Collection("leads")
	propertyCollection := client.Database(cfg.MongoDB).Collection("property")
	// tokenCollection := client.Database(cfg.MongoDB).Collection("token")
	counterCollection := client.Database(cfg.MongoDB).Collection("counters")
	dealerClientCollection := client.Database(cfg.MongoDB).Collection("dealer_clients")
	inquiryCollection := client.Database(cfg.MongoDB).Collection("inquiries")

	// Initialize repositories
	dealerRepo := mongo_repositories.NewMongoDealerRepository(dealerCollection)
	leadRepo := mongo_repositories.NewMongoLeadRepository(leadCollection, propertyCollection)
	propertyRepo := mongo_repositories.NewMongoPropertyRepository(propertyCollection, counterCollection, redisClient)
	// tokenRepo := mongo_repositories.NewMongoTokenRepository(tokenCollection)
	dealerClientRepo := mongo_repositories.NewMongoDealerClientRepository(dealerClientCollection)
	inquiryRepo := mongo_repositories.NewMongoInquiryRepository(inquiryCollection)

	// Initialize services with repositories
	dealerService := &services.DealerService{
		DealerRepo: dealerRepo,
		// TokenRepo:  tokenRepo,
		JWTSecret: cfg.JWTSecret,
	}
	dealerHandler := &handlers.DealerHandler{Service: dealerService}

	leadService := &services.LeadService{
		Repo: leadRepo,
	}

	propertyService := &services.PropertyService{
		Repo:        propertyRepo,
		RedisClient: redisClient,
	}
	dealerClientService := &services.DealerClientService{
		Repo: dealerClientRepo,
		PropertyRepo: propertyRepo,
	}
	inquiryService := services.NewInquiryService(inquiryRepo)

	leadHandler := &handlers.LeadHandler{
		Service:         leadService,
		PropertyService: propertyService,
	}

	propertyHandler := &handlers.PropertyHandler{Service: propertyService, CloudflarePublicURL: cfg.CloudflarePublicURL, DealerService: dealerService}

	dealerClientHandler := &handlers.DealerClientHandler{Service: dealerClientService, CloudflarePublicURL: cfg.CloudflarePublicURL}
	inquiryHandler := handlers.NewInquiryHandler(inquiryService)

	cloudfareHandler := &handlers.CloudfareHandler{
		Service: r2Service,
	}

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
	routes.RegisterCloudFareRoutes(r, cloudfareHandler, cfg.JWTSecret)
	routes.RegisterDealerClientRoutes(r, dealerClientHandler, cfg.JWTSecret)
	routes.SetupInquiryRoutes(r, inquiryHandler, cfg.JWTSecret)

	corsHandler := h.CORS(
		h.AllowedOrigins([]string{"*"}),
		h.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		h.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)(r)

	log.Printf("🚀 Server running on port %s\n", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, corsHandler))
}
