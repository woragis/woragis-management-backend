package management

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"woragis-management-service/internal/domains/projects"
	"woragis-management-service/internal/domains/ideas"
	"woragis-management-service/internal/domains/chats"
	"woragis-management-service/internal/domains/clients"
	"woragis-management-service/internal/domains/finances"
	"woragis-management-service/internal/domains/experiences"
	"woragis-management-service/internal/domains/userpreferences"
	"woragis-management-service/internal/domains/userprofiles"
	"woragis-management-service/internal/domains/apikeys"
	"woragis-management-service/internal/domains/languages"
	"woragis-management-service/internal/domains/scheduler"
	"woragis-management-service/internal/domains/testimonials"
	"woragis-management-service/pkg/authservice"
	"woragis-management-service/pkg/aiservice"
	"woragis-management-service/pkg/middleware"
)

// SetupRoutes sets up all management service routes
func SetupRoutes(api fiber.Router, db *gorm.DB, authServiceURL, aiServiceURL string, logger *slog.Logger) {
	// Initialize Auth Service client
	authClient := authservice.NewClient(authServiceURL)
	
	// Initialize AI Service client
	var aiClient *aiservice.Client
	if aiServiceURL != "" {
		aiClient = aiservice.NewClient(aiServiceURL)
	}

	// Apply auth validation middleware to all routes
	api.Use(middleware.AuthValidationMiddleware(middleware.DefaultAuthValidationConfig(authClient)))

	// Initialize repositories
	projectRepo := projects.NewGormRepository(db)
	ideaRepo := ideas.NewGormRepository(db)
	chatRepo := chats.NewGormRepository(db)
	clientRepo := clients.NewGormRepository(db)
	financeRepo := finances.NewGormRepository(db)
	experienceRepo := experiences.NewGormRepository(db)
	userPreferenceRepo := userpreferences.NewGormRepository(db)
	userProfileRepo := userprofiles.NewGormRepository(db)
	apiKeyRepo := apikeys.NewGormRepository(db)
	languageRepo := languages.NewGormRepository(db)
	schedulerRepo := scheduler.NewGormRepository(db)
	testimonialRepo := testimonials.NewGormRepository(db)

	// Initialize services
	projectService := projects.NewService(projectRepo, logger)
	ideaService := ideas.NewService(ideaRepo, logger)
	chatService := chats.NewService(chatRepo, aiClient, logger, "", "", "startup", nil) // provider, model, agent, stream
	clientService := clients.NewService(clientRepo, logger)
	financeService := finances.NewService(financeRepo, logger)
	experienceService := experiences.NewService(experienceRepo, logger)
	userPreferenceService := userpreferences.NewService(userPreferenceRepo, logger)
	userProfileService := userprofiles.NewService(userProfileRepo, logger)
	apiKeyService := apikeys.NewService(apiKeyRepo, logger)
	languageService := languages.NewService(languageRepo, logger)
	schedulerService := scheduler.NewService(schedulerRepo, nil, logger) // reports service nil for now
	testimonialService := testimonials.NewService(testimonialRepo, logger)

	// Initialize handlers (simplified - without translation enricher for now)
	projectHandler := projects.NewHandler(projectService, nil, nil, logger) // enricher, translationService
	ideaHandler := ideas.NewHandler(ideaService, logger)
	chatHandler := chats.NewHandler(chatService, logger, nil) // stream hub nil for now
	clientHandler := clients.NewHandler(clientService, logger)
	financeHandler := finances.NewHandler(financeService, logger)
	experienceHandler := experiences.NewHandler(experienceService, logger)
	userPreferenceHandler := userpreferences.NewHandler(userPreferenceService, logger)
	userProfileHandler := userprofiles.NewHandler(userProfileService, logger)
	apiKeyHandler := apikeys.NewHandler(apiKeyService, logger)
	languageHandler := languages.NewHandler(languageService, logger)
	schedulerHandler := scheduler.NewHandler(schedulerService, logger)
	testimonialHandler := testimonials.NewHandler(testimonialService, nil, nil, logger) // enricher, translationService

	// Setup routes
	projects.SetupRoutes(api.Group("/projects"), projectHandler)
	ideas.SetupRoutes(api.Group("/ideas"), ideaHandler)
	chats.SetupRoutes(api.Group("/chats"), chatHandler)
	clients.SetupRoutes(api.Group("/clients"), clientHandler)
	finances.SetupRoutes(api.Group("/finance"), financeHandler) // Note: finances domain uses "/finance" not "/finances"
	experiences.SetupRoutes(api.Group("/experiences"), experienceHandler)
	userpreferences.SetupRoutes(api.Group("/user-preferences"), userPreferenceHandler)
	userprofiles.SetupRoutes(api.Group("/user-profiles"), userProfileHandler)
	apikeys.SetupRoutes(api.Group("/api-keys"), apiKeyHandler)
	languages.SetupRoutes(api.Group("/languages"), languageHandler)
	scheduler.SetupRoutes(api.Group("/scheduler"), schedulerHandler)
	testimonials.SetupRoutes(api.Group("/testimonials"), testimonialHandler)
}
