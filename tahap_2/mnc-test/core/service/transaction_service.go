package service

import (
	"errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"mnc-test/core/entity"
	"mnc-test/core/repository"
	"mnc-test/data_source/postgres_datasource"
)

type TransactionService struct {
	Repo   *repository.Repository
	Logger *zap.Logger
}

func NewTransactionService(repo *repository.Repository, logger *zap.Logger) *TransactionService {
	return &TransactionService{
		Repo:   repo,
		Logger: logger,
	}
}
func (service *TransactionService) GetAllTransactions(user *entity.User) (transaction *[]entity.Transaction, err error) {
	_, transaction, err = service.Repo.GetAllTransactions(nil, postgres_datasource.QueryOption{
		Filter: map[string]interface{}{
			"userId": map[string]interface{}{
				"field":      "user_id",
				"searchType": "text",
				"match":      "exact",
				"keyword":    user.UserID,
			},
		},
	})
	if err != nil {
		return
	}
	return
}

func (service *TransactionService) CreatePayment(request *entity.RequestCreatePayment, user *entity.User) (transaction *entity.Transaction, err error) {
	newUUID := uuid.New()

	userBalance, err := service.Repo.GetUserBalanceByUid(nil, user.UserID)
	if err != nil {
		return
	}

	if userBalance.Balance < request.Amount {
		err = errors.New("Balance is not enough")
	}

	transaction = &entity.Transaction{
		UserID:          user.UserID,
		PaymentID:       newUUID.String(),
		Remarks:         request.Remarks,
		Amount:          request.Amount,
		TransactionType: entity.Payment,
		BalanceBefore:   userBalance.Balance,
		BalanceAfter:    userBalance.Balance - request.Amount,
	}

	err = service.Repo.CreateTransaction(nil, transaction)

	if err != nil {
		return
	}

	propertyMap := map[string]interface{}{
		"balance": userBalance.Balance - request.Amount,
	}

	err = service.Repo.UpdateUserBalance(nil, userBalance, propertyMap)

	if err != nil {
		return
	}

	return

}

func (service *TransactionService) CreateTopUp(request *entity.RequestCreateTopUp, user *entity.User) (transaction *entity.Transaction, err error) {
	newUUID := uuid.New()

	userBalance, err := service.Repo.GetUserBalanceByUid(nil, user.UserID)
	if err != nil {
		return
	}

	transaction = &entity.Transaction{
		UserID:          user.UserID,
		TopUpID:         newUUID.String(),
		Remarks:         request.Remarks,
		Amount:          request.Amount,
		TransactionType: entity.TopUp,
		BalanceBefore:   userBalance.Balance,
		BalanceAfter:    userBalance.Balance + request.Amount,
	}

	err = service.Repo.CreateTransaction(nil, transaction)

	if err != nil {
		return
	}

	propertyMap := map[string]interface{}{
		"balance": userBalance.Balance + request.Amount,
	}

	err = service.Repo.UpdateUserBalance(nil, userBalance, propertyMap)

	if err != nil {
		return
	}

	return

}

func (service *TransactionService) CreateTransfer(request *entity.RequestCreateTransfer, user *entity.User) (transaction *entity.Transaction, err error) {
	newUUID := uuid.New()

	userBalance, err := service.Repo.GetUserBalanceByUid(nil, user.UserID)
	if err != nil {
		return
	}

	if userBalance.Balance < request.Amount {
		err = errors.New("Balance is not enough")
	}

	transaction = &entity.Transaction{
		UserID:          user.UserID,
		TopUpID:         newUUID.String(),
		Remarks:         request.Remarks,
		Amount:          request.Amount,
		TransactionType: entity.Transfer,
		BalanceBefore:   userBalance.Balance,
		BalanceAfter:    userBalance.Balance - request.Amount,
	}

	err = service.Repo.CreateTransaction(nil, transaction)

	if err != nil {
		return
	}

	propertyMap := map[string]interface{}{
		"balance": userBalance.Balance - request.Amount,
	}

	err = service.Repo.UpdateUserBalance(nil, userBalance, propertyMap)

	if err != nil {
		return
	}

	targetUser, err := service.Repo.GetUserByUid(nil, request.TargetUser)
	if err != nil {
		return
	}

	targetUserBalance, err := service.Repo.GetUserBalanceByUid(nil, targetUser.UserID)
	if err != nil {
		return
	}

	propertyMap = map[string]interface{}{
		"balance": targetUserBalance.Balance + request.Amount,
	}

	err = service.Repo.UpdateUserBalance(nil, targetUserBalance, propertyMap)

	if err != nil {
		return
	}

	return

}
