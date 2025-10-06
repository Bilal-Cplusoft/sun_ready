package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Bilal-Cplusoft/sun_ready/internal/client"
	"github.com/Bilal-Cplusoft/sun_ready/internal/database"
	"github.com/Bilal-Cplusoft/sun_ready/internal/handler"
	"github.com/Bilal-Cplusoft/sun_ready/internal/repo"
	"github.com/Bilal-Cplusoft/sun_ready/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/Bilal-Cplusoft/sun_ready/docs"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// @title Sun Ready API
// @version 1.0
// @description API for Sun Ready project management system
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@sunready.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	fmt.Printf("JWT_SECRET: %s\n", jwtSecret)
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	lightFusionURL := os.Getenv("LIGHTFUSION_API")
	if lightFusionURL == "" {
		lightFusionURL = "http://localhost:8085"
	}
	lightFusionAPIKey := os.Getenv("LIGHTFUSION_API_KEY")
	lightFusionEmail := os.Getenv("LIGHTFUSION_EMAIL")
	lightFusionPassword := os.Getenv("LIGHTFUSION_PASSWORD")

	db, err := database.New(databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	userRepo := repo.NewUserRepo(db)
	companyRepo := repo.NewCompanyRepo(db)
	projectRepo := repo.NewProjectRepo(db)
	dealRepo := repo.NewDealRepo(db)

	lightFusionClient := client.NewLightFusionClient(lightFusionURL, lightFusionAPIKey)
	log.Printf("LightFUSION API integration enabled: %s", lightFusionURL)

	if lightFusionEmail != "" && lightFusionPassword != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		token, err := lightFusionClient.Login(ctx, lightFusionEmail, lightFusionPassword)
		if err != nil {
			log.Printf("Warning: Failed to authenticate with LightFusion API: %v", err)
		} else {
			log.Printf("\n Session token obtained: %s \n", token)
		}
	} else {
		log.Println("LightFusion credentials not provided. Set LIGHTFUSION_EMAIL and LIGHTFUSION_PASSWORD in .env")
	}

	authService := service.NewAuthService(userRepo, jwtSecret)
	userService := service.NewUserService(userRepo)
	companyService := service.NewCompanyService(companyRepo)
	projectService := service.NewProjectService(projectRepo)
	dealService := service.NewDealService(dealRepo)

	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	companyHandler := handler.NewCompanyHandler(companyService, userService)
	projectHandler := handler.NewProjectHandler(projectService)
	project3DHandler := handler.NewProject3DHandler(lightFusionClient)
	dealHandler := handler.NewDealHandler(dealService)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Post("/api/auth/register", authHandler.Register)
	r.Post("/api/auth/login", authHandler.Login)

	r.Get("/api/users/{id}", userHandler.GetByID)
	r.Put("/api/users/{id}", userHandler.Update)
	r.Delete("/api/users/{id}", userHandler.Delete)
	r.Get("/api/users", userHandler.List)

	r.Post("/api/companies", companyHandler.Create)
	r.Post("/api/companies/add", companyHandler.AddCompany)
	r.Get("/api/companies/all", companyHandler.FindAll)
	r.Get("/api/companies/slug/{slug}", companyHandler.GetBySlug)
	r.Get("/api/companies/{id}", companyHandler.GetByID)
	r.Put("/api/companies/{id}", companyHandler.Update)
	r.Delete("/api/companies/{id}", companyHandler.Delete)
	r.Get("/api/companies", companyHandler.List)

	r.Post("/api/projects", projectHandler.Create)
	r.Get("/api/projects/{id}", projectHandler.GetByID)
	r.Put("/api/projects/{id}", projectHandler.Update)
	r.Delete("/api/projects/{id}", projectHandler.Delete)
	r.Get("/api/projects", projectHandler.ListByCompany)
	r.Get("/api/projects/user", projectHandler.ListByUser)

	r.Post("/api/deals", dealHandler.Create)
	r.Get("/api/deals/uuid/{uuid}", dealHandler.GetByUUID)
	r.Get("/api/deals/company/{company_id}", dealHandler.ListByCompany)
	r.Get("/api/deals/company/{company_id}/signed", dealHandler.ListSigned)
	r.Get("/api/deals/{id}", dealHandler.GetByID)
	r.Put("/api/deals/{id}", dealHandler.Update)
	r.Delete("/api/deals/{id}", dealHandler.Delete)
	r.Post("/api/deals/{id}/archive", dealHandler.Archive)
	r.Post("/api/deals/{id}/unarchive", dealHandler.Unarchive)
	r.Get("/api/deals", dealHandler.List)

	r.Post("/api/projects/3d", project3DHandler.Create3DProject)
	r.Get("/api/projects/3d/{id}", project3DHandler.GetProjectStatus)
	r.Get("/api/projects/3d/{id}/files", project3DHandler.GetProjectFiles3D)

	fileServer := http.StripPrefix("/media/", http.FileServer(http.Dir("./media")))
	r.Handle("/media/*", fileServer)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"status": "ready",
			"project_name": "sun_ready",
			"version": "v1.0.0"
		}`))
	})

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
