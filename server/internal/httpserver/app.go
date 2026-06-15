package httpserver

import (
	devprojectsvc "github.com/woragis/management/backend/server/internal/devproject/service"
	"gorm.io/gorm"
)

type App struct {
	DB           *gorm.DB
	AdminAPIKey  string
	MediaDir     string
	MediaBaseURL string
	SecretsKey   []byte

	DevProjects *devprojectsvc.Service
}
