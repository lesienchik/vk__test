package api

import (
	"errors"
	"strings"

	"github.com/valyala/fasthttp"

	m "github.com/lesienchik/vk__test/internal/models"
)

// Мидлвара, которая каждый раз проверяет авторизацию пользователя, с помощью access токена.
func (a *Api) middlVerify(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
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
		userId, errs := a.logic.UserVerify(access)
		if errs != nil {
			a.respErrs(ctx, errs)
			return
		}
		ctx.SetUserValue("userId", userId)
		next(ctx)
	}
}
