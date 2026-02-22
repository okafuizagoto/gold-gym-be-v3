package echo

import (
	"fmt"
	"gold-gym-be/internal/entity/firebase"
	goldEntity "gold-gym-be/internal/entity/goldgym"
	goldStockEntity "gold-gym-be/internal/entity/stock"
	"gold-gym-be/pkg/response"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go.uber.org/zap"
)

func (h *Handler) InsertGoldGymEcho(c echo.Context) error {
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
	case "insertuser":
		if err := c.Bind(&insertgolduser); err != nil {
			return c.JSON(400, map[string]string{"error": err.Error()})
		}
		result, err = h.goldgymSvc.InsertGoldUser(ctx, insertgolduser)

	case "insertuserfirebase":
		if err := c.Bind(&insertUserFirebase); err != nil {
			return c.JSON(400, map[string]string{"error": err.Error()})
		}
		result, err = h.goldgymSvcStock.CreateUser(ctx, insertUserFirebase)

	case "loginuser":
		if err := c.Bind(&insertgoldloginuser); err != nil {
			return c.JSON(400, map[string]string{"error": err.Error()})
		}
		host := c.Request().Host
		result, metadata, err = h.goldgymSvc.LoginUser(ctx, insertgoldloginuser.GoldEmail, insertgoldloginuser.GoldPassword, host)
		if err != nil {
			log.Println("err", err)
		}
		fmt.Println("result :", result)
		fmt.Println("metadata :", metadata)
		fmt.Println("err :", err)
		fmt.Println("Result :", insertgoldloginuser)

	case "insertsubsuser":
		if err := c.Bind(&insertgoldsubsuser); err != nil {
			return c.JSON(400, map[string]string{"error": err.Error()})
		}
		result, err = h.goldgymSvc.InsertSubscriptionUser(ctx, insertgoldsubsuser)

	case "insertsubsuserdetail":
		if err := c.Bind(&insertgoldsubsuserdetail); err != nil {
			return c.JSON(400, map[string]string{"error": err.Error()})
		}
		result, err, metadata = h.goldgymSvc.InsertSubscriptionDetail(ctx, insertgoldsubsuserdetail)

	case "insertstock":
		if err := c.Bind(&insertstock); err != nil {
			return c.JSON(400, map[string]string{"error": err.Error()})
		}
		result, err = h.goldgymSvcStock.InsertStockSales(ctx, insertstock)

	case "uploadimages":
		file, _, err := c.Request().FormFile("file")
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Unable to get file"})
		}
		defer file.Close()

		imageBytes, err := ioutil.ReadAll(file)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Unable to read file"})
		}

		testings := goldEntity.Testings{
			TestingImages: imageBytes,
		}
		result, err = h.goldgymSvc.UploadTestingImages(ctx, testings)
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
