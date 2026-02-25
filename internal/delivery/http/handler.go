package http

import (
	"context"
	"errors"
	"gold-gym-be/internal/delivery/http/middleware"
	"gold-gym-be/pkg/response"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"

	beegoWeb "github.com/beego/beego/v2/server/web"
)

// // Handler will initialize mux router and register handler
// func (s *Server) Handler() *mux.Router {
// 	r := mux.NewRouter()
// 	// Jika tidak ditemukan, jangan diubah.
// 	r.NotFoundHandler = http.HandlerFunc(notFoundHandler)
// 	// Health Check
// 	r.HandleFunc("", defaultHandler).Methods("GET")
// 	r.HandleFunc("/", defaultHandler).Methods("GET")

// 	// Tambahan Prefix di depan API endpoint
// 	router := r.PathPrefix("/gold-gym").Subrouter()

// 	router.HandleFunc("", defaultHandler).Methods("GET")
// 	router.HandleFunc("/", defaultHandler).Methods("GET")

// 	sub := router.PathPrefix("/v2").Subrouter()

// 	// Routes
// 	goldgym := sub.PathPrefix("/userdata").Subrouter()

// 	goldgym.HandleFunc("", s.Goldgym.GetGoldGym).Methods("GET")
// 	goldgym.HandleFunc("", s.Goldgym.InsertGoldGym).Methods("POST")
// 	goldgym.HandleFunc("", s.Goldgym.UpdateGoldGym).Methods("PUT")
// 	goldgym.HandleFunc("", s.Goldgym.DeleteGoldGym).Methods("DELETE")

// 	goldgym.HandleFunc("/login", s.Auth.LoginUser).Methods("POST")

// 	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
// 	return r
// }

// func defaultHandler(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte("Example Service API"))
// }

// Handler will initialize Gin router and register handler
func (s *Server) Handler() *gin.Engine {
	r := gin.New()

	// recovery
	if s.Config.Server.Env == "local" {
		r.Use(gin.Recovery())
	} else {
		// r.Use(gin.CustomRecovery(func(c gin.Context, recovered interface{}) {
		// 	s.Logger.For(context.Background()).Error(
		// 		"panic",
		// 		zap.Any("err", recovered),
		// 	)
		// 	c.JSON(500, gin.H{"error": "internal server error"})
		// }))
		r.Use(func(c *gin.Context) {
			defer func() {
				if err := recover(); err != nil {
					log.Printf("panic: %v", err)
					c.JSON(500, gin.H{"error": "internal error"})
				}
			}()
			c.Next()
		})
	}

	// metrics + access log
	r.Use(middleware.PrometheusMetrics())
	r.Use(middleware.AccessLogger())

	// timeout
	r.Use(middleware.Timeout(5 * time.Second))

	// r.Use(middleware.Timeout(5 * time.Second))
	// r.Use(func(c *gin.Context) {
	// 	defer func() {
	// 		if err := recover(); err != nil {
	// 			log.Printf("panic: %v", err)
	// 			c.JSON(500, gin.H{"error": "internal error"})
	// 		}
	// 	}()
	// 	c.Next()
	// })
	r.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		s.Logger.For(context.Background()).Error(
			"panic recovered",
			zap.Any("err", recovered),
		)

		c.JSON(500, gin.H{
			"error": "internal server error",
		})
	}))
	// Health Check
	// r.GET("/", defaultHandler)
	r.GET("", defaultHandler)
	r.GET("/healthz", s.Health.Check)

	// Tambahan Prefix di depan API endpoint
	router := r.Group("/gold-gym")

	// Routes
	goldgym := router.Group("/v2/userdata")
	{
		// Define the routes for GoldGym
		goldgym.GET("", s.Goldgym.GetGoldGymGin)                                      // GET
		goldgym.POST("", s.Middleware.CheckUniqueRequest, s.Goldgym.InsertGoldGymGin) // POST
		goldgym.PUT("", s.Goldgym.UpdateGoldGymGin)                                   // PUT
		goldgym.DELETE("", s.Goldgym.DeleteGoldGymGin)                                // DELETE

		// Auth routes
		goldgym.POST("/login", s.Auth.LoginUser) // POST
	}

	// Elastic routes
	elastic := router.Group("/v2/elastic")
	{
		elastic.GET("", s.Elastic.GetElasticGin)   // GET: search or getbyid
		elastic.POST("", s.Elastic.PostElasticGin) // POST: index document
	}

	// Swagger route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Prometheus metrics endpoint — scraped by Prometheus server
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	return r
}

