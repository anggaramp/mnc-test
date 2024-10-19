package user_transport

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"mnc-test/core/entity"
	"mnc-test/core/service"
	"mnc-test/shared"
	"mnc-test/shared/utils"
	"net/http"
)

type HttpHandlerUser struct {
	UserService *service.UserService
}

func (httpHandler *HttpHandlerUser) NewUserHttpHandler(e *echo.Echo, userService *service.UserService) {
	httpHandler.UserService = userService

	//e.GET("/api/user", httpHandler.GetAllUser)
	e.POST("/migration", httpHandler.Migration)
	e.POST("/register", httpHandler.CreateUser)
	e.POST("/login", httpHandler.LoginUser)
	e.PUT("/profile", httpHandler.UpdateUser, utils.AuthMiddleware)

}

func (httpHandler *HttpHandlerUser) GetUser(c echo.Context) error {
	var err error
	var user *entity.User

	uid := c.Param("uid")
	if uid == "" {
		return c.JSON(http.StatusOK, shared.MakeFailedReponse("empty uid"))
	}

	user, err = httpHandler.UserService.GetUser(&uid)

	if err != nil {
		return c.JSON(http.StatusOK, shared.MakeFailedReponse(err.Error()))
	}

	return c.JSON(http.StatusOK, shared.MakeSuccessResponse(user))
}

func (httpHandler *HttpHandlerUser) Migration(c echo.Context) error {
	var err error

	err = httpHandler.UserService.Migration()

	if err != nil {
		return c.JSON(http.StatusOK, shared.MakeFailedReponse(err.Error()))
	}

	return c.JSON(http.StatusOK, nil)
}

func (httpHandler *HttpHandlerUser) LoginUser(c echo.Context) error {
	var err error
	var user *entity.User

	request := new(entity.RequestLoginUser)
	if err = c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, shared.MakeFailedReponse(err.Error()))
	}

	if err = utils.ValidateStruct(request); err != nil {
		return c.JSON(http.StatusBadRequest, shared.MakeFailedReponse(err.Error()))
	}

	user, err = httpHandler.UserService.GetUserByPhoneNumber(&request.PhoneNumber)

	if err != nil {
		return c.JSON(http.StatusOK, shared.MakeFailedReponse("Phone Number and PIN doesn’t match."))
	}
	fmt.Println(user.PIN, request.PIN)
	if user.PIN != request.PIN {
		fmt.Println("not match")

		return c.JSON(http.StatusUnauthorized, shared.MakeFailedReponse("Phone Number and PIN doesn’t match."))
	}

	token, err := utils.GenerateToken(*user)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusOK, shared.MakeFailedReponse("Phone Number and PIN doesn’t match."))
	}

	return c.JSON(http.StatusOK, shared.MakeSuccessResponse(token))
}

func (httpHandler *HttpHandlerUser) CreateUser(c echo.Context) error {
	var err error
	var user *entity.User

	request := new(entity.RequestCreateUser)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusOK, shared.MakeFailedReponse(err.Error()))
	}

	if err = utils.ValidateStruct(request); err != nil {
		return c.JSON(http.StatusBadRequest, shared.MakeFailedReponse(err.Error()))
	}

	user, err = httpHandler.UserService.GetUserByPhoneNumber(&request.PhoneNumber)

	if err == nil && user != nil {
		return c.JSON(http.StatusBadRequest, shared.MakeFailedReponse("User Already registered"))
	}

	user, err = httpHandler.UserService.CreateUser(request)

	if err != nil {
		return c.JSON(http.StatusOK, shared.MakeFailedReponse(err.Error()))
	}

	return c.JSON(http.StatusOK, shared.MakeSuccessResponse(user))
}

func (httpHandler *HttpHandlerUser) UpdateUser(c echo.Context) error {
	var err error
	var user entity.User
	if val, ok := c.Get("user").(entity.User); ok {
		user = val
	} else {
		return c.JSON(http.StatusUnauthorized, shared.MakeFailedReponse("Unauthenticated"))
	}

	request := new(entity.RequestUpdateUser)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusOK, shared.MakeFailedReponse(err.Error()))
	}

	if err = utils.ValidateStruct(request); err != nil {
		return c.JSON(http.StatusBadRequest, shared.MakeFailedReponse(err.Error()))
	}

	userUpdate, err := httpHandler.UserService.UpdateUser(&user.UserID, request)

	if err != nil {
		return c.JSON(http.StatusOK, shared.MakeFailedReponse(err.Error()))
	}

	return c.JSON(http.StatusOK, shared.MakeSuccessResponse(userUpdate))
}
