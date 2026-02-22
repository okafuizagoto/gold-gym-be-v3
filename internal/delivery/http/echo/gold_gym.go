package echo

import (
	"context"
	jaegerLog "gold-gym-be/pkg/log"
	"gold-gym-be/pkg/response"

	authEntity "gold-gym-be/internal/entity/auth/v2"
	firebaseEntity "gold-gym-be/internal/entity/firebase"
	goldEntity "gold-gym-be/internal/entity/goldgym"
	goldStockEntity "gold-gym-be/internal/entity/stock"

	"github.com/opentracing/opentracing-go"
)

type IgoldgymSvc interface {
	GetGoldUser(ctx context.Context) ([]goldEntity.GetGoldUser, error)
	InsertGoldUser(ctx context.Context, user goldEntity.GetGoldUsers) (interface{}, error)
	GetGoldUserByEmail(ctx context.Context, email string) (string, error)
	LoginUser(ctx context.Context, _user, _password string, _host string) (authEntity.Token, map[string]interface{}, error)
	GetAllSubscription(ctx context.Context) ([]goldEntity.Subscription, error)
	InsertSubscriptionUser(ctx context.Context, subs goldEntity.InsertSubsAll) (string, error)
	DeleteSubscriptionHeader(ctx context.Context, subs goldEntity.DeleteSubs) (string, error)
	UpdateSubscriptionDetail(ctx context.Context, subs goldEntity.UpdateSubs) (string, error)
	UpdateDataPeserta(ctx context.Context, subs goldEntity.UpdatePassword) (string, error)
	UpdateNama(ctx context.Context, subs goldEntity.UpdateNama) (string, error)
	UpdateKartu(ctx context.Context, subs goldEntity.UpdateKartu) (string, error)
	Logout(ctx context.Context, subs goldEntity.Logout) (string, error)
	GetSubsWithUser(ctx context.Context) ([]goldEntity.GetSubsWithUser, error)
	UpdateValidationOTP(ctx context.Context, otp string, email string) (string, error)
	UpdateOTP(ctx context.Context, email string) (string, error)
	InsertSubscriptionDetail(ctx context.Context, user goldEntity.SubscriptionDetail) (string, error, response.Response)
	UpdateOTPSubscription(ctx context.Context, email string) (string, error)
	UpdatePayment(ctx context.Context, otp string, email string) (string, error, response.Response)
	GetSubscriptionHeaderTotalHarga(ctx context.Context, email string) (goldEntity.SubscriptionHeaderPayment, error)
	UploadTestingImages(ctx context.Context, testing goldEntity.Testings) (string, error)
	GetTestingImage(ctx context.Context, id int) ([]byte, error)
}

type IgoldgymSvcStock interface {
	GetOneStockProduct(ctx context.Context, stockcode string, stocknmame string, stockid string) (goldStockEntity.GetOneStock, error)
	InsertStockSales(ctx context.Context, stock goldStockEntity.InsertStockData) (string, error)
	GetAllStockHeader(ctx context.Context) ([]goldStockEntity.GetOneStock, error)
	GetAllStockHeaderToRedis(ctx context.Context) ([]goldStockEntity.GetOneStock, error)
	GetFromFirebase(ctx context.Context, userID string) (*firebaseEntity.User, error)
	CreateUser(ctx context.Context, user firebaseEntity.User) (string, error)
}

type Handler struct {
	goldgymSvc      IgoldgymSvc
	goldgymSvcStock IgoldgymSvcStock
	tracer          opentracing.Tracer
	logger          jaegerLog.Factory
}

func New(is IgoldgymSvc, isst IgoldgymSvcStock, tracer opentracing.Tracer, logger jaegerLog.Factory) *Handler {
	return &Handler{
		goldgymSvc:      is,
		goldgymSvcStock: isst,
		tracer:          tracer,
		logger:          logger,
	}
}
