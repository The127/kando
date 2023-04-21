package queries

import (
	"context"
	"database/sql"
	"kando-backend/ioc"
	"kando-backend/mediator"
	"kando-backend/services"
)

func testContext(dbConn *sql.DB) context.Context {
	dpb := ioc.NewDependencyProviderBuilder()

	ioc.AddSingleton(dpb, func(dp *ioc.DependencyProvider) *sql.DB {
		return dbConn
	})

	ioc.AddScoped(dpb, func(dp *ioc.DependencyProvider) *services.RequestContextService {
		return services.NewRequestContextService(dp)
	})
	ioc.AddCloseHandler[*services.RequestContextService](dpb, func(rcs *services.RequestContextService) error {
		return rcs.Close()
	})

	ioc.AddScoped(dpb, func(dp *ioc.DependencyProvider) *mediator.Mediator {
		m := mediator.NewMediator()

		return m
	})

	dp := dpb.Build()

	return context.WithValue(context.TODO(), "scope", dp)
}
