package subscription

import (
	"fmt"
	"net/http"
	"reflect"
	"time"

	echo "github.com/labstack/echo/v4"
)

type Controller interface {
	CreateSubscription(ctx echo.Context) error
	CancelSubscription(ctx echo.Context) error
	GetSubscription(ctx echo.Context) error
	ListSubscriptions(ctx echo.Context) error
}

type ControllerImpl struct {
	CreateSubscriptionUseCase CreateSubscriptionUseCase
	CancelSubscriptionUseCase CancelSubscriptionUseCase
	GetSubscriptionUseCase    GetSubscriptionUseCase
	ListSubscriptionsUseCase  ListSubscriptionsUseCase
}

// CancelSubscription implements Controller.
func (*ControllerImpl) CancelSubscription(ctx echo.Context) error {
	panic("unimplemented")
}

// CreateSubscription implements Controller.
func (*ControllerImpl) CreateSubscription(ctx echo.Context) error {
	panic("unimplemented")
}

// GetSubscription implements Controller.
func (*ControllerImpl) GetSubscription(ctx echo.Context) error {
	panic("unimplemented")
}

func (ctl *ControllerImpl) ListSubscriptions(c echo.Context) error {
	params := &ListSubscriptionsParams{}

	if err := echo.QueryParamsBinder(c).
		String("name_like", &params.NameLike).
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
		return err
	}

	result, err := ctl.ListSubscriptionsUseCase.Call(c.Request().Context(), params)
	if err != nil {
		return err
	}

	response := &ListSubscriptionsResponse{
		Page:          result.Page,
		PageCount:     0,
		PageSize:      result.Limit,
		Size:          result.Size,
		Subscriptions: SubscriptionsResponseFromSubscriptions(result.Subscriptions),
	}

	return c.JSON(http.StatusOK, response)
}

func NewController(sets ...func(*ControllerImpl)) Controller {
	ctl := &ControllerImpl{}

	for _, set := range sets {
		set(ctl)
	}

	return ctl
}

func With[Mod any, Deps any](field string, deps Deps) func(mod Mod) {
	return func(mod Mod) {
		depValue := reflect.ValueOf(deps)
		modValue := reflect.ValueOf(mod).Elem()
		if modValue.Kind() == reflect.Struct {
			modField := modValue.FieldByName(field)
			if modField.IsValid() && modField.CanSet() {
				modField.Set(depValue)
			}
		}
	}
}
