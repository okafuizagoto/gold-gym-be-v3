package fiber

import (
	"log"
	"net/http"

	goldEntity "gold-gym-be/internal/entity/goldgym"
	"gold-gym-be/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go.uber.org/zap"
)

func (h *Handler) DeleteGoldGymFiber(c *fiber.Ctx) error {
	var (
		result             interface{}
		metadata           interface{}
		err                error
		resp               response.Response
		deletegoldsubsuser goldEntity.DeleteSubs
	)

	ctx := c.UserContext()

	headerCarrier := opentracing.TextMapCarrier{}
	c.Request().Header.VisitAll(func(k, v []byte) {
		headerCarrier[string(k)] = string(v)
	})
	spanCtx, _ := h.tracer.Extract(opentracing.TextMap, headerCarrier)
	span := h.tracer.StartSpan("DeleteGoldGymFiber", ext.RPCServerOption(spanCtx))
	defer span.Finish()

	ctx = opentracing.ContextWithSpan(ctx, span)
	h.logger.For(ctx).Info("HTTP request received",
		zap.String("method", c.Method()),
		zap.String("url", c.OriginalURL()))

	types := c.Query("type")
	switch types {
	case "deletesubsuser":
		if err = c.BodyParser(&deletegoldsubsuser); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		result, err = h.goldgymSvc.DeleteSubscriptionHeader(ctx, deletegoldsubsuser)
		if err != nil {
			log.Println("err", err)
		}
	}

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	resp.Data = result
	resp.Metadata = metadata
	log.Printf("[INFO] %s %s\n", c.Method(), c.OriginalURL())
	h.logger.For(ctx).Info("HTTP request done",
		zap.String("method", c.Method()),
		zap.String("url", c.OriginalURL()))

	return c.JSON(resp)
}
