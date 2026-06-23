package httpserver

import (
	devprojectsvc "github.com/woragis/management/backend/server/internal/devproject/service"
	personalitysvc "github.com/woragis/management/backend/server/internal/agent/personality/service"
	contentsvc "github.com/woragis/management/backend/server/internal/content/service"
	contactssvc "github.com/woragis/management/backend/server/internal/contacts/service"
	financesvc "github.com/woragis/management/backend/server/internal/finance/service"
	mediarepo "github.com/woragis/management/backend/server/internal/media/repository"
	mediasvc "github.com/woragis/management/backend/server/internal/media/service"
	messagingsvc "github.com/woragis/management/backend/server/internal/messaging/service"
	presencesvc "github.com/woragis/management/backend/server/internal/presence/service"
	profilesvc "github.com/woragis/management/backend/server/internal/profile/service"
	"github.com/woragis/management/backend/server/internal/messaging/executor"
	"gorm.io/gorm"
)

type App struct {
	DB           *gorm.DB
	AdminAPIKey  string
	AgentAPIKey  string
	WorkerAPIKey string
	MediaBaseURL string
	SecretsKey   []byte

	DevProjects *devprojectsvc.Service
	Contacts    *contactssvc.Service
	Finance     *financesvc.Service
	Content     *contentsvc.Service
	Media       *mediasvc.Service
	MediaRepo   *mediarepo.Repository
	Profile     *profilesvc.Service
	Personality *personalitysvc.Service
	Messaging   *messagingsvc.Service
	Presence    *presencesvc.Service
	Scheduler   *executor.Executor
}
