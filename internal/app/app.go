package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sensors-generator/config"
	"sensors-generator/internal/generator"
	"sensors-generator/internal/group"
	"sensors-generator/internal/middleware"
	"sensors-generator/internal/mocks"
	"sensors-generator/internal/sensor"
	sensordata "sensors-generator/internal/sensorData"
	"sensors-generator/internal/spiece"
	"sensors-generator/pkg/client/postgresql"
	"sensors-generator/pkg/client/redis"
	"sensors-generator/pkg/logging"
	"sensors-generator/pkg/metric"
	"time"

	_ "sensors-generator/docs"

	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const (
	WriteTimeout = 15
	ReadTimeout  = 15
)

type App struct {
	cfg        *config.Config
	logger     *logging.Logger
	router     *gin.Engine
	httpServer *http.Server
}

func NewApp(ctx context.Context, cfg *config.Config, logger *logging.Logger) (App, error) {
	logger.Info("DB init")
	dbClient, err := postgresql.NewClient(cfg.PgConfig)
	if err != nil {
		logger.Errorf("Failed to connect database, due to error: %v", err)
		return App{}, err
	}

	logger.Info("Redis init")
	redisCache := redis.NewRedisCache(cfg.RedisConfig)

	logger.Info("Gin init")
	router := gin.Default()
	router.Use(middleware.HandleErrors())

	logger.Info("Swagger docs init")
	router.GET("/swagger", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	heartbeatHandler := &metric.Handler{}
	heartbeatHandler.Register(router)

	logger.Info("Create sensor group repo.")
	sensorGroupRepo := group.NewPostgresqlRepository(dbClient, logger, cfg)
	logger.Info("Create sensor group service.")
	sensorGroupService := group.NewService(sensorGroupRepo, redisCache, logger, cfg)
	logger.Info("Create sensor group handler.")
	sensorGroupHandler := group.NewHandler(sensorGroupService, logger)
	logger.Info("Register router for sensor group handler.")
	sensorGroupHandler.Register(router)

	logger.Info("Create sensor repo.")
	sensorRepo := sensor.NewPostgresqlRepository(dbClient, logger, cfg)
	logger.Info("Create sensor service.")
	sensorService := sensor.NewService(sensorRepo, logger, cfg)
	logger.Info("Create sensor handler.")
	sensorHandler := sensor.NewHandler(sensorService, logger)
	logger.Info("Register router for sensor handler.")
	sensorHandler.Register(router)

	logger.Info("Create sensor data repo.")
	sensorDataRepo := sensordata.NewPostgresqlRepository(dbClient, logger, cfg)
	logger.Info("Create sensor data service.")
	sensorDataService := sensordata.NewService(sensorDataRepo, logger, cfg)

	logger.Info("Create spiece repo.")
	spieceRepo := spiece.NewPostgresqlRepository(dbClient, logger, cfg)
	logger.Info("Create spiece service.")
	spieceService := spiece.NewService(spieceRepo, logger, cfg)

	logger.Info("Create Main Entities Generator.")
	meGen := generator.NewMainEntitiesGenerator(generator.MainEntities{
		Groups:  mocks.CreateSensorGroups,
		Sensors: mocks.CreateSensors,
		Spieces: mocks.CreateSpieces,
	}, generator.Services{
		SensorService:      sensorService,
		SensorGroupService: sensorGroupService,
		SpieceService:      spieceService,
	})

	meGen.Generate()

	logger.Info("Create Data Generator.")
	dataGen := generator.NewDataGenerator(generator.Services{
		SensorService:      sensorService,
		SensorGroupService: sensorGroupService,
		SpieceService:      spieceService,
		SensorDataService:  sensorDataService,
	}, generator.NewRandomGenerator())

	// Start data generator only if main entities created.
	for {
		if meGen.IsGenerated() {
			dataGen.Generate()
			break
		}

		time.Sleep(1 * time.Second)
	}

	// Start app only if data generator started.
	return App{
		cfg:    cfg,
		logger: logger,
		router: router,
	}, nil
}

func (a *App) Run() {
	a.startHTTP()
}

func (a *App) startHTTP() error {
	a.logger.Info("Start http")

	var listener net.Listener

	if a.cfg.Listen.Type == "sock" {
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))

		if err != nil {
			a.logger.Fatal(err)
		}

		a.logger.Info("Create socket")
		socketPath := path.Join(appDir, a.cfg.Listen.SocketFile)

		a.logger.Debugf("Socket path: %s", socketPath)

		a.logger.Info("Listen unix socket")

		listener, err = net.Listen("unix", socketPath)
		a.logger.Infof("Server is listening unix socket: %s", socketPath)

	} else {
		var err error

		a.logger.Info("Listen tcp")
		listener, err = net.Listen("tcp", fmt.Sprintf("%s:%s", a.cfg.Listen.BindIP, a.cfg.Listen.Port))
		a.logger.Infof("Server is listening %s:%s", a.cfg.Listen.BindIP, a.cfg.Listen.Port)
		if err != nil {
			a.logger.Fatal(err)
		}
	}

	c := cors.New(cors.Options{
		AllowedMethods:     a.cfg.CorsConfig.AllowedMethods,
		AllowedOrigins:     a.cfg.CorsConfig.AllowedOrigins,
		AllowCredentials:   a.cfg.CorsConfig.AllowCredentials,
		AllowedHeaders:     a.cfg.CorsConfig.AllowedHeaders,
		OptionsPassthrough: a.cfg.CorsConfig.OptionsPassthrough,
		ExposedHeaders:     a.cfg.CorsConfig.ExposedHeaders,
		Debug:              a.cfg.IsDebug,
	})

	handler := c.Handler(a.router)

	a.httpServer = &http.Server{
		Handler:      handler,
		WriteTimeout: WriteTimeout * time.Second,
		ReadTimeout:  ReadTimeout * time.Second,
	}

	if err := a.httpServer.Serve(listener); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			a.logger.Warning("server shutdown")
		default:
			a.logger.Fatal(err)
		}
	}
	err := a.httpServer.Shutdown(context.Background())
	if err != nil {
		a.logger.Fatal(err)
	}
	return err
}
