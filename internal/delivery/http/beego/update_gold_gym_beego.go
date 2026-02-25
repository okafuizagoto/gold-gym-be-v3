package beego

import (
	"encoding/json"
	"log"
	"net/http"

	goldEntity "gold-gym-be/internal/entity/goldgym"
	"gold-gym-be/pkg/response"

	beegoCtx "github.com/beego/beego/v2/server/web/context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go.uber.org/zap"
)

func (h *Handler) UpdateGoldGymBeego(ctx *beegoCtx.Context) {
	var (
		result             interface{}
		metadata           interface{}
		err                error
		resp               response.Response
		updategoldsubsuser goldEntity.UpdateSubs
		updatepassword     goldEntity.UpdatePassword
		updatenama         goldEntity.UpdateNama
		updatekartu        goldEntity.UpdateKartu
		logout             goldEntity.Logout
	)

	reqCtx := ctx.Request.Context()

	spanCtx, _ := h.tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(ctx.Request.Header))
	span := h.tracer.StartSpan("UpdateGoldGymBeego", ext.RPCServerOption(spanCtx))
	defer span.Finish()

	reqCtx = opentracing.ContextWithSpan(reqCtx, span)
	h.logger.For(reqCtx).Info("HTTP request received",
		zap.String("method", ctx.Request.Method),
		zap.Stringer("url", ctx.Request.URL))

	types := ctx.Input.Query("type")
	switch types {
	case "updatesubsuser":
		json.Unmarshal(ctx.Input.RequestBody, &updategoldsubsuser)
		result, err = h.goldgymSvc.UpdateSubscriptionDetail(reqCtx, updategoldsubsuser)
		if err != nil {
			log.Println("err", err)
		}

	case "updatepassword":
		json.Unmarshal(ctx.Input.RequestBody, &updatepassword)
		result, err = h.goldgymSvc.UpdateDataPeserta(reqCtx, updatepassword)
		if err != nil {
			log.Println("err", err)
		}

	case "updatenama":
		json.Unmarshal(ctx.Input.RequestBody, &updatenama)
		result, err = h.goldgymSvc.UpdateNama(reqCtx, updatenama)
		if err != nil {
			log.Println("err", err)
		}

	case "updatekartu":
		json.Unmarshal(ctx.Input.RequestBody, &updatekartu)
		result, err = h.goldgymSvc.UpdateKartu(reqCtx, updatekartu)
		if err != nil {
			log.Println("err", err)
		}

	case "logout":
		json.Unmarshal(ctx.Input.RequestBody, &logout)
		result, err = h.goldgymSvc.Logout(reqCtx, logout)
		if err != nil {
			log.Println("err", err)
		}
	}

	if err != nil {
		ctx.Output.SetStatus(http.StatusBadRequest)
		ctx.Output.JSON(map[string]string{"error": err.Error()}, false, false)
		return
	}

	resp.Data = result
	resp.Metadata = metadata
	log.Printf("[INFO] %s %s\n", ctx.Request.Method, ctx.Request.URL)
	h.logger.For(reqCtx).Info("HTTP request done",
		zap.String("method", ctx.Request.Method),
		zap.Stringer("url", ctx.Request.URL))

	ctx.Output.SetStatus(http.StatusOK)
	ctx.Output.JSON(resp, false, false)
}
