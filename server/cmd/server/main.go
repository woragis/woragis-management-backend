package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/woragis/management/backend/server/internal/httpserver"
	"github.com/woragis/management/backend/server/internal/middleware"
	"github.com/woragis/management/backend/server/internal/migrate"
	"github.com/woragis/management/backend/server/internal/models"
	"github.com/woragis/management/backend/server/internal/platform/listen"
	devprojectrepo "github.com/woragis/management/backend/server/internal/devproject/repository"
	devprojectsvc "github.com/woragis/management/backend/server/internal/devproject/service"
	contactsrepo "github.com/woragis/management/backend/server/internal/contacts/repository"
	contactssvc "github.com/woragis/management/backend/server/internal/contacts/service"
	personalitycache "github.com/woragis/management/backend/server/internal/agent/personality/cache"
	personalityrepo "github.com/woragis/management/backend/server/internal/agent/personality/repository"
	personalitysvc "github.com/woragis/management/backend/server/internal/agent/personality/service"
	messagingrepo "github.com/woragis/management/backend/server/internal/messaging/repository"
	messagingsvc "github.com/woragis/management/backend/server/internal/messaging/service"
	presencerepo "github.com/woragis/management/backend/server/internal/presence/repository"
	presencereminder "github.com/woragis/management/backend/server/internal/presence/reminder"
	presencesvc "github.com/woragis/management/backend/server/internal/presence/service"
	"github.com/woragis/management/backend/server/internal/messaging/executor"
	contentrepo "github.com/woragis/management/backend/server/internal/content/repository"
	contentsvc "github.com/woragis/management/backend/server/internal/content/service"
	"github.com/woragis/management/backend/server/internal/agentworkerclient"
	"github.com/woragis/management/backend/server/internal/creativesclient"
	msgtemplaterender "github.com/woragis/management/backend/server/internal/messaging/templaterender"
	"github.com/woragis/management/backend/server/internal/telegramworkerclient"
	"github.com/woragis/management/backend/server/internal/whatsappworkerclient"
	financerepo "github.com/woragis/management/backend/server/internal/finance/repository"
	financesvc "github.com/woragis/management/backend/server/internal/finance/service"
	mediarepo "github.com/woragis/management/backend/server/internal/media/repository"
	mediasvc "github.com/woragis/management/backend/server/internal/media/service"
	mediastore "github.com/woragis/management/backend/server/internal/media/storage"
	profilerepo "github.com/woragis/management/backend/server/internal/profile/repository"
	profilesvc "github.com/woragis/management/backend/server/internal/profile/service"
	"github.com/woragis/management/backend/server/internal/platform/postgres"
)

