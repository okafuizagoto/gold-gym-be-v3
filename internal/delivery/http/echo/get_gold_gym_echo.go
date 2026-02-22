package echo

import (
	"bytes"
	"errors"
	"gold-gym-be/internal/entity"
	"gold-gym-be/pkg/response"
	"image"
	"image/png"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go.uber.org/zap"
)

func (h *Handler) GetGoldGymEcho(c echo.Context) error {
	var (
		result   interface{}
		metadata interface{}
		err      error
		resp     response.Response
	)

	ctx := c.Request().Context()

	spanCtx, _ := h.tracer.Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(c.Request().Header),
	)
	span := h.tracer.StartSpan("GetGoldGym", ext.RPCServerOption(spanCtx))
	defer span.Finish()

	ctx = opentracing.ContextWithSpan(ctx, span)

	h.logger.For(ctx).Info("HTTP request received",
		zap.String("method", c.Request().Method),
		zap.Stringer("url", c.Request().URL))

	types := c.QueryParam("type")

	switch types {
	case "getgoldgym":
		result, err = h.goldgymSvc.GetGoldUser(ctx)
		log.Println("deliverygolduser", result)

	case "golduserbyemail":
		result, err = h.goldgymSvc.GetGoldUserByEmail(ctx, c.QueryParam("email"))

	case "allsubscription":
		result, err = h.goldgymSvc.GetAllSubscription(ctx)

	case "getuserandsubsdetail":
		result, err = h.goldgymSvc.GetSubsWithUser(ctx)

	case "gettotalpayment":
		result, err = h.goldgymSvc.GetSubscriptionHeaderTotalHarga(ctx, c.QueryParam("email"))

	// Stock operations
	case "getonestock":
		result, err = h.goldgymSvcStock.GetOneStockProduct(ctx,
			c.QueryParam("stockcode"),
			c.QueryParam("stockname"),
			c.QueryParam("stockid"))
		log.Printf("testDelivery %+v", result)

	case "getallstock":
		result, err = h.goldgymSvcStock.GetAllStockHeader(ctx)

	case "getallstockredis":
		result, err = h.goldgymSvcStock.GetAllStockHeaderToRedis(ctx)

	case "getfromfirebase":
		result, err = h.goldgymSvcStock.GetFromFirebase(ctx, c.QueryParam("userid"))

	case "getimages":
		id, _ := strconv.Atoi(c.QueryParam("id"))
		result, err = h.goldgymSvc.GetTestingImage(ctx, id)

		// Type assertion
		imgData, ok := result.([]byte)
		if !ok {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "invalid image data",
			})
		}

		// Create buffer from image data
		imgBuffer := bytes.NewReader(imgData)

		// Decode image
		img, _, err := image.Decode(imgBuffer)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Unable to decode image",
			})
		}

		// Set header for PNG image
		c.Response().Header().Set("Content-Type", "image/png")

		// Encode image to PNG and write to response
		if err := png.Encode(c.Response().Writer, img); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "unable to encode image",
			})
		}
		return nil
	}

	if err != nil {
		switch {
		case errors.Is(err, entity.ErrNotFound):
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		case errors.Is(err, entity.ErrInvalid):
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		case errors.Is(err, entity.ErrUnauthorized):
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
		}
	}

	resp.Data = result
	resp.Metadata = metadata

	log.Printf("[INFO] %s %s\n", c.Request().Method, c.Request().URL)
	h.logger.For(ctx).Info(
		"HTTP request done",
		zap.String("method", c.Request().Method),
		zap.String("url", c.Request().URL.String()),
	)

	return c.JSON(http.StatusOK, resp)
}
