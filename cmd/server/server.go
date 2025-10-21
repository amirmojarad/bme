package server

import (
	"bme/conf"
	"bme/database"
	"bme/internal/controller"
	"bme/internal/controller/middleware"
	"bme/internal/repository"
	"bme/internal/service"
	"bme/pkg/errorext"
	"bme/pkg/jwt"
	"bme/pkg/logger"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pressly/goose/v3"
	gormPg "gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func Run(cfg *conf.AppConfig) error {
	bmsJwt := jwt.New(cfg)

	sqlDb, err := database.ConnectToPostgres(cfg)
	if err != nil {
		logger.GetLogger().Fatal(err)
	}

	gormDB, err := getGormDB(cfg, sqlDb)
	if err != nil {
		logger.GetLogger().Fatal(err)

		return err
	}

	gormWrapper, err := setupDb(cfg, sqlDb, gormDB)
	if err != nil {
		logger.GetLogger().Fatal(err)
	}

	userRepo := repository.NewUserRepository(gormWrapper)
	deviceRepo := repository.NewDevice(gormWrapper)

	userSvc := service.NewUserService(userRepo)
	authSvc := service.NewAuth(userSvc, userRepo)
	deviceSvc := service.NewDevice(deviceRepo)

	authController := controller.NewAuth(
		authSvc,
		logger.GetLogger().WithField("name", "auth-controller"),
		bmsJwt,
	)

	deviceCtrl := controller.NewDevice(deviceSvc, logger.GetLogger().WithField("name", "device-controller"))

	// Middlewares
	authMiddleware := middleware.NewAuth(bmsJwt, cfg.Jwt.AccessSecret)

	ginEngine := gin.New()

	v1Group := ginEngine.Group("/v1")
	authGroup := v1Group.Group("/auth")
	deviceGroup := v1Group.Group("/device", authMiddleware.Authorize)

	controller.SetupAuthRoutes(authGroup, authController)
	controller.SetupDeviceRoutes(deviceGroup, deviceCtrl)

	return ginEngine.Run(fmt.Sprintf(":%d", cfg.App.Port))
}

func getGormDB(cfg *conf.AppConfig, postgresDB *sql.DB) (*gorm.DB, error) {
	var (
		gormCfg *gorm.Config
		err     error
	)

	gormCfg = &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Silent),
	}

	if cfg.IsEnvDebug() {
		gormCfg = &gorm.Config{
			Logger: gormLogger.Default.LogMode(gormLogger.Info),
		}
	}

	psql, err := gorm.Open(gormPg.New(gormPg.Config{Conn: postgresDB}), gormCfg)
	if err != nil {
		return nil, errorext.New(err, errorext.ErrGeneralOccurrence)
	}

	return psql, nil
}

func setupDb(cfg *conf.AppConfig, sqlDb *sql.DB, gormDB *gorm.DB) (*database.GormWrapper, error) {
	if err := goose.SetDialect("postgres"); err != nil {
		return nil, errorext.New(err, errorext.ErrGeneralOccurrence)
	}
	if err := goose.Up(sqlDb, cfg.Postgres.MigrationPath); err != nil && err.Error() != "no migration files found" {
		return nil, errorext.New(err, errorext.ErrGeneralOccurrence)
	}

	gormWrapper := database.NewGormWrapper(gormDB)

	return gormWrapper, nil
}
