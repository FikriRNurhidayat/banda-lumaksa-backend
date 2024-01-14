package subscription

import (
	"context"
	"time"

	"github.com/fikrirnurhidayat/banda-lumaksa/pkg/exists"
)

type ListSubscriptionsUseCase interface {
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
}
type ListSubscriptionsResult struct {
	Size          uint32
	Page          uint32
	Limit         uint32
	Subscriptions []Subscription
}
type ListSubscriptionsUseCaseImpl struct {
	subscriptionRepository Repository
}

func (u *ListSubscriptionsUseCaseImpl) Call(ctx context.Context, params *ListSubscriptionsParams) (*ListSubscriptionsResult, error) {
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
		Page:          1,
		Limit:         10,
	}, nil
}

func NewListSubscriptionsUseCase() ListSubscriptionsUseCase {
	return &ListSubscriptionsUseCaseImpl{
		subscriptionRepository: nil,
	}
}
