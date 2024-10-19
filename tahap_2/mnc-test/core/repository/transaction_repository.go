package repository

import (
	"fmt"
	"gorm.io/gorm"
	"mnc-test/core/entity"
	"mnc-test/data_source/postgres_datasource"
)

// repo for transaction
func (r *Repository) CreateTransaction(tx *gorm.DB, transaction *entity.Transaction) (err error) {
	err = r.PostgresDatasource.Create(tx, transaction)
	return
}

func (r *Repository) GetUserBalanceByUid(tx *gorm.DB, userUid string) (userBalance *entity.UserBalance, err error) {
	var resp interface{}
	condition := []map[string]interface{}{
		{
			"key":      "user_id",
			"operator": "=",
			"value":    userUid,
		},
	}

	resp, err = r.PostgresDatasource.Get(tx, &entity.UserBalance{}, condition)
	userBalance, ok := resp.(*entity.UserBalance)
	if !ok {
		return nil, fmt.Errorf("invalid parsing data to struct")
	}

	return
}

func (r *Repository) UpdateUserBalance(tx *gorm.DB, userBalance *entity.UserBalance, updateValue map[string]interface{}) (err error) {
	uid := fmt.Sprintf("%s", userBalance.UserID)
	err = r.PostgresDatasource.Update(tx, "user_id", uid, userBalance, updateValue)
	return
}

func (r *Repository) GetAllTransactions(tx *gorm.DB, queryOption postgres_datasource.QueryOption) (total int64, transaction *[]entity.Transaction, err error) {
	var resp interface{}

	total, resp, err = r.PostgresDatasource.GetList(tx, entity.Transaction{}, &[]entity.Transaction{}, queryOption)

	transaction, ok := resp.(*[]entity.Transaction)
	if !ok {
		return 0, nil, fmt.Errorf("invalid parsing data to struct")
	}

	return
}
