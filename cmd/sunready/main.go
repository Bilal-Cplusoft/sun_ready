package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/Bilal-Cplusoft/sun_ready/internal/database"
	"github.com/Bilal-Cplusoft/sun_ready/internal/handler"
	custommw "github.com/Bilal-Cplusoft/sun_ready/internal/middleware"
	"github.com/Bilal-Cplusoft/sun_ready/internal/repo"
	"github.com/Bilal-Cplusoft/sun_ready/internal/service"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Get configuration from environment
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize database
	db, err := database.New(databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize repositories
	userRepo := repo.NewUserRepo(db)
	companyRepo := repo.NewCompanyRepo(db)
	projectRepo := repo.NewProjectRepo(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, jwtSecret)
	userService := service.NewUserService(userRepo)
	companyService := service.NewCompanyService(companyRepo)
	projectService := service.NewProjectService(projectRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	projectHandler := handler.NewProjectHandler(projectService)

	// Initialize router
	r := chi.NewRouter()

	// Middleware
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

	// Public routes
	r.Post("/api/auth/register", authHandler.Register)
	r.Post("/api/auth/login", authHandler.Login)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(custommw.AuthMiddleware(authService))

		// User routes
		r.Get("/api/users/{id}", userHandler.GetByID)
		r.Put("/api/users/{id}", userHandler.Update)
		r.Delete("/api/users/{id}", userHandler.Delete)
		r.Get("/api/users", userHandler.List)

		// Project routes
		r.Post("/api/projects", projectHandler.Create)
		r.Get("/api/projects/{id}", projectHandler.GetByID)
		r.Put("/api/projects/{id}", projectHandler.Update)
		r.Delete("/api/projects/{id}", projectHandler.Delete)
		r.Get("/api/projects", projectHandler.ListByCompany)
		r.Get("/api/projects/user", projectHandler.ListByUser)
	})

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
