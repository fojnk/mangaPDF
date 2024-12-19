package transport

import (
	"net/http"
	"regexp"

	"github.com/fojnk/Task-Test-devBack/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	ipv6_regex = `^(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))$`
	ipv4_regex = `^(((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|$)){4})`
)

type TokenPair struct {
	AccesToken   string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type InputRegister struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type InputLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// @Summary Login
// @Tags Auth
// @Description Generate tokens
// @ID login
// @Param input body InputLogin true "account info"
// @Param Ip header string true "Ip"
// @Produce  json
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} transort_error
// @Failure 500 {object} transort_error
// @Failure default {object} transort_error
// @Router /auth/login [post]
func (h *Handler) login(c *gin.Context) {
	header := c.GetHeader("Ip")
	var input InputLogin

	match, err := regexp.MatchString(ipv4_regex+`|`+ipv6_regex, header)
	if !match || err != nil {
		NewTransportErrorResponse(c, http.StatusBadRequest, "bad IP format ")
		return
	}

	if err := c.BindJSON(&input); err != nil {
		NewTransportErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.services.GetUserByUsername(input.Username, input.Password)
	if err != nil {
		NewTransportErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	accessToken, refreshToken, err := h.services.GenerateTokens(user.Id)
	if err != nil {
		NewTransportErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

// @Summary Refresh
// @Param Ip header string true "Ip"
// @Tags Auth
// @Description Refresh
// @ID refresh
// @Accept  json
// @Produce  json
// @Param input body TokenPair true "tokens"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} transort_error
// @Failure 500 {object} transort_error
// @Failure default {object} transort_error
// @Router /auth/refresh [post]
func (h *Handler) refresh(c *gin.Context) {
	header := c.GetHeader("Ip")

	match, err := regexp.MatchString(ipv4_regex+`|`+ipv6_regex, header)
	if !match || err != nil {
		NewTransportErrorResponse(c, http.StatusBadRequest, "bad IP format ")
		return
	}

	var input TokenPair

	if err := c.BindJSON(&input); err != nil {
		NewTransportErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	newAccess, newRefresh, err := h.services.Refresh(input.AccesToken, input.RefreshToken)
	if err != nil {
		NewTransportErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"accessToken":  newAccess,
		"refreshToken": newRefresh,
	})
}

// @Summary Register
// @Tags Auth
// @Description create account
// @ID create-account
// @Param Ip header string true "Ip"
// @Accept  json
// @Produce  json
// @Param input body InputRegister true "account info"
// @Success 200 {integer} integer 1
// @Failure 400,404,500,default {object} transort_error
// @Router /auth/register [post]
func (h *Handler) register(c *gin.Context) {
	header := c.GetHeader("Ip")
	var input InputRegister

	match, err := regexp.MatchString(ipv4_regex+`|`+ipv6_regex, header)
	if !match || err != nil {
		NewTransportErrorResponse(c, http.StatusBadRequest, "bad IP format ")
		return
	}

	if err := c.BindJSON(&input); err != nil {
		NewTransportErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	user := models.User{
		Email:    input.Email,
		Password: input.Password,
		Username: input.Username,
	}

	logrus.Printf("create user with %s, %s", input.Email, input.Username)

	id, err := h.services.Authorization.CreateUser(user)
	if err != nil {
		NewTransportErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	logrus.Printf("generate tokens for user: %d", id)

	accessToken, refreshToken, err := h.services.Authorization.GenerateTokens(id)
	if err != nil {
		NewTransportErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id":           id,
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

// @Summary Get Account Ingo
// @Security ApiKeyAuth
// @Tags Account
// @Description Get accound by id
// @ID get-account
// @Produce  json
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} transort_error
// @Failure 500 {object} transort_error
// @Failure default {object} transort_error
// @Router /api/v1/account [get]
func (h *Handler) getAccountInfo(c *gin.Context) {
	userId, ok := c.Get(UserId)
	if !ok {
		NewTransportErrorResponse(c, http.StatusBadRequest, "You are not authorized!!!")
		return
	}

	user, err := h.services.Authorization.GetUserById(userId.(int))

	if err != nil {
		NewTransportErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"user": user,
	})
}
