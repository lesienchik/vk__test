package api

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/valyala/fasthttp"

	m "github.com/lesienchik/vk__test/internal/models"
)

// @Summary userRegister
// @Tags User
// @Description Валидирует и отправляет на почту пользователя ссылку для завершения регистрации
// @ID userRegister
// @Accept json
// @Produce json
// @Param input body models.UserRegReq true "Данные пользователя"
// @Success 200 {object} models.RespSucc
// @Failure default {object} models.RespErr
// @Router /api/v1/user/register [post]
func (a *Api) userRegister(ctx *fasthttp.RequestCtx) {
	var userReq m.UserRegReq
	if err := json.Unmarshal(ctx.PostBody(), &userReq); err != nil {
		a.respErrs(ctx, &m.Err{
			Code:  fasthttp.StatusBadRequest,
			Error: err,
		})
		return
	}

	errs := a.logic.UserRegister(&userReq)
	if errs != nil {
		a.respErrs(ctx, errs)
		return
	}
	a.respSucc(ctx, fasthttp.StatusAccepted, "request accepted")
}

// @Summary userConfirm
// @Tags User
// @Description Завершает регистрацию пользователя и выдает access/refresh пару токенов
// @ID userConfirm
// @Accept json
// @Produce json
// @Param code query string true "Код подтверждения (с почты)"
// @Success 200 {object} models.RespSucc
// @Failure default {object} models.RespErr
// @Router /api/v1/user/confirm/registration [get]
func (a *Api) userConfirm(ctx *fasthttp.RequestCtx) {
	code := ctx.QueryArgs().Peek("code")
	if len(code) == 0 {
		a.respErrs(ctx, &m.Err{
			Code:  fasthttp.StatusBadRequest,
			Error: errors.New("empty verification code"),
		})
		return
	}

	userId, errs := a.logic.UserConfirm(string(code))
	if errs != nil {
		a.respErrs(ctx, errs)
		return
	}
	a.userSetJwtTokens(ctx, userId)
}

// Устанавливает refresh-токен в cookie и возвращает в ответе access-токен.
func (a *Api) userSetJwtTokens(ctx *fasthttp.RequestCtx, userId int) {
	tokens, errs := a.logic.UserSetJwtTokens(userId)
	if errs != nil {
		a.respErrs(ctx, errs)
		return
	}

	if len(tokens) != 2 {
		a.respErrs(ctx, &m.Err{
			Code:  fasthttp.StatusInternalServerError,
			Error: errors.New("error creating a token pair"),
		})
		return
	}

	access, refresh := tokens[0], tokens[1]

	refreshCookie := fasthttp.AcquireCookie()
	defer fasthttp.ReleaseCookie(refreshCookie)

	refreshCookie.SetKey("refresh_token")
	refreshCookie.SetValue(refresh)
	refreshCookie.SetPath("/")
	refreshCookie.SetHTTPOnly(true)
	refreshCookie.SetExpire(time.Now().Add(30 * 24 * time.Hour))
	ctx.Response.Header.SetCookie(refreshCookie)

	data := m.UserAccessResp{Token: access}
	a.respSucc(ctx, fasthttp.StatusOK, data)
}

// @Summary userAuth
// @Tags User
// @Description Аутентифицирует пользователя с помощью логина и пароля. Выдает ему access/refresh пару токенов
// @ID userAuth
// @Accept json
// @Produce json
// @Param input body models.UserAuthReq true "Данные пользователя"
// @Success 200 {object} models.RespSucc
// @Failure default {object} models.RespErr
// @Router /api/v1/user/auth [post]
func (a *Api) userAuth(ctx *fasthttp.RequestCtx) {
	var userReq m.UserAuthReq
	if err := json.Unmarshal(ctx.PostBody(), &userReq); err != nil {
		a.respErrs(ctx, &m.Err{
			Code:  fasthttp.StatusBadRequest,
			Error: err,
		})
		return
	}

	userId, errs := a.logic.UserAuth(&userReq)
	if errs != nil {
		a.respErrs(ctx, errs)
		return
	}
	a.userSetJwtTokens(ctx, userId)
}

// @Summary userRefresh
// @Security ApiKeyAuth
// @Tags User
// @Description Выдает новый access токен пользователю с помощью refresh токена
// @ID userRefresh
// @Accept json
// @Produce json
// @Success 200 {object} models.RespSucc
// @Failure default {object} models.RespErr
// @Router /api/v1/user/refresh [get]
func (a *Api) userRefresh(ctx *fasthttp.RequestCtx) {
	authHeader := ctx.Request.Header.Peek("Authorization")
	if len(authHeader) == 0 {
		a.respErrs(ctx, &m.Err{
			Code:  fasthttp.StatusUnauthorized,
			Error: errors.New("empty authorization header"),
		})
		return
	}

	authParts := strings.Split(string(authHeader), " ")
	if len(authParts) != 2 {
		a.respErrs(ctx, &m.Err{
			Code:  fasthttp.StatusUnauthorized,
			Error: errors.New("authorization header does not consist of 2 parts"),
		})
		return
	}

	if authParts[0] != "Bearer" {
		a.respErrs(ctx, &m.Err{
			Code:  fasthttp.StatusUnauthorized,
			Error: errors.New("authorization header does not start with Bearer"),
		})
		return
	}

	access := authParts[1]
	refresh := ctx.Request.Header.Cookie("refresh_token")
	if len(refresh) == 0 {
		a.respErrs(ctx, &m.Err{
			Code:  fasthttp.StatusUnauthorized,
			Error: errors.New("refresh token not found"),
		})
		return
	}

	token, errs := a.logic.UserRefresh(access, string(refresh))
	if errs != nil {
		a.respErrs(ctx, errs)
		return
	}
	a.respSucc(ctx, fasthttp.StatusOK, m.UserAccessResp{Token: token})
}
