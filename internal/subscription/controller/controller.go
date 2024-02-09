package subscription_controller

import (
	"net/http"
	"time"

	common_schema "github.com/fikrirnurhidayat/banda-lumaksa/internal/common/schema"
	common_service "github.com/fikrirnurhidayat/banda-lumaksa/internal/common/service"
	"github.com/fikrirnurhidayat/banda-lumaksa/internal/infra/logger"
	"github.com/google/uuid"
	echo "github.com/labstack/echo/v4"

	subscription_service "github.com/fikrirnurhidayat/banda-lumaksa/internal/subscription/service"
	subscription_types "github.com/fikrirnurhidayat/banda-lumaksa/internal/subscription/types"
)

type SubscriptionController interface {
	Register(*echo.Echo)
	CreateSubscription(c echo.Context) error
	CancelSubscription(c echo.Context) error
	GetSubscription(c echo.Context) error
	ListSubscriptions(c echo.Context) error
}

type SubscriptionControllerImpl struct {
	logger logger.Logger
	subscriptionService subscription_service.SubscriptionService
}

func (ctl *SubscriptionControllerImpl) Register(e *echo.Echo) {
	e.POST("/v1/subscriptions", ctl.CreateSubscription)
	e.DELETE("/v1/subscriptions/:id", ctl.CancelSubscription)
	e.GET("/v1/subscriptions/:id", ctl.GetSubscription)
	e.GET("/v1/subscriptions", ctl.ListSubscriptions)
}

func (ctl *SubscriptionControllerImpl) CancelSubscription(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.Logger().Error(err.Error())
		return err
	}

	params := &subscription_service.CancelSubscriptionParams{
		ID: id,
	}

	if _, err := ctl.subscriptionService.CancelSubscription(c.Request().Context(), params); err != nil {
		c.Logger().Error(err.Error())
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (ctl *SubscriptionControllerImpl) CreateSubscription(c echo.Context) error {
	requestJSON := &CreateSubscriptionRequest{}

	if err := c.Bind(&requestJSON); err != nil {
		c.Logger().Error(err.Error())
		return err
	}

	result, err := ctl.subscriptionService.CreateSubscription(c.Request().Context(), &subscription_service.CreateSubscriptionParams{
		Name:      requestJSON.Subscription.Name,
		Fee:       requestJSON.Subscription.Fee,
		Type:      subscription_types.GetType(requestJSON.Subscription.Type),
		StartedAt: requestJSON.Subscription.StartedAt,
		EndedAt:   requestJSON.Subscription.EndedAt,
		DueAt:     requestJSON.Subscription.DueAt,
	})

	if err != nil {
		c.Logger().Error(err.Error())
		return err
	}

	response := &CreateSubscriptionResponse{
		Subscription: NewSubscriptionResponse(result.Subscription),
	}

	return c.JSON(http.StatusCreated, response)
}

func (ctl *SubscriptionControllerImpl) GetSubscription(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.Logger().Error(err.Error())
		return err
	}

	params := &subscription_service.GetSubscriptionParams{
		ID: id,
	}

	result, err := ctl.subscriptionService.GetSubscription(c.Request().Context(), params)
	if err != nil {
		c.Logger().Error(err.Error())
		return err
	}

	response := &GetSubscriptionResponse{
		Subscription: NewSubscriptionResponse(result.Subscription),
	}

	return c.JSON(http.StatusOK, response)
}

func (ctl *SubscriptionControllerImpl) ListSubscriptions(c echo.Context) error {
	params := &subscription_service.ListSubscriptionsParams{
		TypeIs: -1,
		Pagination: common_service.PaginationParams{},
	}

	if err := echo.QueryParamsBinder(c).
		String("name_like", &params.NameLike).
		Uint32("page", &params.Pagination.Page).
		Uint32("page_size", &params.Pagination.PageSize).
		Time("created_from", &params.CreatedFrom, time.RFC3339).
		Time("created_to", &params.CreatedFrom, time.RFC3339).
		Time("started_from", &params.StartedFrom, time.RFC3339).
		Time("started_to", &params.StartedFrom, time.RFC3339).
		Time("ended_from", &params.EndedFrom, time.RFC3339).
		Time("ended_to", &params.EndedFrom, time.RFC3339).
		Time("due_from", &params.DueFrom, time.RFC3339).
		Time("due_to", &params.DueFrom, time.RFC3339).
		CustomFunc("type_is", func(values []string) []error {
			params.TypeIs = subscription_types.GetType(values[0])

			return nil
		}).
		FailFast(true).
		BindError(); err != nil {
		c.Logger().Error(err.Error())
		return err
	}

	result, err := ctl.subscriptionService.ListSubscriptions(c.Request().Context(), params)
	if err != nil {
		c.Logger().Error(err.Error())
		return err
	}

	response := &ListSubscriptionsResponse{
		PaginationResponse: common_schema.NewPaginationResponse(result.Pagination),
		Subscriptions:      NewSubscriptionsResponse(result.Subscriptions),
	}

	return c.JSON(http.StatusOK, response)
}

func New(logger logger.Logger, subscriptionService subscription_service.SubscriptionService) SubscriptionController {
	return &SubscriptionControllerImpl{
		logger:              logger,
		subscriptionService: subscriptionService,
	}
}
