package service

import (
	"github.com/google/uuid"
	"go.uber.org/zap"
	"mnc-test/core/entity"
	"mnc-test/core/repository"
)

type UserService struct {
	Repo   *repository.Repository
	Logger *zap.Logger
}

func NewUserService(repo *repository.Repository, logger *zap.Logger) *UserService {
	return &UserService{
		Repo:   repo,
		Logger: logger,
	}
}

func (service *UserService) Migration() error {

	err := service.Repo.AutoMigration(nil)

	if err != nil {
		return err
	}

	return nil

}

func (service *UserService) GetUser(uid *string) (user *entity.User, err error) {
	user, err = service.Repo.GetUserByUid(nil, *uid)

	if err != nil {
		return
	}
	return

}
func (service *UserService) GetUserByPhoneNumber(phoneNumber *string) (user *entity.User, err error) {
	user, err = service.Repo.GetUserByPhoneNumber(nil, *phoneNumber)

	if err != nil {
		return
	}
	return

}

func (service *UserService) CreateUser(request *entity.RequestCreateUser) (user *entity.User, err error) {
	newUUID := uuid.New()
	user = &entity.User{
		UserID:      newUUID.String(),
		FirstName:   request.FirstName,
		LastName:    request.LastName,
		PhoneNumber: request.PhoneNumber,
		PIN:         request.PIN,
		Address:     request.Address,
	}

	userBalance := &entity.UserBalance{
		UserID:  user.UserID,
		Balance: 0,
	}

	err = service.Repo.CreateUser(nil, user)

	if err != nil {
		return
	}

	err = service.Repo.CreateUserBalance(nil, userBalance)

	if err != nil {
		return
	}
	return

}

func (service *UserService) UpdateUser(uid *string, request *entity.RequestUpdateUser) (user *entity.User, err error) {
	tx := service.Repo.GetDBPostgres(nil)
	propertyMap := map[string]interface{}{
		"first_name": request.FirstName,
		"last_name":  request.LastName,
		"address":    request.Address,
	}

	user, err = service.Repo.GetUserByUid(tx, *uid)

	if err != nil {
		goto failed
	}

	err = service.Repo.UpdateUser(tx, user, propertyMap)

	if err != nil {
		goto failed
	}

	tx.Commit()
	return

failed:
	tx.Rollback()

	return
}
func (service *UserService) DeleteUser(uid *string) (user *entity.User, err error) {
	tx := service.Repo.GetDBPostgres(nil)
	propertyMap := []map[string]interface{}{}

	user, err = service.Repo.GetUserByUid(nil, *uid)

	if err != nil {
		goto failed
	}

	propertyMap = []map[string]interface{}{
		{
			"key":      "uid",
			"operator": "=",
			"value":    *uid,
		},
	}

	err = service.Repo.DeleteUser(nil, &entity.User{}, propertyMap)

	if err != nil {
		goto failed
	}

	tx.Commit()
	return

failed:
	tx.Rollback()

	return
}
