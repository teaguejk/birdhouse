package api

import (
	"api/internal/admin"
	"api/internal/api/interfaces"
	"api/internal/auth"
	"api/internal/device"
	"api/internal/health"
	"api/internal/upload"
	"api/pkg/ai"
	"api/pkg/database"
	"api/pkg/logging"
	"api/pkg/storage"
)

type Repositories struct {
	Admin  interfaces.AdminRepository
	Device interfaces.DeviceRepository
	Upload interfaces.UploadRepository
}

type Services struct {
	Admin  interfaces.AdminService
	Device interfaces.DeviceService
	Upload interfaces.UploadService
}

type Handlers struct {
	Auth   interfaces.AuthHandler
	Device interfaces.DeviceHandler
	Health interfaces.HealthHandler
	Upload interfaces.UploadHandler
}

func InitRepositories(db *database.PostgresDB) *Repositories {
	return &Repositories{
		Admin:  admin.NewPostgreSQLRepository(db),
		Device: device.NewPostgreSQLRepository(db),
		Upload: upload.NewPostgreSQLRepository(db),
	}
}

func InitServices(repos *Repositories, logger *logging.Logger, storage storage.Provider, aiClient ai.Client) *Services {
	s := &Services{}

	s.Admin = admin.NewService(logger.WithField("service", "admin"), repos.Admin)
	s.Device = device.NewService(logger.WithField("service", "device"), repos.Device)
	s.Upload = upload.NewService(logger.WithField("service", "upload"), repos.Upload, storage)

	return s
}

func InitHandlers(services *Services, logger *logging.Logger, db database.Database) *Handlers {
	return &Handlers{
		Auth:   auth.NewHandler(logger.WithField("handler", "auth"), services.Admin),
		Device: device.NewHandler(logger.WithField("handler", "device"), services.Device),
		Health: health.NewHandler(logger.WithField("handler", "health"), db),
		Upload: upload.NewHandler(logger.WithField("handler", "upload"), services.Upload),
	}
}
