package echo

import (
	goldEntity "gold-gym-be/internal/entity/goldgym"
	"gold-gym-be/pkg/response"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go.uber.org/zap"
)

func (h *Handler) DeleteGoldGymEcho(c echo.Context) error {
	var (
		result             interface{}
		metadata           interface{}
		err                error
		resp               response.Response
		deletegoldsubsuser goldEntity.DeleteSubs
	)

	ctx := c.Request().Context()

	spanCtx, _ := h.tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request().Header))
	span := h.tracer.StartSpan("Getgoldgym", ext.RPCServerOption(spanCtx))
	defer span.Finish()

	ctx = opentracing.ContextWithSpan(ctx, span)
	h.logger.For(ctx).Info("HTTP request received",
		zap.String("method", c.Request().Method),
		zap.Stringer("url", c.Request().URL))

	types := c.QueryParam("type")
	switch types {
	case "deletesubsuser":
		if err := c.Bind(&deletegoldsubsuser); err != nil {
			return c.JSON(400, map[string]string{"error": err.Error()})
		}
		result, err = h.goldgymSvc.DeleteSubscriptionHeader(ctx, deletegoldsubsuser)
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
