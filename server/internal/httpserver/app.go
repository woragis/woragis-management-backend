package httpserver

import "gorm.io/gorm"

type App struct {
	DB            *gorm.DB
	AdminAPIKey   string
	MediaDir      string
	MediaBaseURL  string
	SecretsKey    []byte
}
