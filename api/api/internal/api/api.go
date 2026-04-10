package api

import (
	"api/internal/api/interfaces"
	"api/internal/device"
	"api/internal/health"
	"api/internal/upload"
	"api/pkg/ai"
	"api/pkg/database"
	"api/pkg/logging"
	"api/pkg/storage"
)

type Repositories struct {
	Device interfaces.DeviceRepository
	Upload interfaces.UploadRepository
}

type Services struct {
	Device interfaces.DeviceService
	Upload interfaces.UploadService
}

type Handlers struct {
	Device interfaces.DeviceHandler
	Health interfaces.HealthHandler
	Upload interfaces.UploadHandler
}

func InitRepositories(db *database.PostgresDB) *Repositories {
	return &Repositories{
		Device: device.NewPostgreSQLRepository(db),
		Upload: upload.NewPostgreSQLRepository(db),
	}
}

func InitServices(repos *Repositories, logger *logging.Logger, storage storage.Provider, aiClient ai.Client) *Services {
	s := &Services{}

	s.Device = device.NewService(logger.WithField("service", "device"), repos.Device)
	s.Upload = upload.NewService(logger.WithField("service", "upload"), repos.Upload, storage)

	return s
}

func InitHandlers(services *Services, logger *logging.Logger, db database.Database) *Handlers {
	return &Handlers{
		Device: device.NewHandler(logger.WithField("handler", "device"), services.Device),
		Health: health.NewHandler(logger.WithField("handler", "health"), db),
		Upload: upload.NewHandler(logger.WithField("handler", "upload"), services.Upload),
	}
}
