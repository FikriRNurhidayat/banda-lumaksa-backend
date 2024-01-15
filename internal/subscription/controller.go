package subscription

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	echo "github.com/labstack/echo/v4"
)

type Controller interface {
	Register(*echo.Echo)
	CreateSubscription(c echo.Context) error
	CancelSubscription(c echo.Context) error
	GetSubscription(c echo.Context) error
	ListSubscriptions(c echo.Context) error
}

type ControllerImpl struct {
	CreateSubscriptionService CreateSubscriptionService
	CancelSubscriptionService CancelSubscriptionService
	GetSubscriptionService    GetSubscriptionService
	ListSubscriptionsService  ListSubscriptionsService
}

func (ctl *ControllerImpl) Register(e *echo.Echo) {
	e.POST("/v1/subscriptions", ctl.CreateSubscription)
	e.DELETE("/v1/subscriptions/:id", ctl.CancelSubscription)
	e.GET("/v1/subscriptions/:id", ctl.GetSubscription)
	e.GET("/v1/subscriptions", ctl.ListSubscriptions)
}

func (ctl *ControllerImpl) CancelSubscription(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.Logger().Error(err.Error())
		return err
	}

	params := &CancelSubscriptionParams{
		ID: id,
	}

	if _, err := ctl.CancelSubscriptionService.Call(c.Request().Context(), params); err != nil {
		c.Logger().Error(err.Error())
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (ctl *ControllerImpl) CreateSubscription(c echo.Context) error {
	payload := &CreateSubscriptionRequest{}

	if err := c.Bind(&payload); err != nil {
		c.Logger().Error(err.Error())
		return err
	}

	result, err := ctl.CreateSubscriptionService.Call(c.Request().Context(), &CreateSubscriptionParams{
		Name:      payload.Name,
		Fee:       payload.Fee,
		Type:      GetType(payload.Type),
		StartedAt: payload.StartedAt,
		EndedAt:   payload.EndedAt,
		DueAt:     payload.DueAt,
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

func (ctl *ControllerImpl) GetSubscription(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.Logger().Error(err.Error())
		return err
	}

	params := &GetSubscriptionParams{
		ID: id,
	}

	result, err := ctl.GetSubscriptionService.Call(c.Request().Context(), params)
	if err != nil {
		c.Logger().Error(err.Error())
		return err
	}

	response := &GetSubscriptionResponse{
		Subscription: NewSubscriptionResponse(result.Subscription),
	}

	return c.JSON(http.StatusOK, response)
}

func (ctl *ControllerImpl) ListSubscriptions(c echo.Context) error {
	params := &ListSubscriptionsParams{}

	if err := echo.QueryParamsBinder(c).
		String("name_like", &params.NameLike).
		Uint32("page", &params.Page).
		Uint32("page_size", &params.PageSize).
		Time("created_from", &params.CreatedFrom, time.RFC3339).
		Time("created_to", &params.CreatedFrom, time.RFC3339).
		Time("started_from", &params.StartedFrom, time.RFC3339).
		Time("started_to", &params.StartedFrom, time.RFC3339).
		Time("ended_from", &params.EndedFrom, time.RFC3339).
		Time("ended_to", &params.EndedFrom, time.RFC3339).
		Time("due_from", &params.DueFrom, time.RFC3339).
		Time("due_to", &params.DueFrom, time.RFC3339).
		CustomFunc("type_is", func(values []string) []error {
			params.TypeIs = GetType(values[0])

			return nil
		}).
		FailFast(true).
		BindError(); err != nil {
		c.Logger().Error(err.Error())
		return err
	}

	result, err := ctl.ListSubscriptionsService.Call(c.Request().Context(), params)
	if err != nil {
		c.Logger().Error(err.Error())
		return err
	}

	response := &ListSubscriptionsResponse{
		Page:          result.Page,
		PageCount:     result.PageCount,
		PageSize:      result.PageSize,
		Size:          result.Size,
		Subscriptions: NewSubscriptionsResponse(result.Subscriptions),
	}

	return c.JSON(http.StatusOK, response)
}

func NewController(
	listSubscriptionsSvc ListSubscriptionsService,
	getSubscriptionSvc GetSubscriptionService,
	createSubscriptionSvc CreateSubscriptionService,
	cancelSubscriptionSvc CancelSubscriptionService,
) Controller {
	return &ControllerImpl{
		CreateSubscriptionService: createSubscriptionSvc,
		CancelSubscriptionService: cancelSubscriptionSvc,
		GetSubscriptionService:    getSubscriptionSvc,
		ListSubscriptionsService:  listSubscriptionsSvc,
	}
}
