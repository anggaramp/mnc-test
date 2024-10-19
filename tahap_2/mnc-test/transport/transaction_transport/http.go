package transaction_transport

import (
	"github.com/labstack/echo/v4"
	"mnc-test/core/entity"
	"mnc-test/core/service"
	"mnc-test/shared"
	"mnc-test/shared/utils"
	"net/http"
)

type HttpHandlerTransaction struct {
	TransactionService *service.TransactionService
}

func (httpHandler *HttpHandlerTransaction) NewTransactionHttpHandler(e *echo.Echo, transactionService *service.TransactionService) {
	httpHandler.TransactionService = transactionService

	//e.GET("/api/user", httpHandler.GetAllTransaction)
	e.POST("/topup", httpHandler.TopUp, utils.AuthMiddleware)
	e.POST("/pay", httpHandler.Payment, utils.AuthMiddleware)
	e.POST("/transfer", httpHandler.Transfer, utils.AuthMiddleware)
	e.GET("/transactions", httpHandler.Transactions, utils.AuthMiddleware)

}

func (httpHandler *HttpHandlerTransaction) TopUp(c echo.Context) error {
	var err error
	var user entity.User
	if val, ok := c.Get("user").(entity.User); ok {
		user = val
	} else {
		return c.JSON(http.StatusUnauthorized, shared.MakeFailedReponse("Unauthenticated"))
	}

	request := new(entity.RequestCreateTopUp)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusOK, shared.MakeFailedReponse(err.Error()))
	}

	if err = utils.ValidateStruct(request); err != nil {
		return c.JSON(http.StatusBadRequest, shared.MakeFailedReponse(err.Error()))
	}

	transaction, err := httpHandler.TransactionService.CreateTopUp(request, &user)

	if err != nil {
		return c.JSON(http.StatusOK, shared.MakeFailedReponse(err.Error()))
	}

	return c.JSON(http.StatusOK, shared.MakeSuccessResponse(transaction))
}

func (httpHandler *HttpHandlerTransaction) Payment(c echo.Context) error {
	var err error
	var user entity.User
	if val, ok := c.Get("user").(entity.User); ok {
		user = val
	} else {
		return c.JSON(http.StatusUnauthorized, shared.MakeFailedReponse("Unauthenticated"))
	}

	request := new(entity.RequestCreatePayment)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusOK, shared.MakeFailedReponse(err.Error()))
	}

	if err = utils.ValidateStruct(request); err != nil {
		return c.JSON(http.StatusBadRequest, shared.MakeFailedReponse(err.Error()))
	}

	transaction, err := httpHandler.TransactionService.CreatePayment(request, &user)

	if err != nil {
		return c.JSON(http.StatusOK, shared.MakeFailedReponse(err.Error()))
	}

	return c.JSON(http.StatusOK, shared.MakeSuccessResponse(transaction))
}

func (httpHandler *HttpHandlerTransaction) Transfer(c echo.Context) error {
	var err error
	var user entity.User
	if val, ok := c.Get("user").(entity.User); ok {
		user = val
	} else {
		return c.JSON(http.StatusUnauthorized, shared.MakeFailedReponse("Unauthenticated"))
	}

	request := new(entity.RequestCreateTransfer)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusOK, shared.MakeFailedReponse(err.Error()))
	}

	if err = utils.ValidateStruct(request); err != nil {
		return c.JSON(http.StatusBadRequest, shared.MakeFailedReponse(err.Error()))
	}

	transaction, err := httpHandler.TransactionService.CreateTransfer(request, &user)

	if err != nil {
		return c.JSON(http.StatusOK, shared.MakeFailedReponse(err.Error()))
	}

	return c.JSON(http.StatusOK, shared.MakeSuccessResponse(transaction))
}

func (httpHandler *HttpHandlerTransaction) Transactions(c echo.Context) error {
	var err error
	var user entity.User
	if val, ok := c.Get("user").(entity.User); ok {
		user = val
	} else {
		return c.JSON(http.StatusUnauthorized, shared.MakeFailedReponse("Unauthenticated"))
	}

	transaction, err := httpHandler.TransactionService.GetAllTransactions(&user)

	if err != nil {
		return c.JSON(http.StatusOK, shared.MakeFailedReponse(err.Error()))
	}

	return c.JSON(http.StatusOK, shared.MakeSuccessResponse(transaction))
}
