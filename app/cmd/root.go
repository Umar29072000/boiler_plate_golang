package cmd

import (
	"fmt"
	"os"
	"time"

	"boiler_plate_be_golang/app/config"
	"boiler_plate_be_golang/internal/database/migrations"
	"boiler_plate_be_golang/internal/repository"
	"boiler_plate_be_golang/internal/service"
	"boiler_plate_be_golang/pkg/logger"
	"boiler_plate_be_golang/pkg/redis"
	"boiler_plate_be_golang/pkg/utils"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLog "gorm.io/gorm/logger"
)

var (
	EnvFilePath string
	rootCmd     = &cobra.Command{
		Use:   "app",
		Short: "Start Fiber Boilerplate Application",
	}

	// Global variables
	DbGorm     *gorm.DB
	rootConfig config.Root

	// Services
	userService service.IUserService
	authService service.IAuthService
)

func init() {
	cobra.OnInitialize(func() {
		initConfigReader()
		initLogger()
		initJWT()
		InitGorm(rootConfig.Postgres)
		initRedis()
		initApp()
	})
}

func initConfigReader() {
	rootConfig = config.Load(EnvFilePath)
	logrus.Info("Config loaded - App: ", rootConfig.App.ServiceName)
	logrus.Info("Postgres loaded: ", rootConfig.Postgres.Host)
}

func initLogger() {
	logger.Init(rootConfig.App.Env)
}

func initJWT() {
	utils.InitJWT(rootConfig.JWT.Secret)
	logrus.Info("JWT initialized")
}

func initLogrus() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)
}

func InitGorm(conf config.Postgres) {
	connectionString := conf.GetDSN()

	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{
		Logger: gormLog.Default.LogMode(gormLog.Info),
	})
	if err != nil || db.Error != nil {
		logger.Fatal("Failed to Initialize to Postgres").Err(err).Send()
	}

	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println(err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DbGorm = db

	// Run migrations if DB_MIGRATE env is set to "yes"
	if os.Getenv("DB_MIGRATE") == "yes" {
		if err := migrations.Migrate(DbGorm); err != nil {
			logrus.Fatal("Failed to run migrations: ", err)
		}
	}

	fmt.Println("PostgreSQL Successfully Connected")
}

func initRedis() {
	// Connect to Redis (optional, non-fatal)
	if err := redis.Connect(redis.RedisConfig{
		Host:     rootConfig.Redis.Host,
		Port:     rootConfig.Redis.Port,
		Password: rootConfig.Redis.Password,
		DB:       rootConfig.Redis.DB,
	}); err != nil {
		logrus.Warn("Redis connection failed, using fallback: ", err)
	}
}

func initApp() {
	// TODO: List of Repositories
	userRepository := repository.NewUserRepository(DbGorm)

	// TODO: List of Services
	userService = service.NewUserService(userRepository)
	authService = service.NewAuthService(userRepository, rootConfig.JWT)

	logrus.Info("Application initialized successfully")
}

// Execute will call the root command execute
func Execute() {
	rootCmd.PersistentFlags().StringVarP(&EnvFilePath, "env", "e", ".env", ".env file to read from")
	if err := rootCmd.Execute(); err != nil {
		logrus.Error("Can't start the CLI: ", err)
		os.Exit(1)
	}
}
