package beego

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"gold-gym-be/internal/entity/firebase"
	goldEntity "gold-gym-be/internal/entity/goldgym"
	goldStockEntity "gold-gym-be/internal/entity/stock"
	"gold-gym-be/pkg/response"

	beegoCtx "github.com/beego/beego/v2/server/web/context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go.uber.org/zap"
)

func (h *Handler) InsertGoldGymBeego(ctx *beegoCtx.Context) {
	var (
		result                   interface{}
		metadata                 interface{}
		err                      error
		resp                     response.Response
		insertgolduser           goldEntity.GetGoldUsers
		insertgoldloginuser      goldEntity.LogUser
		insertgoldsubsuser       goldEntity.InsertSubsAll
		insertUserFirebase       firebase.User
		insertgoldsubsuserdetail goldEntity.SubscriptionDetail
		insertstock              goldStockEntity.InsertStockData
	)

	reqCtx := ctx.Request.Context()

	spanCtx, _ := h.tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(ctx.Request.Header))
	span := h.tracer.StartSpan("InsertGoldGymBeego", ext.RPCServerOption(spanCtx))
	defer span.Finish()

	reqCtx = opentracing.ContextWithSpan(reqCtx, span)
	h.logger.For(reqCtx).Info("HTTP request received",
		zap.String("method", ctx.Request.Method),
		zap.Stringer("url", ctx.Request.URL))

	types := ctx.Input.Query("type")
	switch types {
	case "insertuser":
		if err = json.Unmarshal(ctx.Input.RequestBody, &insertgolduser); err != nil {
			ctx.Output.SetStatus(http.StatusBadRequest)
			ctx.Output.JSON(map[string]string{"error": err.Error()}, false, false)
			return
		}
		result, err = h.goldgymSvc.InsertGoldUser(reqCtx, insertgolduser)

	case "insertuserfirebase":
		if err = json.Unmarshal(ctx.Input.RequestBody, &insertUserFirebase); err != nil {
			ctx.Output.SetStatus(http.StatusBadRequest)
			ctx.Output.JSON(map[string]string{"error": err.Error()}, false, false)
			return
		}
		result, err = h.goldgymSvcStock.CreateUser(reqCtx, insertUserFirebase)

	case "loginuser":
		if err = json.Unmarshal(ctx.Input.RequestBody, &insertgoldloginuser); err != nil {
			ctx.Output.SetStatus(http.StatusBadRequest)
			ctx.Output.JSON(map[string]string{"error": err.Error()}, false, false)
			return
		}
		host := ctx.Request.Host
		result, metadata, err = h.goldgymSvc.LoginUser(reqCtx, insertgoldloginuser.GoldEmail, insertgoldloginuser.GoldPassword, host)
		if err != nil {
			log.Println("err", err)
		}
		fmt.Println("result :", result)
		fmt.Println("metadata :", metadata)

	case "insertsubsuser":
		if err = json.Unmarshal(ctx.Input.RequestBody, &insertgoldsubsuser); err != nil {
			ctx.Output.SetStatus(http.StatusBadRequest)
			ctx.Output.JSON(map[string]string{"error": err.Error()}, false, false)
			return
		}
		result, err = h.goldgymSvc.InsertSubscriptionUser(reqCtx, insertgoldsubsuser)

	case "insertsubsuserdetail":
		if err = json.Unmarshal(ctx.Input.RequestBody, &insertgoldsubsuserdetail); err != nil {
			ctx.Output.SetStatus(http.StatusBadRequest)
			ctx.Output.JSON(map[string]string{"error": err.Error()}, false, false)
			return
		}
		result, err, metadata = h.goldgymSvc.InsertSubscriptionDetail(reqCtx, insertgoldsubsuserdetail)

	case "insertstock":
		if err = json.Unmarshal(ctx.Input.RequestBody, &insertstock); err != nil {
			ctx.Output.SetStatus(http.StatusBadRequest)
			ctx.Output.JSON(map[string]string{"error": err.Error()}, false, false)
			return
		}
		result, err = h.goldgymSvcStock.InsertStockSales(reqCtx, insertstock)
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
