package subscription

import (
	"context"
	"math"
	"time"

	"github.com/fikrirnurhidayat/banda-lumaksa/pkg/exists"
)

type ListSubscriptionsService interface {
	Call(ctx context.Context, params *ListSubscriptionsParams) (*ListSubscriptionsResult, error)
}
type ListSubscriptionsParams struct {
	NameLike    string
	TypeIs      Type
	StartedFrom time.Time
	StartedTo   time.Time
	EndedFrom   time.Time
	EndedTo     time.Time
	DueFrom     time.Time
	DueTo       time.Time
	CreatedFrom time.Time
	CreatedTo   time.Time
	Page        uint32
	PageSize    uint32
}

type ListSubscriptionsResult struct {
	Size          uint32
	Page          uint32
	PageSize      uint32
	PageCount     uint32
	Subscriptions []Subscription
}

type ListSubscriptionsServiceImpl struct {
	subscriptionRepository Repository
}

func (u *ListSubscriptionsServiceImpl) Call(ctx context.Context, params *ListSubscriptionsParams) (*ListSubscriptionsResult, error) {
	specs := []Specification{}

	if exists.String(params.NameLike) {
		specs = append(specs, NameLike(params.NameLike))
	}

	if params.TypeIs != -1 {
		specs = append(specs, TypeIs(params.TypeIs))
	}

	if exists.Date(params.StartedFrom) && exists.Date(params.StartedTo) {
		specs = append(specs, StartedBetween(params.StartedFrom, params.StartedTo))
	}

	if exists.Date(params.EndedFrom) && exists.Date(params.EndedTo) {
		specs = append(specs, EndedBetween(params.EndedFrom, params.EndedTo))
	}

	if exists.Date(params.CreatedFrom) && exists.Date(params.CreatedTo) {
		specs = append(specs, CreatedBetween(params.CreatedFrom, params.CreatedTo))
	}

	if exists.Date(params.DueFrom) && exists.Date(params.DueTo) {
		specs = append(specs, DueBetween(params.DueFrom, params.DueTo))
	}

	if !exists.Number(params.Page) {
		params.Page = 1
	}

	if !exists.Number(params.PageSize) {
		params.PageSize = 10
	}

	specs = append(specs, Limit(params.PageSize))
	specs = append(specs, Offset((params.Page-1)*params.PageSize))

	subs, err := u.subscriptionRepository.List(ctx, specs...)
	if err != nil {
		return nil, err
	}

	size, err := u.subscriptionRepository.Size(ctx, specs...)
	if err != nil {
		return nil, err
	}

	return &ListSubscriptionsResult{
		Subscriptions: subs,
		Size:          size,
		Page:          params.Page,
		PageSize:      params.PageSize,
		PageCount:     uint32(math.Ceil(float64(size) / float64(params.PageSize))),
	}, nil
}

func NewListSubscriptionsService(subscriptionRepository Repository) ListSubscriptionsService {
	return &ListSubscriptionsServiceImpl{
		subscriptionRepository: subscriptionRepository,
	}
}
