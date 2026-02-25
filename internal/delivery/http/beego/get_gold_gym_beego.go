package beego

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"gold-gym-be/internal/entity"
	"gold-gym-be/pkg/response"

	beegoCtx "github.com/beego/beego/v2/server/web/context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go.uber.org/zap"
)

func (h *Handler) GetGoldGymBeego(ctx *beegoCtx.Context) {
	var (
		result   interface{}
		metadata interface{}
		err      error
		resp     response.Response
	)

	reqCtx := ctx.Request.Context()

	spanCtx, _ := h.tracer.Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(ctx.Request.Header),
	)
	span := h.tracer.StartSpan("GetGoldGymBeego", ext.RPCServerOption(spanCtx))
	defer span.Finish()

	reqCtx = opentracing.ContextWithSpan(reqCtx, span)

	h.logger.For(reqCtx).Info("HTTP request received",
		zap.String("method", ctx.Request.Method),
		zap.Stringer("url", ctx.Request.URL))

	types := ctx.Input.Query("type")

	switch types {
	case "getgoldgym":
		result, err = h.goldgymSvc.GetGoldUser(reqCtx)
		log.Println("deliverygolduser", result)

	case "golduserbyemail":
		result, err = h.goldgymSvc.GetGoldUserByEmail(reqCtx, ctx.Input.Query("email"))

	case "allsubscription":
		result, err = h.goldgymSvc.GetAllSubscription(reqCtx)

	case "getuserandsubsdetail":
		result, err = h.goldgymSvc.GetSubsWithUser(reqCtx)

	case "gettotalpayment":
		result, err = h.goldgymSvc.GetSubscriptionHeaderTotalHarga(reqCtx, ctx.Input.Query("email"))

	case "getonestock":
		result, err = h.goldgymSvcStock.GetOneStockProduct(reqCtx,
			ctx.Input.Query("stockcode"),
			ctx.Input.Query("stockname"),
			ctx.Input.Query("stockid"))
		log.Printf("testDelivery %+v", result)

	case "getallstock":
		result, err = h.goldgymSvcStock.GetAllStockHeader(reqCtx)

	case "getallstockredis":
		result, err = h.goldgymSvcStock.GetAllStockHeaderToRedis(reqCtx)

	case "getfromfirebase":
		result, err = h.goldgymSvcStock.GetFromFirebase(reqCtx, ctx.Input.Query("userid"))

	case "getimages":
		id, _ := strconv.Atoi(ctx.Input.Query("id"))
		result, err = h.goldgymSvc.GetTestingImage(reqCtx, id)
		if err == nil {
			imgData, ok := result.([]byte)
			if !ok {
				ctx.Output.SetStatus(http.StatusInternalServerError)
				ctx.Output.JSON(map[string]string{"error": "invalid image data"}, false, false)
				return
			}
			ctx.Output.Header("Content-Type", "image/png")
			ctx.Output.Body(imgData)
			return
		}
	}

	if err != nil {
		switch {
		case errors.Is(err, entity.ErrNotFound):
			ctx.Output.SetStatus(http.StatusNotFound)
			ctx.Output.JSON(map[string]string{"error": err.Error()}, false, false)
		case errors.Is(err, entity.ErrInvalid):
			ctx.Output.SetStatus(http.StatusBadRequest)
			ctx.Output.JSON(map[string]string{"error": err.Error()}, false, false)
		case errors.Is(err, entity.ErrUnauthorized):
			ctx.Output.SetStatus(http.StatusUnauthorized)
			ctx.Output.JSON(map[string]string{"error": "unauthorized"}, false, false)
		default:
			ctx.Output.SetStatus(http.StatusInternalServerError)
			ctx.Output.JSON(map[string]string{"error": "internal error"}, false, false)
		}
		return
	}

	resp.Data = result
	resp.Metadata = metadata

	log.Printf("[INFO] %s %s\n", ctx.Request.Method, ctx.Request.URL)
	h.logger.For(reqCtx).Info(
		"HTTP request done",
		zap.String("method", ctx.Request.Method),
		zap.String("url", ctx.Request.URL.String()),
	)

	ctx.Output.SetStatus(http.StatusOK)
	ctx.Output.JSON(resp, false, false)
}
