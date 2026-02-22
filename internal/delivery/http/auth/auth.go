package goldgym

import (
	"errors"
	"gold-gym-be/pkg/response"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
func (h *Handler) LoginUser(c *gin.Context) {
	// func (h *Handler) DeleteGoldGymGin(c *gin.Context) {
	log.Println("testDelivery")
	resp := response.Response{}
	// defer resp.RenderJSON(w, r)

	// ctx := c.Request.Context()
	ctx := c.Request.Context()

	user, password, ok := c.Request.BasicAuth()
	if !ok {
		log.Printf("[ERROR] %s %s - %s\n", c.Request.Method, c.Request.URL, errors.New("403 Forbidden"))
		return
	}
	log.Println("testDelivery2")

	result, metadata, err := h.goldgymSvc.LoginUser(ctx, user, password, c.Request.RemoteAddr)
	log.Println("testDelivery3")
	if err != nil {
		// Return error message with HTTP 200 OK
		resp.SetError(err, http.StatusOK)

		log.Printf("[ERROR] %s %s - %s\n", c.Request.Method, c.Request.URL, err.Error())
		return
	}

	resp.Data = result
	resp.Metadata = metadata

	log.Printf("[INFO] %s %s\n", c.Request.Method, c.Request.URL)
}
