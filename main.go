package main

import (
	"go-bank/cmd/api"
	"go-bank/cmd/grpc"
	"go-bank/config"
	"go-bank/db"
	"go-bank/internal/mail"
	"go-bank/internal/worker"
	"os"

	sqlc "go-bank/db/sqlc"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg, err := config.Load(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cant read config file")
	}

	if cfg.ENVIRONMENT == "dev" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	db, err := db.Start(cfg.DB_DSN)

	if err != nil {
		log.Fatal().Err(err).Msg("cant start db")
	}

	mailer := mail.NewEmailSender(cfg.EMAIL_NAME, cfg.EMAIL_ADDRESS, cfg.EMAIL_PASSWORD)

	redisOpt := asynq.RedisClientOpt{
		Addr: cfg.REDIS_ADDR,
	}

	distributor := worker.NewRedisTaskDistributor(redisOpt)

	go startTaskProcessor(redisOpt, db, mailer)

	go startGrpcServer(cfg, db, distributor)

	startGatewayServer(cfg, db, distributor)

	//startHTTPServer(cfg, db)
}

func startHTTPServer(cfg config.Config, db sqlc.Store) {
	httpServer, err := api.NewHTTPServer(cfg, db)
	if err != nil {
		log.Fatal().Err(err).Msg("cant create http server")
	}

	err = httpServer.Start(cfg.SRV_ADDR)
	if err != nil {
		log.Fatal().Err(err).Msg("cant start http server")
	}
}

func startGrpcServer(cfg config.Config, db sqlc.Store, distributor worker.TaskDistributor) {
	grpcServer, err := grpc.NewGrpcServer(cfg, db, distributor)
	if err != nil {
		log.Fatal().Err(err).Msg("cant create grpc server")
	}

	err = grpcServer.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("cant start grpc server")
	}
}

func startGatewayServer(cfg config.Config, db sqlc.Store, distributor worker.TaskDistributor) {
	grpcServer, err := grpc.NewGrpcServer(cfg, db, distributor)
	if err != nil {
		log.Fatal().Err(err).Msg("cant create grpc gateway server")
	}

	err = grpcServer.StartGateway()
	if err != nil {
		log.Fatal().Err(err).Msg("cant start grpc gateway server")
	}
}

func startTaskProcessor(redisOpt asynq.RedisClientOpt, store sqlc.Store, mailer mail.EmailSender) {
	processor := worker.NewRedisTaskProcessor(redisOpt, store, mailer)

	log.Info().Msg("starting distributed task processor")
	err := processor.Start()

	if err != nil {
		log.Fatal().Err(err).Msg("failed to start distributed task processor")
	}
}
