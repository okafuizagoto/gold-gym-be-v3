package goldgym

import (
	"encoding/json"
	"gold-gym-be/pkg/response"
	"io/ioutil"
	"log"
	"net/http"

	"gold-gym-be/internal/entity/firebase"
	goldEntity "gold-gym-be/internal/entity/goldgym"
	goldStockEntity "gold-gym-be/internal/entity/stock"

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
func (h *Handler) InsertGoldGymMux(w http.ResponseWriter, r *http.Request) {
	var (
		result   interface{}
		metadata interface{}
		err      error

		resp           response.Response
		types          string
		insertgolduser goldEntity.GetGoldUsers
		// insertgoldloginuser      goldEntity.LogUser
		insertgoldsubsuser       goldEntity.InsertSubsAll
		insertUserFirebase       firebase.User
		insertgoldsubsuserdetail goldEntity.SubscriptionDetail
		insertstock              goldStockEntity.InsertStockData
		// header                   http.Header
		// testings                 goldEntity.Testings
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
	case "insertuser":
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &insertgolduser)
		result, err = h.goldgymSvc.InsertGoldUser(ctx, insertgolduser)
	case "insertuserfirebase":
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &insertUserFirebase)
		result, err = h.goldgymSvcStock.CreateUser(ctx, insertUserFirebase)
	// case "loginuser":
	// 	body, _ := ioutil.ReadAll(r.Body)
	// 	json.Unmarshal(body, &insertgoldloginuser)
	// 	fmt.Println("Result :", insertgoldloginuser)
	// 	fmt.Println("Result2 :", &insertgoldloginuser)
	// 	result, metadata, err = h.goldgymSvc.LoginUser(ctx, insertgoldloginuser)
	// 	if err != nil {
	// 		log.Println("err", err)
	// 	}
	// 	fmt.Println("result :", result)
	// 	fmt.Println("metadata :", metadata)
	// 	fmt.Println("err :", err)
	// 	fmt.Println("Result :", insertgoldloginuser)
	case "insertsubsuser":
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &insertgoldsubsuser)
		result, err = h.goldgymSvc.InsertSubscriptionUser(ctx, insertgoldsubsuser)
		if err != nil {
			log.Println("err", err)
		}
	case "insertsubsuserdetail":
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &insertgoldsubsuserdetail)
		result, err, resp = h.goldgymSvc.InsertSubscriptionDetail(ctx, insertgoldsubsuserdetail)
		if err != nil {
			log.Println("err", err)
		}
	case "insertstock":
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &insertstock)
		result, err = h.goldgymSvcStock.InsertStockSales(ctx, insertstock)
		if err != nil {
			log.Println("err", err)
		}
		// case "":
	case "uploadimages":

		// err := r.ParseMultipartForm(10 << 20) // 10 MB limit
		// if err != nil {
		// 	http.Error(w, "Unable to parse form", http.StatusBadRequest)
		// 	return
		// }
		// -- body, _ := ioutil.ReadAll(r.Body)
		// Retrieve the file from the form field named "image"
		file, _, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Unable to get file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			http.Error(w, "Unable to read file", http.StatusInternalServerError)
			return
		}

		// fmt.Println("fileBytes", fileBytes)

		// fmt.Println("files2", file)

		test := goldEntity.Testings{
			ID:            r.FormValue("id"),
			TestingImages: fileBytes,
		}

		result, err = h.goldgymSvc.UploadTestingImages(ctx, test)
		if err != nil {
			log.Println("err", err)
		}
	}

	if err != nil {
		resp.SetError(err, http.StatusInternalServerError)
		resp.StatusCode = 500
		log.Printf("[ERROR] %s %s - %s\n", r.Method, r.URL, err.Error())
		return
	}

	resp.Data = result
	resp.Metadata = metadata
	log.Printf("[INFO] %s %s\n", r.Method, r.URL)
	h.logger.For(ctx).Info("HTTP request done", zap.String("method", r.Method), zap.Stringer("url", r.URL))

	return
}
