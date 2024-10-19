package transport

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"mnc-test/core/repository"
	"mnc-test/core/service"
	"mnc-test/data_source/postgres_datasource"
	"mnc-test/transport/transaction_transport"
	"mnc-test/transport/user_transport"
)

type Transport struct {
	user_transport.HttpHandlerUser
	transaction_transport.HttpHandlerTransaction
}

func Setup(e *echo.Echo, client *gorm.DB, logger *zap.Logger) (transport *Transport) {
	transport = &Transport{}
	//datasource
	var postgresDatasource = postgres_datasource.NewPostgresDatasource(client)

	//repository
	var repo = repository.NewRepository(postgresDatasource)

	//service
	//var websocketService = websocket_service.New(tag, isLocal, config, repository)
	var (
		userService        = service.NewUserService(repo, logger)
		transactionService = service.NewTransactionService(repo, logger)
	)

	//setup transport

	transport.NewUserHttpHandler(e, userService)
	transport.NewTransactionHttpHandler(e, transactionService)

	return
}
