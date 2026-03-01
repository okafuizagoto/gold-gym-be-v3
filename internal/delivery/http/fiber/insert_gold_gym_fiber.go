package fiber

import (
	"fmt"
	"io"
	"log"
	"net/http"

	firebaseEntity "gold-gym-be/internal/entity/firebase"
	goldEntity "gold-gym-be/internal/entity/goldgym"
	goldStockEntity "gold-gym-be/internal/entity/stock"
	"gold-gym-be/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go.uber.org/zap"
)

func (h *Handler) InsertGoldGymFiber(c *fiber.Ctx) error {
	var (
		result                   interface{}
		metadata                 interface{}
		err                      error
		resp                     response.Response
		insertgolduser           goldEntity.GetGoldUsers
		insertgoldloginuser      goldEntity.LogUser
		insertgoldsubsuser       goldEntity.InsertSubsAll
		insertUserFirebase       firebaseEntity.User
		insertgoldsubsuserdetail goldEntity.SubscriptionDetail
		insertstock              goldStockEntity.InsertStockData
	)

	ctx := c.UserContext()

	headerCarrier := opentracing.TextMapCarrier{}
	c.Request().Header.VisitAll(func(k, v []byte) {
		headerCarrier[string(k)] = string(v)
	})
	spanCtx, _ := h.tracer.Extract(opentracing.TextMap, headerCarrier)
	span := h.tracer.StartSpan("InsertGoldGymFiber", ext.RPCServerOption(spanCtx))
	defer span.Finish()

	ctx = opentracing.ContextWithSpan(ctx, span)
	h.logger.For(ctx).Info("HTTP request received",
		zap.String("method", c.Method()),
		zap.String("url", c.OriginalURL()))

	types := c.Query("type")
	switch types {
	case "insertuser":
		if err = c.BodyParser(&insertgolduser); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		result, err = h.goldgymSvc.InsertGoldUser(ctx, insertgolduser)

	case "insertuserfirebase":
		if err = c.BodyParser(&insertUserFirebase); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		result, err = h.goldgymSvcStock.CreateUser(ctx, insertUserFirebase)

	case "loginuser":
		if err = c.BodyParser(&insertgoldloginuser); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		host := c.Hostname()
		result, metadata, err = h.goldgymSvc.LoginUser(ctx, insertgoldloginuser.GoldEmail, insertgoldloginuser.GoldPassword, host)
		if err != nil {
			log.Println("err", err)
		}
		fmt.Println("result :", result)
		fmt.Println("metadata :", metadata)
		fmt.Println("err :", err)
		fmt.Println("Result :", insertgoldloginuser)

	case "insertsubsuser":
		if err = c.BodyParser(&insertgoldsubsuser); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		result, err = h.goldgymSvc.InsertSubscriptionUser(ctx, insertgoldsubsuser)

	case "insertsubsuserdetail":
		if err = c.BodyParser(&insertgoldsubsuserdetail); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		result, err, metadata = h.goldgymSvc.InsertSubscriptionDetail(ctx, insertgoldsubsuserdetail)

	case "insertstock":
		if err = c.BodyParser(&insertstock); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		result, err = h.goldgymSvcStock.InsertStockSales(ctx, insertstock)

	case "uploadimages":
		fileHeader, fileErr := c.FormFile("file")
		if fileErr != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Unable to get file"})
		}
		f, fileErr := fileHeader.Open()
		if fileErr != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Unable to open file"})
		}
		defer f.Close()

		imageBytes, readErr := io.ReadAll(f)
		if readErr != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Unable to read file"})
		}

		testings := goldEntity.Testings{TestingImages: imageBytes}
		result, err = h.goldgymSvc.UploadTestingImages(ctx, testings)
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