func main() {
	addr := listen.Addr()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is required")
	}

	adminKey := strings.TrimSpace(os.Getenv("ADMIN_API_KEY"))
	if adminKey == "" {
		log.Fatal("ADMIN_API_KEY is required")
	}

	db, err := postgres.Open(dsn)
	if err != nil {
		log.Fatalf("database: %v", err)
	}

	if skip := strings.TrimSpace(os.Getenv("SKIP_SQL_MIGRATIONS")); skip != "1" && !strings.EqualFold(skip, "true") {
		dir := strings.TrimSpace(os.Getenv("MIGRATIONS_DIR"))
		if dir == "" {
			dir = migrate.ResolveDir()
		}
		if dir != "" {
			sqlDB, err := db.DB()
			if err != nil {
				log.Fatalf("sql db: %v", err)
			}
			mctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			err = migrate.Up(mctx, sqlDB, dir)
			cancel()
			if err != nil {
				log.Fatalf("sql migrate: %v", err)
			}
			log.Printf("sql migrations applied from %s", dir)
		} else {
			log.Print("warning: SQL migrations skipped (set MIGRATIONS_DIR or run from a tree that contains migrations/)")
		}
	}

	if err := db.AutoMigrate(
		&models.Project{},
		&models.ProjectLink{},
		&models.ProjectDomain{},
		&models.ProjectSecret{},
		&models.ProjectGallery{},
		&models.ProjectEnv{},
		&models.IncomeSource{},
		&models.Expense{},
		&models.Transaction{},
		&models.Invoice{},
		&models.InvoiceItem{},
		&models.BudgetPlan{},
		&models.MediaAsset{},
		&models.Profile{},
		&models.LeetcodeVideo{},
		&models.ContentThumbnail{},
		&models.ContentPromptTemplate{},
		&models.LeetcodeChannelSettings{},
		&models.WhatsappMessageTemplate{},
		&models.Contact{},
		&models.ContactInteraction{},
		&models.AgentPersonality{},
		&models.ChannelDestination{},
		&models.MessageTemplate{},
		&models.ScheduledJob{},
		&models.MessageDelivery{},
		&models.SocialCampaign{},
		&models.PostTemplate{},
		&models.SocialPost{},
		&models.PresenceSettings{},
	); err != nil {
		log.Fatalf("automigrate: %v", err)
	}

	devRepo := devprojectrepo.New(db)
	devSvc := devprojectsvc.New(devRepo, loadSecretsKey())

	financeRepo := financerepo.New(db)
	financeSvc := financesvc.New(financeRepo)

	contactsRepo := contactsrepo.New(db)
	contactsSvc := contactssvc.New(contactsRepo)
	financeSvc.SetContactValidator(contactsSvc)

	var personalityCache personalitycache.Store = personalitycache.Noop{}
	if redisURL := strings.TrimSpace(os.Getenv("REDIS_URL")); redisURL != "" {
		rc, err := personalitycache.NewRedis(redisURL)
		if err != nil {
			log.Fatalf("redis: %v", err)
		}
		if err := rc.Ping(context.Background()); err != nil {
			log.Fatalf("redis ping: %v", err)
		}
		personalityCache = rc
		log.Print("redis connected for agent personality cache")
	} else {
		log.Print("warning: REDIS_URL not set; agent personality uses postgres only")
	}
	personalityRepo := personalityrepo.New(db)
	personalitySvc := personalitysvc.New(personalityRepo, personalityCache)
	if _, err := personalityRepo.EnsureDefault(context.Background()); err != nil {
		log.Fatalf("default agent personality: %v", err)
	}

	mediaBaseURL := strings.TrimSpace(os.Getenv("MEDIA_PUBLIC_BASE_URL"))
	if mediaBaseURL == "" {
		mediaBaseURL = "http://127.0.0.1:8080/v1/public/media"
	}
	mediaStore, driver, err := mediastore.NewFromEnv()
	if err != nil {
		log.Fatalf("media storage: %v", err)
	}
	log.Printf("media storage driver: %s", driver)
	mediaRepo := mediarepo.New(db)
	mediaSvc := mediasvc.New(mediaRepo, mediaStore, mediaBaseURL)

	profileRepo := profilerepo.New(db)
	profileSvc := profilesvc.New(profileRepo)
	if _, err := profileRepo.EnsureDefault(context.Background()); err != nil {
		log.Fatalf("default profile: %v", err)
	}

	creativesClient := creativesclient.New(creativesclient.Config{
		BaseURL: os.Getenv("CREATIVES_API_URL"),
		APIKey:  os.Getenv("CREATIVES_API_KEY"),
	})
	whatsappWorkerClient := whatsappworkerclient.New(whatsappworkerclient.Config{
		BaseURL:      os.Getenv("WHATSAPP_WORKER_URL"),
		WorkerAPIKey: os.Getenv("WORKER_API_KEY"),
	})
	contentRepo := contentrepo.New(db)
	contentSvc := contentsvc.New(
		contentRepo,
		mediaSvc,
		creativesClient,
		whatsappWorkerClient,
		strings.TrimSpace(os.Getenv("MANAGEMENT_WEBHOOK_URL")),
		envOrDefault("CONTENT_THUMBNAIL_DEFAULT_SIZE", "1280x720"),
	)

	if err := contentSvc.EnsureWhatsappDefaults(context.Background()); err != nil {
		log.Fatalf("whatsapp defaults: %v", err)
	}

	messagingRepo := messagingrepo.New(db)
	messagingSvc := messagingsvc.New(messagingRepo)
	contentSvc.SetMessagingTemplates(messagingSvc)

	msgRenderer := msgtemplaterender.NewEngine(contentSvc, devSvc)
	agentWorkerClient := agentworkerclient.New(agentworkerclient.Config{
		BaseURL:     os.Getenv("AGENT_WORKER_URL"),
		AgentAPIKey: strings.TrimSpace(os.Getenv("AGENT_API_KEY")),
	})
	telegramWorkerClient := telegramworkerclient.New(telegramworkerclient.Config{
		BaseURL:      os.Getenv("TELEGRAM_WORKER_URL"),
		WorkerAPIKey: os.Getenv("WORKER_API_KEY"),
	})
	schedulerExec := executor.New(messagingSvc, contentSvc, whatsappWorkerClient, telegramWorkerClient, agentWorkerClient, msgRenderer)

	presenceRepo := presencerepo.New(db)
	presenceSvc := presencesvc.New(presenceRepo)
	if _, err := presenceRepo.EnsureSettings(context.Background()); err != nil {
		log.Fatalf("default presence settings: %v", err)
	}
	presenceReminderExec := presencereminder.New(presenceSvc, messagingSvc, devSvc, whatsappWorkerClient)

	workerAPIKey := strings.TrimSpace(os.Getenv("WORKER_API_KEY"))
	if workerAPIKey == "" {
		log.Print("warning: WORKER_API_KEY not set; internal worker routes disabled")
	}

	agentAPIKey := strings.TrimSpace(os.Getenv("AGENT_API_KEY"))
	if agentAPIKey == "" {
		log.Print("warning: AGENT_API_KEY not set; internal agent routes disabled")
	}

	app := &httpserver.App{
		DB:           db,
		AdminAPIKey:  adminKey,
		AgentAPIKey:  agentAPIKey,
		WorkerAPIKey: workerAPIKey,
		MediaBaseURL: mediaBaseURL,
		SecretsKey:   loadSecretsKey(),
		DevProjects:  devSvc,
		Contacts:     contactsSvc,
		Finance:      financeSvc,
		Media:        mediaSvc,
		MediaRepo:    mediaRepo,
		Profile:      profileSvc,
		Content:      contentSvc,
		Personality:  personalitySvc,
		Messaging:    messagingSvc,
		Presence:         presenceSvc,
		PresenceReminder: presenceReminderExec,
		Scheduler:        schedulerExec,
	}

	cfg := middleware.LoadConfigFromEnv()
	handler := httpserver.NewHandler(app, cfg)

	srv := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("shutdown: %v", err)
	}
}

func loadSecretsKey() []byte {
	raw := strings.TrimSpace(os.Getenv("SECRETS_ENCRYPTION_KEY"))
	if raw == "" {
		return nil
	}
	return []byte(raw)
}

func envOrDefault(key, def string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return def
}
