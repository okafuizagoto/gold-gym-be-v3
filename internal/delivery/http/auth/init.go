package goldgym

import (
	"context"
	"gold-gym-be/internal/entity/auth/v2"
	jaegerLog "gold-gym-be/pkg/log"

	"github.com/opentracing/opentracing-go"
)

type IgoldgymSvc interface {
	LoginUser(ctx context.Context, _user, _password string, _host string) (auth.Token, map[string]interface{}, error)
}

type Handler struct {
	goldgymSvc      IgoldgymSvc
	tracer          opentracing.Tracer
	logger          jaegerLog.Factory
}

// New for bridging product handler initialization
func New(is IgoldgymSvc, tracer opentracing.Tracer, logger jaegerLog.Factory) *Handler {
	return &Handler{
		goldgymSvc:      is,
		tracer:          tracer,
		logger:          logger,
	}
}
