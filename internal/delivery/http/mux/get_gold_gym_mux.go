package goldgym

import (
	httpHelper "gold-gym-be/internal/delivery/http"
	"gold-gym-be/pkg/response"
	"log"
	"net/http"
	"strconv"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"go.uber.org/zap"
)

// Getgoldgym godoc
// @Summary Get entries of all goldgyms
// @Description Get entries of all goldgyms
// @Tags goldgym
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Success 200

// @Router /v1/profiles [get]
func (h *Handler) GetGoldGymMux(w http.ResponseWriter, r *http.Request) {
	var (
		result   interface{}
		metadata interface{}
		err      error
		resp     response.Response
		types    string
	)
	defer resp.RenderJSON(w, r)

	spanCtx, _ := h.tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
	span := h.tracer.StartSpan("Getgoldgym", ext.RPCServerOption(spanCtx))
	defer span.Finish()

	ctx := r.Context()
	ctx = opentracing.ContextWithSpan(ctx, span)
	h.logger.For(ctx).Info("HTTP request received", zap.String("method", r.Method), zap.Stringer("url", r.URL))

	// Your code here
	types = r.FormValue("type")
	switch types {
	case "getgoldgym":
		result, err = h.goldgymSvc.GetGoldUser(ctx)
		log.Println("deliverygolduser", result)
	case "golduserbyemail":
		result, err = h.goldgymSvc.GetGoldUserByEmail(ctx, r.FormValue("email"))
	case "allsubscription":
		result, err = h.goldgymSvc.GetAllSubscription(ctx)
	case "getuserandsubsdetail":
		result, err = h.goldgymSvc.GetSubsWithUser(ctx)
	case "gettotalpayment":
		result, err = h.goldgymSvc.GetSubscriptionHeaderTotalHarga(ctx, r.FormValue("email"))
	// stock -----------------------------------------------------------------------------------------------
	case "getonestock":
		result, err = h.goldgymSvcStock.GetOneStockProduct(ctx, r.FormValue("stockcode"), r.FormValue("stockname"), r.FormValue("stockid"))
		// stock -----------------------------------------------------------------------------------------------
		log.Printf("testDelivery %+v", result)
	case "getallstock":
		result, err = h.goldgymSvcStock.GetAllStockHeader(ctx)
		// stock -----------------------------------------------------------------------------------------------
		// log.Printf("testDelivery %+v", result)
	case "getallstockredis":
		result, err = h.goldgymSvcStock.GetAllStockHeaderToRedis(ctx)
	case "getfromfirebase":
		result, err = h.goldgymSvcStock.GetFromFirebase(ctx, r.FormValue("userid"))
	case "getimages":
		id, _ := strconv.Atoi(r.FormValue("id"))
		result, err = h.goldgymSvc.GetTestingImage(ctx, id)

		// // Example image data in []byte (this should be replaced with your actual data source)
		// imgData := []byte{ /* your PNG image data here */ }

		// // Type assertion
		// imgData, ok := result.([]byte)
		// if !ok {
		// 	log.Fatal("The result is not of type []byte")
		// }

		// // Create a buffer from the image data
		// imgBuffer := bytes.NewReader(imgData)

		// // Decode the image data to get the image.Image object
		// img, _, err := image.Decode(imgBuffer)
		// if err != nil {
		// 	http.Error(w, "Unable to decode image", http.StatusInternalServerError)
		// 	return
		// }

		// // Set the appropriate header for PNG image
		// w.Header().Set("Content-Type", "image/png")

		// // Encode the image to PNG format and write it to the response
		// err = png.Encode(w, img)
		// if err != nil {
		// 	http.Error(w, "Unable to encode image", http.StatusInternalServerError)
		// 	return
		// }
	}

	if err != nil {
		resp = httpHelper.ParseErrorCode(err.Error())

		log.Printf("[ERROR] %s %s - %v\n", r.Method, r.URL, err)
		h.logger.For(ctx).Error("HTTP request error", zap.String("method", r.Method), zap.Stringer("url", r.URL), zap.Error(err))
		return
	}

	resp.Data = result
	resp.Metadata = metadata
	log.Printf("[INFO] %s %s\n", r.Method, r.URL)
	h.logger.For(ctx).Info("HTTP request done", zap.String("method", r.Method), zap.Stringer("url", r.URL))

	return
}
