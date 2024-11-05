package api

import (
	"log"
	"runtime/debug"
	"strings"
	"time"

	fasthttprouter "github.com/fasthttp/router"
	"github.com/sirupsen/logrus"
	fasthttpswagger "github.com/swaggo/fasthttp-swagger"
	"github.com/valyala/fasthttp"

	_ "github.com/lesienchik/vk__test/docs"
	"github.com/lesienchik/vk__test/internal/config"
	"github.com/lesienchik/vk__test/internal/logic"
)

type Api struct {
	addr   string
	logger *logrus.Logger
	router *fasthttprouter.Router
	server *fasthttp.Server
	logic  *logic.Logic
}

func New(cfg *config.Config, logger *logrus.Logger, logic *logic.Logic) *Api {
	api := new(Api)
	router := fasthttprouter.New()
	httpServer := &fasthttp.Server{
		Name:               "vktest",
		MaxRequestBodySize: 1_000_000,
		ReadTimeout:        time.Duration(cfg.Api.ReadTimeout) * time.Second,
		WriteTimeout:       time.Duration(cfg.Api.WriteTimeout) * time.Second,
		Handler:            httpHandle(api),
	}

	api.addr = cfg.Api.Addr
	api.logger = logger
	api.router = router
	api.server = httpServer
	api.logic = logic

	return api
}

func httpHandle(api *Api) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[panic]: Panic during api operation: %s\n%s\n", r, string(debug.Stack()))
			}
		}()
		api.handle(ctx)
	}
}

func (a *Api) handle(ctx *fasthttp.RequestCtx) {
	path, method := string(ctx.Path()), string(ctx.Method())
	switch {
	case ctx.IsOptions():
		ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
		ctx.Response.Header.Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		ctx.Response.Header.Set("Access-Control-Allow-Headers", "Content-Type")
		ctx.SetStatusCode(fasthttp.StatusNoContent) // 204 No Content

	// User
	case path == "/api/v1/user/register" && method == fasthttp.MethodPost:
		a.userRegister(ctx)
	case path == "/api/v1/user/confirm/registration" && method == fasthttp.MethodGet:
		a.userConfirm(ctx)
	case path == "/api/v1/user/auth" && method == fasthttp.MethodPost:
		a.userAuth(ctx)
	case path == "/api/v1/user/refresh" && method == fasthttp.MethodGet:
		a.userRefresh(ctx)
	// Swagger docs
	case strings.HasPrefix(path, "/swagger"):
		fasthttpswagger.WrapHandler(fasthttpswagger.InstanceName("swagger"))(ctx)

	// Test api
	case path == "/status" && method == fasthttp.MethodGet:
		a.status(ctx)
	default:
		ctx.SetStatusCode(fasthttp.StatusNotFound)
	}
}

func (a *Api) Start() error {
	return a.server.ListenAndServe(a.addr)
}

func (a *Api) Shutdown() error {
	return a.server.Shutdown()
}

// @Summary status
// @Tags Liveness
// @Description Показывает статус запуска приложения (сервера).
// @ID status
// @Accept json
// @Produce json
// @Success 200 {object} models.RespSucc
// @Router /status [get]
func (a *Api) status(ctx *fasthttp.RequestCtx) {
	a.respSucc(ctx, fasthttp.StatusOK, "Hi, welcome! I'm the server for the site Vktest!")
}
