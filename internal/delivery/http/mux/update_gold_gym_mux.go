package goldgym

import (
	"encoding/json"
	httpHelper "gold-gym-be/internal/delivery/http"
	"gold-gym-be/pkg/response"
	"io/ioutil"
	"log"
	"net/http"

	goldEntity "gold-gym-be/internal/entity/goldgym"

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
func (h *Handler) UpdateGoldGymMux(w http.ResponseWriter, r *http.Request) {
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
	case "updatesubsuser":
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &updategoldsubsuser)
		result, err = h.goldgymSvc.UpdateSubscriptionDetail(ctx, updategoldsubsuser)
		if err != nil {
			log.Println("err", err)
		}
	case "updatepassword":
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &updatepassword)
		result, err = h.goldgymSvc.UpdateDataPeserta(ctx, updatepassword)
		if err != nil {
			log.Println("err", err)
		}
	case "updatenama":
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &updatenama)
		result, err = h.goldgymSvc.UpdateNama(ctx, updatenama)
		if err != nil {
			log.Println("err", err)
		}
	case "updatekartu":
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &updatekartu)
		result, err = h.goldgymSvc.UpdateKartu(ctx, updatekartu)
		if err != nil {
			log.Println("err", err)
		}
	case "logout":
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &logout)
		result, err = h.goldgymSvc.Logout(ctx, logout)
		if err != nil {
			log.Println("err", err)
		}
	case "updatevalidationemail":
		// body, _ := ioutil.ReadAll(r.Body)
		// json.Unmarshal(body, &logout)
		// fmt.Println("Result :", logout)
		// fmt.Println("Result2 :", &logout)
		result, err = h.goldgymSvc.UpdateValidationOTP(ctx, r.FormValue("otp"), r.FormValue("email"))
		if err != nil {
			log.Println("err", err)
		}
	case "updateotp":
		// body, _ := ioutil.ReadAll(r.Body)
		// json.Unmarshal(body, &logout)
		// fmt.Println("Result :", logout)
		// fmt.Println("Result2 :", &logout)
		result, err = h.goldgymSvc.UpdateOTP(ctx, r.FormValue("email"))
		if err != nil {
			log.Println("err", err)
		}
	case "updateotpsubscription":
		// body, _ := ioutil.ReadAll(r.Body)
		// json.Unmarshal(body, &logout)
		// fmt.Println("Result :", logout)
		// fmt.Println("Result2 :", &logout)
		// result, _, err = h.goldgymSvc.UpdateOTPSubscription(ctx, r.FormValue("email"))
		result, err = h.goldgymSvc.UpdateOTPSubscription(ctx, r.FormValue("email"))
		if err != nil {
			log.Println("err", err)
		}
	case "updatepaymentsubscription":
		result, err, resp = h.goldgymSvc.UpdatePayment(ctx, r.FormValue("otp"), r.FormValue("email"))
		// case "":
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
