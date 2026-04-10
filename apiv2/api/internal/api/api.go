package api

import (
	"api/internal/api/interfaces"
	"api/internal/health"
	"api/internal/upload"
	"api/pkg/ai"
	"api/pkg/database"
	"api/pkg/logging"
	"api/pkg/storage"
)

type Repositories struct {
	Upload interfaces.UploadRepository
}

type Services struct {
	Upload interfaces.UploadService
}

type Handlers struct {
	Health interfaces.HealthHandler
	Upload interfaces.UploadHandler
}

func InitRepositories(db *database.PostgresDB) *Repositories {
	return &Repositories{
		Upload: upload.NewPostgreSQLRepository(db),
	}
}

func InitServices(repos *Repositories, logger *logging.Logger, storage storage.Provider, aiClient ai.Client) *Services {
	s := &Services{}

	s.Upload = upload.NewService(logger.WithField("service", "upload"), repos.Upload, storage)

	return s
}

func InitHandlers(services *Services, logger *logging.Logger, db database.Database) *Handlers {
	return &Handlers{
		Health: health.NewHandler(logger.WithField("handler", "health"), db),
		Upload: upload.NewHandler(logger.WithField("handler", "upload"), services.Upload),
	}
}
