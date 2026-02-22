package http

import (
	"context"
	"net/http"

	"gold-gym-be/internal/config"
	"gold-gym-be/pkg/grace"
	jaegerLog "gold-gym-be/pkg/log"

	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
	"github.com/rs/cors"
)

// GoldGymHandler ...
type GoldGymHandler interface {
	// GetGoldGym(w http.ResponseWriter, r *http.Request)
	// InsertGoldGym(w http.ResponseWriter, r *http.Request)
	// DeleteGoldGym(w http.ResponseWriter, r *http.Request)
	// UpdateGoldGym(w http.ResponseWriter, r *http.Request)

	GetGoldGymGin(c *gin.Context)
	InsertGoldGymGin(c *gin.Context)
	DeleteGoldGymGin(c *gin.Context)
	UpdateGoldGymGin(c *gin.Context)

	// PrintSelisih(w http.ResponseWriter, r *http.Request)
	// PrintExpiredTerpajang(w http.ResponseWriter, r *http.Request)
	// PrintExpiredTerkumpul(w http.ResponseWriter, r *http.Request)

	// PrintBatch(w http.ResponseWriter, r *http.Request)
	// PrintBatchFull(w http.ResponseWriter, r *http.Request)

	// //Trans Out
	// InsertTransOut(w http.ResponseWriter, r *http.Request)
	// InsertSales(w http.ResponseWriter, r *http.Request)
	// DeleteSalesByPeriod(w http.ResponseWriter, r *http.Request)
	// RemoveSalesByOutcode(w http.ResponseWriter, r *http.Request)
	// InsertBatchData(w http.ResponseWriter, r *http.Request)
}

// AuthHandler ...
type AuthHandler interface {
	// LoginUser(w http.ResponseWriter, r *http.Request)
	LoginUser(c *gin.Context)
}

type MiddlewareHandler interface {
	CheckUniqueRequest(c *gin.Context)
}

type HealthHandler interface {
	Check(c *gin.Context)
}

type EchoGoldGymHandler interface {
	GetGoldGymEcho(c echo.Context) error
	InsertGoldGymEcho(c echo.Context) error
	UpdateGoldGymEcho(c echo.Context) error
	DeleteGoldGymEcho(c echo.Context) error
}

type MuxGoldGymHandler interface {
	GetGoldGymMux(w http.ResponseWriter, r *http.Request)
	InsertGoldGymMux(w http.ResponseWriter, r *http.Request)
	DeleteGoldGymMux(w http.ResponseWriter, r *http.Request)
	UpdateGoldGymMux(w http.ResponseWriter, r *http.Request)
}

// Server ...
type Server struct {
	Goldgym     GoldGymHandler
	Auth        AuthHandler
	Middleware  MiddlewareHandler
	EchoGoldGym EchoGoldGymHandler
	MuxGoldGym  MuxGoldGymHandler

	engine     *gin.Engine
	echoEngine *echo.Echo
	server     *http.Server

	Health HealthHandler

	Logger jaegerLog.Factory
	Config *config.Config
}

// Serve is serving HTTP gracefully on port x ...
// func (s *Server) Serve(port string) error {
// 	handler := cors.AllowAll().Handler(s.Handler())
// 	return grace.Serve(port, handler)
// }

// func (s *Server) Serve(port string) error {
// 	handler := s.Handler()         // Change this to use Gin handler instead of mux
// 	return handler.Run(":" + port) // Gin's way to serve the app
// }

func (s *Server) Serve(port string) error {
	s.engine = s.Handler()

	s.server = &http.Server{
		Addr:    ":" + port,
		Handler: s.engine,
	}

	return s.server.ListenAndServe()
}

func (s *Server) ServeEcho(port string) error {
	s.echoEngine = s.EchoHandler()
	return s.echoEngine.Start(":" + port)
}

func (s *Server) ServeMux(port string) error {
	handler := cors.AllowAll().Handler(s.Handler())
	return grace.Serve(port, handler)
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.server != nil {
		if err := s.server.Shutdown(ctx); err != nil {
			return err
		}
	}

	if s.echoEngine != nil {
		if err := s.echoEngine.Shutdown(ctx); err != nil {
			return err
		}
	}

	return nil
}
