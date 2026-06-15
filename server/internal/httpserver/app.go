package httpserver

import (
	devprojectsvc "github.com/woragis/management/backend/server/internal/devproject/service"
	mediasvc "github.com/woragis/management/backend/server/internal/media/service"
	"gorm.io/gorm"
)

type App struct {
	DB           *gorm.DB
	AdminAPIKey  string
	MediaDir     string
	MediaBaseURL string
	SecretsKey   []byte

	DevProjects *devprojectsvc.Service
	Media       *mediasvc.Service
}
