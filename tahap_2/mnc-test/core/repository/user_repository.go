package repository

import (
	"fmt"
	"gorm.io/gorm"
	"mnc-test/core/entity"
	"mnc-test/data_source/postgres_datasource"
)

// repo for user
func (r *Repository) CreateUser(tx *gorm.DB, user *entity.User) (err error) {
	err = r.PostgresDatasource.Create(tx, user)
	return
}

func (r *Repository) CreateUserBalance(tx *gorm.DB, userBalance *entity.UserBalance) (err error) {
	err = r.PostgresDatasource.Create(tx, userBalance)
	return
}
func (r *Repository) UpdateUser(tx *gorm.DB, user *entity.User, updateValue map[string]interface{}) (err error) {
	uid := fmt.Sprintf("%s", user.UserID)
	err = r.PostgresDatasource.Update(tx, "user_id", uid, user, updateValue)
	return
}

func (r *Repository) DeleteUser(tx *gorm.DB, user *entity.User, condition []map[string]interface{}) (err error) {
	err = r.PostgresDatasource.Delete(tx, user, condition)
	return
}
func (r *Repository) GetAllUser(tx *gorm.DB, queryOption postgres_datasource.QueryOption) (total int64, attendances *[]entity.User, err error) {
	var resp interface{}

	total, resp, err = r.PostgresDatasource.GetList(tx, entity.User{}, &[]entity.User{}, queryOption)

	attendances, ok := resp.(*[]entity.User)
	if !ok {
		return 0, nil, fmt.Errorf("invalid parsing data to struct")
	}

	return
}

func (r *Repository) GetUserByUid(tx *gorm.DB, userUid string) (user *entity.User, err error) {
	var resp interface{}
	condition := []map[string]interface{}{
		{
			"key":      "user_id",
			"operator": "=",
			"value":    userUid,
		},
	}

	resp, err = r.PostgresDatasource.Get(tx, &entity.User{}, condition)
	user, ok := resp.(*entity.User)
	if !ok {
		return nil, fmt.Errorf("invalid parsing data to struct")
	}

	return
}
func (r *Repository) GetUserByPhoneNumber(tx *gorm.DB, phoneNumber string) (user *entity.User, err error) {
	var resp interface{}
	condition := []map[string]interface{}{
		{
			"key":      "phone_number",
			"operator": "=",
			"value":    phoneNumber,
		},
	}

	resp, err = r.PostgresDatasource.Get(tx, &entity.User{}, condition)
	user, ok := resp.(*entity.User)
	if !ok {
		return nil, fmt.Errorf("invalid parsing data to struct")
	}

	return
}
