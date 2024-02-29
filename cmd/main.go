package main

import (
	"fmt"
	logGo "log"
	"strconv"
	"time"
	"user-service/configs"
	addressHandler "user-service/internal/modules/address/handlers"
	addressRepoQuery "user-service/internal/modules/address/repositories/queries"
	addressUsecase "user-service/internal/modules/address/usecases"
	userHandler "user-service/internal/modules/user/handlers"
	userRepoCommand "user-service/internal/modules/user/repositories/commands"
	userRepoQuery "user-service/internal/modules/user/repositories/queries"
	userUsecase "user-service/internal/modules/user/usecases"
	"user-service/internal/pkg/apm"
	"user-service/internal/pkg/databases/mongodb"
	graceful "user-service/internal/pkg/gs"
	"user-service/internal/pkg/helpers"
	kafkaConfluent "user-service/internal/pkg/kafka/confluent"
	"user-service/internal/pkg/log"
	"user-service/internal/pkg/redis"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.elastic.co/apm/module/apmfiber"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// @BasePath	/
func main() {
	// Init Config
	configs.InitConfig()

	// Init Elastic APM Config or DD Apm
	if configs.GetConfig().Datadog.DatadogEnabled == "true" {
		tracer.Start(
			tracer.WithAgentAddr(fmt.Sprintf("%s:%s", configs.GetConfig().Datadog.DatadogHost, configs.GetConfig().Datadog.DatadogPort)),
			tracer.WithEnv(configs.GetConfig().Datadog.DatadogEnv),
			tracer.WithService(configs.GetConfig().Datadog.DatadogService),
		)
		defer tracer.Stop()
	} else {
		// temp handling until we move all to Datadog
		apm.InitConnection(configs.GetConfig().APMElastic.APMUrl, configs.GetConfig().APMElastic.APMSecretToken)
	}
	// Init MongoDB Connection
	mongo := mongodb.MongoImpl{}
	mongo.SetCollections(&mongo)
	mongo.InitConnection(configs.GetConfig().MongoDB.MongoMasterDBUrl, configs.GetConfig().MongoDB.MongoSlaveDBUrl)
	// Init Logger
	logZap := log.SetupLogger()
	log.Init(logZap)

	// Init BlacklistedEmail
	helpers.InitReadBlackListEmail()

	// Init Kafka Config
	kafkaConfluent.InitKafkaConfig(configs.GetConfig().Kafka.KafkaUrl, configs.GetConfig().Kafka.KafkaUsername, configs.GetConfig().Kafka.KafkaPassword)

	// Init instance fiber
	app := fiber.New(fiber.Config{
		BodyLimit: 30 * 1024 * 1024,
	})
	app.Use(apmfiber.Middleware(apmfiber.WithTracer(apm.GetTracer())))
	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(pprof.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
		// Format:       `${time} {"router_activity" : [${status},"${latency}","${method}","${path}"], "query_param":${queryParams}, "body_param":${body}}` + "\n",
		TimeInterval: time.Millisecond,
		TimeFormat:   "02-01-2006 15:04:05",
		TimeZone:     "Indonesia/Jakarta",
	}))
	shutdownDelay, _ := strconv.Atoi(configs.GetConfig().ShutDownDelay)
	// graceful shutdown setup
	gs := &graceful.GracefulShutdown{
		Timeout:        5 * time.Second,
		GracefulPeriod: time.Duration(shutdownDelay) * time.Second,
	}
	app.Get("/healthz", gs.LivenessCheck)
	app.Get("/readyz", gs.ReadinessCheck)
	gs.Enable(app)

	setHttp(app, gs)

	//=== listen port ===//
	if err := app.Listen(fmt.Sprintf(":%s", configs.GetConfig().ServicePort)); err != nil {
		logGo.Fatal(err)
	}
}

func setHttp(app *fiber.App, gs *graceful.GracefulShutdown) {
	// Init Redis
	redisClient := redis.InitConnection(configs.GetConfig().Redis.RedisDB, configs.GetConfig().Redis.RedisHost, configs.GetConfig().Redis.RedisPort,
		configs.GetConfig().Redis.RedisPassword, configs.GetConfig().Redis.RedisAppConfig)
	// Init Jwt
	helperImpl := &helpers.JwtImpl{}
	helperImpl.InitConfig(configs.GetConfig().Jwt.JwtPrivateKey, configs.GetConfig().Jwt.JwtPublicKey,
		configs.GetConfig().Jwt.JwtRefreshPrivateKey, configs.GetConfig().Jwt.JwtRefreshPublicKey)

	logger := log.GetLogger()
	mongoMasterClient := mongodb.NewMongoDBLogger(mongodb.GetMasterConn(), mongodb.GetMasterDBName(), logger)
	mongoSlaveClient := mongodb.NewMongoDBLogger(mongodb.GetSlaveConn(), mongodb.GetMasterDBName(), logger)
	kafkaProducer, err := kafkaConfluent.NewProducer(kafkaConfluent.GetConfig().GetKafkaConfig(configs.GetConfig().ServiceName, true), logger)
	if err != nil {
		panic(err)
	}
	gs.Register(
		mongoMasterClient,
		mongoSlaveClient,
		graceful.FnWithError(redisClient.Close),
		kafkaProducer,
	)

	userQueryMongodbRepo := userRepoQuery.NewQueryMongodbRepository(mongoSlaveClient, logger)
	userCommandMongodbRepo := userRepoCommand.NewCommandMongodbRepository(mongoMasterClient, logger)

	addressQueryMongodbRepo := addressRepoQuery.NewQueryMongodbRepository(mongoSlaveClient, logger)
	addressUsecaseQuery := addressUsecase.NewQueryUsecase(addressQueryMongodbRepo, logger)

	userUsecaseQuery := userUsecase.NewQueryUsecase(userQueryMongodbRepo, userCommandMongodbRepo, logger)
	userUsecaseCommand := userUsecase.NewCommandUsecase(userQueryMongodbRepo, userCommandMongodbRepo, logger, redisClient, kafkaProducer, helperImpl, addressQueryMongodbRepo)
	// set module
	userHandler.InitUserHttpHandler(app, userUsecaseCommand, userUsecaseQuery, logger, redisClient)
	addressHandler.InitAddressHttpHandler(app, addressUsecaseQuery, logger, redisClient)

}
