package main

import (
	"database/sql"
	"kando-backend/behaviours"
	"kando-backend/commands"
	"kando-backend/config"
	"kando-backend/db"
	"kando-backend/ioc"
	"kando-backend/log"
	"kando-backend/mediator"
	"kando-backend/server"
	"kando-backend/services"
	"os"
)

func main() {
	config.ReadConfig()
	log.SetupLogging()

	log.Logger.Info("Application starting...")
	log.Logger.Infof("Using environment '%s'", config.C.Environment)

	dbConnection := db.ConnectToDatabase()
	db.MigrateDatabase(dbConnection)

	server.ServeApi(configureServices(dbConnection))

	log.Logger.Info("Application shutting down...")
	os.Exit(0)
}

func configureServices(dbConnection *sql.DB) *ioc.DependencyProvider {
	dpb := ioc.NewDependencyProviderBuilder()

	ioc.AddSingleton(dpb, func(dp *ioc.DependencyProvider) *sql.DB {
		return dbConnection
	})

	ioc.AddScoped(dpb, func(dp *ioc.DependencyProvider) *services.RequestContextService {
		return services.NewRequestContextService(dp)
	})
	ioc.AddCloseHandler[*services.RequestContextService](dpb, func(rcs *services.RequestContextService) error {
		return rcs.Close()
	})

	ioc.AddScoped(dpb, func(dp *ioc.DependencyProvider) *mediator.Mediator {
		m := mediator.NewMediator()

		mediator.RegisterBehaviour[any](m, behaviours.LoggingBehaviour)

		mediator.RegisterHandler(m, commands.CreateUserCommandHandler)

		mediator.RegisterHandler(m, commands.CreateSessionCommandHandler)

		mediator.RegisterHandler(m, commands.CreateManufacturerCommandHandler)

		return m
	})

	return dpb.Build()
}
