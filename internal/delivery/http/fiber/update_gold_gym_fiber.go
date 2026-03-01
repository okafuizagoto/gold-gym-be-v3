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

func (h *Handler) UpdateGoldGymFiber(c *fiber.Ctx) error {
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

	ctx := c.UserContext()

	headerCarrier := opentracing.TextMapCarrier{}
	c.Request().Header.VisitAll(func(k, v []byte) {
		headerCarrier[string(k)] = string(v)
	})
	spanCtx, _ := h.tracer.Extract(opentracing.TextMap, headerCarrier)
	span := h.tracer.StartSpan("UpdateGoldGymFiber", ext.RPCServerOption(spanCtx))
	defer span.Finish()

	ctx = opentracing.ContextWithSpan(ctx, span)
	h.logger.For(ctx).Info("HTTP request received",
		zap.String("method", c.Method()),
		zap.String("url", c.OriginalURL()))

	types := c.Query("type")
	switch types {
	case "updatesubsuser":
		if err = c.BodyParser(&updategoldsubsuser); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		result, err = h.goldgymSvc.UpdateSubscriptionDetail(ctx, updategoldsubsuser)
		if err != nil {
			log.Println("err", err)
		}

	case "updatepassword":
		if err = c.BodyParser(&updatepassword); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		result, err = h.goldgymSvc.UpdateDataPeserta(ctx, updatepassword)
		if err != nil {
			log.Println("err", err)
		}

	case "updatenama":
		if err = c.BodyParser(&updatenama); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		result, err = h.goldgymSvc.UpdateNama(ctx, updatenama)
		if err != nil {
			log.Println("err", err)
		}

	case "updatekartu":
		if err = c.BodyParser(&updatekartu); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		result, err = h.goldgymSvc.UpdateKartu(ctx, updatekartu)
		if err != nil {
			log.Println("err", err)
		}

	case "logout":
		if err = c.BodyParser(&logout); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		result, err = h.goldgymSvc.Logout(ctx, logout)
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
