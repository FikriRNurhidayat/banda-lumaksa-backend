package subscription_service

import (
	"context"
	"time"

	"github.com/fikrirnurhidayat/banda-lumaksa/pkg/exists"

	common_repository "github.com/fikrirnurhidayat/banda-lumaksa/internal/common/repository"
	common_service "github.com/fikrirnurhidayat/banda-lumaksa/internal/common/service"
	common_specification "github.com/fikrirnurhidayat/banda-lumaksa/internal/common/specification"
	subscription_entity "github.com/fikrirnurhidayat/banda-lumaksa/internal/domain/subscription/entity"
	subscription_specification "github.com/fikrirnurhidayat/banda-lumaksa/internal/domain/subscription/specification"
	subscription_types "github.com/fikrirnurhidayat/banda-lumaksa/internal/domain/subscription/types"
)

type ListSubscriptionsParams struct {
	NameLike    string
	TypeIs      subscription_types.Type
	StartedFrom time.Time
	StartedTo   time.Time
	EndedFrom   time.Time
	EndedTo     time.Time
	DueFrom     time.Time
	DueTo       time.Time
	CreatedFrom time.Time
	CreatedTo   time.Time
	Pagination  common_service.PaginationParams
}

type ListSubscriptionsResult struct {
	Pagination    common_service.PaginationResult
	Subscriptions []subscription_entity.Subscription
}

func (s *SubscriptionServiceImpl) ListSubscriptions(ctx context.Context, params *ListSubscriptionsParams) (*ListSubscriptionsResult, error) {
	filters := []subscription_specification.SubscriptionSpecification{}

	if exists.String(params.NameLike) {
		filters = append(filters, subscription_specification.NameLike(params.NameLike))
	}

	if params.TypeIs != -1 {
		filters = append(filters, subscription_specification.TypeIs(params.TypeIs))
	}

	if exists.Date(params.StartedFrom) && exists.Date(params.StartedTo) {
		filters = append(filters, subscription_specification.StartedBetween(params.StartedFrom, params.StartedTo))
	}

	if exists.Date(params.EndedFrom) && exists.Date(params.EndedTo) {
		filters = append(filters, subscription_specification.EndedBetween(params.EndedFrom, params.EndedTo))
	}

	if exists.Date(params.CreatedFrom) && exists.Date(params.CreatedTo) {
		filters = append(filters, subscription_specification.CreatedBetween(params.CreatedFrom, params.CreatedTo))
	}

	if exists.Date(params.DueFrom) && exists.Date(params.DueTo) {
		filters = append(filters, subscription_specification.DueBetween(params.DueFrom, params.DueTo))
	}

	params.Pagination = params.Pagination.Normalize()

	subs, err := s.subscriptionRepository.List(ctx, common_repository.ListArgs[subscription_specification.SubscriptionSpecification]{
		Filters: filters,
		Limit:   common_specification.WithLimit(params.Pagination.Limit()),
		Offset:  common_specification.WithOffset(params.Pagination.Offset()),
	})
	if err != nil {
		s.logger.Error("subscription repository list error", "detail", err.Error())
		return nil, err
	}

	size, err := s.subscriptionRepository.Size(ctx, filters...)
	if err != nil {
		s.logger.Error("subscription repository size error", "detail", err.Error())
		return nil, err
	}

	return &ListSubscriptionsResult{
		Subscriptions: subs,
		Pagination:    common_service.NewPaginationResult(params.Pagination, size),
	}, nil
}
