package config

import (
	"github.com/jakemckenzie/chirpy-server/internal/database"
	"github.com/jakemckenzie/chirpy-server/internal/services"
)

type APIConfig struct {
	MetricsService *services.MetricsService
	TextService    *services.TextService
	DBQueries      *database.Queries
	Platform       string
}
