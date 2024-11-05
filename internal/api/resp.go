package api

import (
	"encoding/json"

	"github.com/valyala/fasthttp"

	m "github.com/lesienchik/vk__test/internal/models"
)

func (a *Api) respSucc(ctx *fasthttp.RequestCtx, code int, raw any) {
	var body *m.RespSuccData
	if raw != nil {
		body = &m.RespSuccData{
			Data: raw,
		}
	}

	resp := m.RespSucc{
		Status: "success",
		Code:   code,
		Body:   body,
	}
	data, err := json.Marshal(resp)
	if err != nil {
		a.logger.Error("api.respSucc(1): %w", err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
	ctx.SetStatusCode(code)
	ctx.Response.SetBodyRaw(data)
}

func (a *Api) respErrs(ctx *fasthttp.RequestCtx, errs *m.Err) {
	var detail string
	if errs.Error == nil {
		detail = ""
	} else {
		detail = errs.Error.Error()
	}

	resp := m.RespErr{
		Status: "error",
		Code:   errs.Code,
		Message: m.RespErrMsg{
			Client: errs.ClientMsg,
			Detail: detail,
		},
	}

	data, err := json.Marshal(resp)
	if err != nil {
		a.logger.Error("api.respErrs(1): %w", err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
	ctx.SetStatusCode(resp.Code)
	ctx.Response.SetBodyRaw(data)
}
