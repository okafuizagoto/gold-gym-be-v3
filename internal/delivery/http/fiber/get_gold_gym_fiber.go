package fiber

import (
	"bytes"
	"errors"
	"image"
	"image/png"
	"log"
	"net/http"
	"strconv"

	"gold-gym-be/internal/entity"
	"gold-gym-be/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go.uber.org/zap"
)

func (h *Handler) GetGoldGymFiber(c *fiber.Ctx) error {
	var (
		result   interface{}
		metadata interface{}
		err      error
		resp     response.Response
	)

	ctx := c.UserContext()

	// Extract tracing headers from Fiber request (fasthttp → opentracing TextMap)
	headerCarrier := opentracing.TextMapCarrier{}
	c.Request().Header.VisitAll(func(k, v []byte) {
		headerCarrier[string(k)] = string(v)
	})
	spanCtx, _ := h.tracer.Extract(opentracing.TextMap, headerCarrier)
	span := h.tracer.StartSpan("GetGoldGymFiber", ext.RPCServerOption(spanCtx))
	defer span.Finish()

	ctx = opentracing.ContextWithSpan(ctx, span)

	h.logger.For(ctx).Info("HTTP request received",
		zap.String("method", c.Method()),
		zap.String("url", c.OriginalURL()))

	types := c.Query("type")

	switch types {
	case "getgoldgym":
		result, err = h.goldgymSvc.GetGoldUser(ctx)
		log.Println("deliverygolduser", result)

	case "golduserbyemail":
		result, err = h.goldgymSvc.GetGoldUserByEmail(ctx, c.Query("email"))

	case "allsubscription":
		result, err = h.goldgymSvc.GetAllSubscription(ctx)

	case "getuserandsubsdetail":
		result, err = h.goldgymSvc.GetSubsWithUser(ctx)

	case "gettotalpayment":
		result, err = h.goldgymSvc.GetSubscriptionHeaderTotalHarga(ctx, c.Query("email"))

	case "getonestock":
		result, err = h.goldgymSvcStock.GetOneStockProduct(ctx,
			c.Query("stockcode"),
			c.Query("stockname"),
			c.Query("stockid"))
		log.Printf("testDelivery %+v", result)

	case "getallstock":
		result, err = h.goldgymSvcStock.GetAllStockHeader(ctx)

	case "getallstockredis":
		result, err = h.goldgymSvcStock.GetAllStockHeaderToRedis(ctx)

	case "getfromfirebase":
		result, err = h.goldgymSvcStock.GetFromFirebase(ctx, c.Query("userid"))

	case "getimages":
		id, _ := strconv.Atoi(c.Query("id"))
		result, err = h.goldgymSvc.GetTestingImage(ctx, id)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "internal error"})
		}

		imgData, ok := result.([]byte)
		if !ok {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "invalid image data"})
		}

		img, _, decErr := image.Decode(bytes.NewReader(imgData))
		if decErr != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Unable to decode image"})
		}

		var buf bytes.Buffer
		if encErr := png.Encode(&buf, img); encErr != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "unable to encode image"})
		}

		c.Set("Content-Type", "image/png")
		return c.Send(buf.Bytes())
	}

	if err != nil {
		switch {
		case errors.Is(err, entity.ErrNotFound):
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		case errors.Is(err, entity.ErrInvalid):
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		case errors.Is(err, entity.ErrUnauthorized):
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		default:
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "internal error"})
		}
	}

	resp.Data = result
	resp.Metadata = metadata

	log.Printf("[INFO] %s %s\n", c.Method(), c.OriginalURL())
	h.logger.For(ctx).Info("HTTP request done",
		zap.String("method", c.Method()),
		zap.String("url", c.OriginalURL()))

	return c.JSON(resp)
}
