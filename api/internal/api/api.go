package api

import (
	"api/internal/admin"
	"api/internal/api/interfaces"
	"api/internal/auth"
	"api/internal/command"
	"api/internal/device"
	"api/internal/health"
	"api/internal/upload"
	"api/pkg/ai"
	"api/pkg/database"
	"api/pkg/logging"
	"api/pkg/storage"
)

type Repositories struct {
	Admin   interfaces.AdminRepository
	Command interfaces.CommandRepository
	Device  interfaces.DeviceRepository
	Upload  interfaces.UploadRepository
}

type Services struct {
	Admin   interfaces.AdminService
	Command interfaces.CommandService
	Device  interfaces.DeviceService
	Upload  interfaces.UploadService
}

type Handlers struct {
	Auth    interfaces.AuthHandler
	Command interfaces.CommandHandler
	Device  interfaces.DeviceHandler
	Health  interfaces.HealthHandler
	Upload  interfaces.UploadHandler
}

func InitRepositories(db *database.PostgresDB) *Repositories {
	return &Repositories{
		Admin:   admin.NewPostgreSQLRepository(db),
		Command: command.NewPostgreSQLRepository(db),
		Device:  device.NewPostgreSQLRepository(db),
		Upload:  upload.NewPostgreSQLRepository(db),
	}
}

func InitServices(repos *Repositories, logger *logging.Logger, storage storage.Provider, aiClient ai.Client) *Services {
	s := &Services{}

	s.Admin = admin.NewService(logger.WithField("service", "admin"), repos.Admin)
	s.Command = command.NewService(logger.WithField("service", "command"), repos.Command)
	s.Device = device.NewService(logger.WithField("service", "device"), repos.Device)
	s.Upload = upload.NewService(logger.WithField("service", "upload"), repos.Upload, storage)

	return s
}

func InitHandlers(services *Services, logger *logging.Logger, db database.Database) *Handlers {
	return &Handlers{
		Auth:    auth.NewHandler(logger.WithField("handler", "auth"), services.Admin),
		Command: command.NewHandler(logger.WithField("handler", "command"), services.Command),
		Device:  device.NewHandler(logger.WithField("handler", "device"), services.Device),
		Health:  health.NewHandler(logger.WithField("handler", "health"), db),
		Upload:  upload.NewHandler(logger.WithField("handler", "upload"), services.Upload),
	}
}
