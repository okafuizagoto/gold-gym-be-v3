package echo

import (
	"encoding/json"
	goldEntity "gold-gym-be/internal/entity/goldgym"
	"gold-gym-be/pkg/response"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go.uber.org/zap"
)

func (h *Handler) UpdateGoldGymEcho(c echo.Context) error {
	var (
		result             interface{}
		metadata           interface{}
		err                error
		resp               response.Response
		types              string
		updategoldsubsuser goldEntity.UpdateSubs
		updatepassword     goldEntity.UpdatePassword
		updatenama         goldEntity.UpdateNama
		updatekartu        goldEntity.UpdateKartu
		logout             goldEntity.Logout
	)

	ctx := c.Request().Context()

	spanCtx, _ := h.tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request().Header))
	span := h.tracer.StartSpan("Getgoldgym", ext.RPCServerOption(spanCtx))
	defer span.Finish()

	ctx = opentracing.ContextWithSpan(ctx, span)
	h.logger.For(ctx).Info("HTTP request received",
		zap.String("method", c.Request().Method),
		zap.Stringer("url", c.Request().URL))

	types = c.FormValue("type")
	switch types {
	case "updatesubsuser":
		body, _ := ioutil.ReadAll(c.Request().Body)
		json.Unmarshal(body, &updategoldsubsuser)
		result, err = h.goldgymSvc.UpdateSubscriptionDetail(ctx, updategoldsubsuser)
		if err != nil {
			log.Println("err", err)
		}

	case "updatepassword":
		body, _ := ioutil.ReadAll(c.Request().Body)
		json.Unmarshal(body, &updatepassword)
		result, err = h.goldgymSvc.UpdateDataPeserta(ctx, updatepassword)
		if err != nil {
			log.Println("err", err)
		}

	case "updatenama":
		body, _ := ioutil.ReadAll(c.Request().Body)
		json.Unmarshal(body, &updatenama)
		result, err = h.goldgymSvc.UpdateNama(ctx, updatenama)
		if err != nil {
			log.Println("err", err)
		}

	case "updatekartu":
		body, _ := ioutil.ReadAll(c.Request().Body)
		json.Unmarshal(body, &updatekartu)
		result, err = h.goldgymSvc.UpdateKartu(ctx, updatekartu)
		if err != nil {
			log.Println("err", err)
		}

	case "logout":
		body, _ := ioutil.ReadAll(c.Request().Body)
		json.Unmarshal(body, &logout)
		result, err = h.goldgymSvc.Logout(ctx, logout)
		if err != nil {
			log.Println("err", err)
		}
	}

	if err != nil {
		return c.JSON(400, map[string]string{"error": err.Error()})
	}

	resp.Data = result
	resp.Metadata = metadata
	log.Printf("[INFO] %s %s\n", c.Request().Method, c.Request().URL)
	h.logger.For(ctx).Info("HTTP request done",
		zap.String("method", c.Request().Method),
		zap.Stringer("url", c.Request().URL))

	return c.JSON(http.StatusOK, resp)
}