func (s *Server) EchoHandler() *echo.Echo {
	e := echo.New()

	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())

	echoGym := e.Group("/echo-gym")
	echoUserdata := echoGym.Group("/v2/userdata")
	{
		echoUserdata.GET("", s.EchoGoldGym.GetGoldGymEcho)       // GET
		echoUserdata.POST("", s.EchoGoldGym.InsertGoldGymEcho)   // POST
		echoUserdata.PUT("", s.EchoGoldGym.UpdateGoldGymEcho)    // PUT
		echoUserdata.DELETE("", s.EchoGoldGym.DeleteGoldGymEcho) // DELETE
	}

	return e
}

func (s *Server) MuxHandler() *mux.Router {
	r := mux.NewRouter()
	// Jika tidak ditemukan, jangan diubah.
	r.NotFoundHandler = http.HandlerFunc(notFoundHandler)
	// Health Check
	r.HandleFunc("", defaultHandlerMux).Methods("GET")
	r.HandleFunc("/", defaultHandlerMux).Methods("GET")

	// Tambahan Prefix di depan API endpoint
	router := r.PathPrefix("/mux-gold-gym").Subrouter()

	router.HandleFunc("", defaultHandlerMux).Methods("GET")
	router.HandleFunc("/", defaultHandlerMux).Methods("GET")

	sub := router.PathPrefix("/v2").Subrouter()

	// Routes
	goldgym := sub.PathPrefix("/userdata").Subrouter()

	goldgym.HandleFunc("", s.MuxGoldGym.GetGoldGymMux).Methods("GET")
	goldgym.HandleFunc("", s.MuxGoldGym.InsertGoldGymMux).Methods("POST")
	goldgym.HandleFunc("", s.MuxGoldGym.UpdateGoldGymMux).Methods("PUT")
	goldgym.HandleFunc("", s.MuxGoldGym.DeleteGoldGymMux).Methods("DELETE")

	// goldgym.HandleFunc("/login", s.Auth.LoginUser).Methods("POST")

	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
	return r
}

func (s *Server) BeegoHandler() *beegoWeb.HttpServer {
	app := beegoWeb.NewHttpSever()

	app.Cfg.WebConfig.AutoRender = false
	app.Cfg.Log.AccessLogs = false

	app.Get("/beego-gym/v2/userdata", s.BeegoGoldGym.GetGoldGymBeego)
	app.Post("/beego-gym/v2/userdata", s.BeegoGoldGym.InsertGoldGymBeego)
	app.Put("/beego-gym/v2/userdata", s.BeegoGoldGym.UpdateGoldGymBeego)
	app.Delete("/beego-gym/v2/userdata", s.BeegoGoldGym.DeleteGoldGymBeego)

	return app
}

func defaultHandlerMux(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Example Service API"))
}

func defaultHandler(c *gin.Context) {
	c.String(200, "Example Service API")
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	var (
		resp   *response.Response
		err    error
		errRes response.Error
	)
	resp = &response.Response{}
	defer resp.RenderJSON(w, r)

	err = errors.New("404 Not Found")

	if err != nil {
		// Error response handling
		errRes = response.Error{
			Code:   404,
			Msg:    "404 Not Found",
			Status: true,
		}

		log.Printf("[ERROR] %s %s - %v\n", r.Method, r.URL, err)
		resp.StatusCode = 404
		resp.Error = errRes
		return
	}
}
