package httpserver

import (
	devprojectsvc "github.com/woragis/management/backend/server/internal/devproject/service"
	contentsvc "github.com/woragis/management/backend/server/internal/content/service"
	financesvc "github.com/woragis/management/backend/server/internal/finance/service"
	mediarepo "github.com/woragis/management/backend/server/internal/media/repository"
	mediasvc "github.com/woragis/management/backend/server/internal/media/service"
	profilesvc "github.com/woragis/management/backend/server/internal/profile/service"
	"gorm.io/gorm"
)

type App struct {
	DB           *gorm.DB
	AdminAPIKey  string
	MediaBaseURL string
	SecretsKey   []byte

	DevProjects *devprojectsvc.Service
	Finance     *financesvc.Service
	Content     *contentsvc.Service
	Media       *mediasvc.Service
	MediaRepo   *mediarepo.Repository
	Profile     *profilesvc.Service
}
