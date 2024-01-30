package subscription

import (
	"context"
	"fmt"
	"time"

	"github.com/fikrirnurhidayat/banda-lumaksa/internal/common/errors"
	"github.com/fikrirnurhidayat/banda-lumaksa/internal/common/manager"
	"github.com/fikrirnurhidayat/banda-lumaksa/internal/common/repository"
	"github.com/fikrirnurhidayat/banda-lumaksa/internal/common/service"
	"github.com/fikrirnurhidayat/banda-lumaksa/internal/common/specification"
	"github.com/fikrirnurhidayat/banda-lumaksa/internal/common/values"
	"github.com/fikrirnurhidayat/banda-lumaksa/internal/transaction"
	"github.com/fikrirnurhidayat/banda-lumaksa/pkg/exists"
	"github.com/google/uuid"
)

type SubscriptionService interface {
	CreateSubscription(ctx context.Context, params *CreateSubscriptionParams) (*CreateSubscriptionResult, error)
	GetSubscription(ctx context.Context, params *GetSubscriptionParams) (*GetSubscriptionResult, error)
	ListSubscriptions(ctx context.Context, params *ListSubscriptionsParams) (*ListSubscriptionsResult, error)
	CancelSubscription(ctx context.Context, params *CancelSubscriptionParams) (*CancelSubscriptionResult, error)
	ChargeSubscription(ctx context.Context, params *ChargeSubscriptionParams) (*ChargeSubscriptionResult, error)
	ChargeSubscriptions(ctx context.Context, params *ChargeSubscriptionsParams) (*ChargeSubscriptionsResult, error)
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
	Pagination  service.PaginationParams
}

type ListSubscriptionsResult struct {
	Pagination    service.PaginationResult
	Subscriptions []Subscription
}

type GetSubscriptionParams struct {
	ID uuid.UUID
}

type GetSubscriptionResult struct {
	Subscription Subscription
}

type CreateSubscriptionParams struct {
	Name      string
	Fee       int32
	Type      Type
	StartedAt time.Time
	EndedAt   time.Time
	DueAt     time.Time
}

type CreateSubscriptionResult struct {
	Subscription Subscription
}

type CancelSubscriptionParams struct {
	ID uuid.UUID
}

type CancelSubscriptionResult struct{}

type ChargeSubscriptionParams struct {
	ID uuid.UUID
}

type ChargeSubscriptionResult struct {
	Subscription Subscription
}

type ChargeSubscriptionsParams struct{}
type ChargeSubscriptionsResult struct{}

type SubscriptionServiceImpl struct {
	subscriptionRepository SubscriptionRepository
	transactionRepository  transaction.TransactionRepository
	transactionManager     manager.TransactionManager
}

func (s *SubscriptionServiceImpl) CancelSubscription(ctx context.Context, params *CancelSubscriptionParams) (*CancelSubscriptionResult, error) {
	subscription, err := s.subscriptionRepository.Get(ctx, WithID(params.ID))
	if err != nil {
		return nil, err
	}

	if err := s.subscriptionRepository.Delete(ctx, WithID(subscription.ID)); err != nil {
		return nil, errors.ErrInternalServer
	}

	return &CancelSubscriptionResult{}, nil
}

func (s *SubscriptionServiceImpl) ChargeSubscriptions(ctx context.Context, params *ChargeSubscriptionsParams) (*ChargeSubscriptionsResult, error) {
	today := time.Now()
	iterator, err := s.subscriptionRepository.Each(ctx, repository.ListArgs[SubscriptionSpecification]{
		Filters: SubscriptionSpecifications{DueBefore(today), NotEnded(today)},
	})
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	for iterator.Next() {
		subscription, err := iterator.Entry()
		if err != nil {
			continue
		}

		if _, err := s.chargeSubscription(ctx, subscription); err != nil {
			continue
		}
	}

	return &ChargeSubscriptionsResult{}, nil
}

func (s *SubscriptionServiceImpl) ChargeSubscription(ctx context.Context, params *ChargeSubscriptionParams) (*ChargeSubscriptionResult, error) {
	subscription, err := s.subscriptionRepository.Get(ctx, WithID(params.ID))
	if err != nil {
		return nil, ErrSubscriptionNotFound
	}

	subscription, err = s.chargeSubscription(ctx, subscription)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	return &ChargeSubscriptionResult{}, nil
}

func (s *SubscriptionServiceImpl) CreateSubscription(ctx context.Context, params *CreateSubscriptionParams) (*CreateSubscriptionResult, error) {
	now := time.Now()
	subscription := Subscription{
		ID:        uuid.New(),
		Name:      params.Name,
		Fee:       params.Fee,
		Type:      params.Type,
		StartedAt: params.StartedAt,
		EndedAt:   params.EndedAt,
		DueAt:     params.DueAt,
	}

	if subscription.DueAt == values.NoTime {
		subscription.DueAt = s.computeDueAt(subscription, params.StartedAt)
	}

	subscription.CreatedAt = now
	subscription.UpdatedAt = now

	if err := s.subscriptionRepository.Save(ctx, subscription); err != nil {
		fmt.Println(err.Error())
		return nil, errors.ErrInternalServer
	}

	return &CreateSubscriptionResult{
		Subscription: subscription,
	}, nil
}

func (s *SubscriptionServiceImpl) GetSubscription(ctx context.Context, params *GetSubscriptionParams) (*GetSubscriptionResult, error) {
	subscription, err := s.subscriptionRepository.Get(ctx, WithID(params.ID))
	if err != nil {
		return nil, err
	}

	if subscription == NoSubscription {
		return nil, ErrSubscriptionNotFound
	}

	return &GetSubscriptionResult{
		Subscription: subscription,
	}, nil
}

func (s *SubscriptionServiceImpl) ListSubscriptions(ctx context.Context, params *ListSubscriptionsParams) (*ListSubscriptionsResult, error) {
	filters := []SubscriptionSpecification{}

	if exists.String(params.NameLike) {
		filters = append(filters, NameLike(params.NameLike))
	}

	if params.TypeIs != -1 {
		filters = append(filters, TypeIs(params.TypeIs))
	}

	if exists.Date(params.StartedFrom) && exists.Date(params.StartedTo) {
		filters = append(filters, StartedBetween(params.StartedFrom, params.StartedTo))
	}

	if exists.Date(params.EndedFrom) && exists.Date(params.EndedTo) {
		filters = append(filters, EndedBetween(params.EndedFrom, params.EndedTo))
	}

	if exists.Date(params.CreatedFrom) && exists.Date(params.CreatedTo) {
		filters = append(filters, CreatedBetween(params.CreatedFrom, params.CreatedTo))
	}

	if exists.Date(params.DueFrom) && exists.Date(params.DueTo) {
		filters = append(filters, DueBetween(params.DueFrom, params.DueTo))
	}

	params.Pagination = params.Pagination.Normalize()

	subs, err := s.subscriptionRepository.List(ctx, repository.ListArgs[SubscriptionSpecification]{
		Filters: filters,
		Limit:   specification.WithLimit(params.Pagination.Limit()),
		Offset:  specification.WithOffset(params.Pagination.Offset()),
	})
	if err != nil {
		fmt.Println(err.Error())
		return nil, errors.ErrInternalServer
	}

	size, err := s.subscriptionRepository.Size(ctx, filters...)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	return &ListSubscriptionsResult{
		Subscriptions: subs,
		Pagination:    service.NewPaginationResult(params.Pagination, size),
	}, nil
}

func (s *SubscriptionServiceImpl) computeDueAt(subscription Subscription, startFrom time.Time) time.Time {
	switch subscription.Type {
	case Daily:
		return startFrom.Add(values.Day)
	case Weekly:
		return startFrom.AddDate(0, 0, 7)
	case Monthly:
		return startFrom.AddDate(0, 1, 0)
	default:
		return values.NoTime
	}
}

func (s *SubscriptionServiceImpl) chargeSubscription(ctx context.Context, subscription Subscription) (Subscription, error) {
	now := time.Now()
	subscription.UpdatedAt = now
	subscription.DueAt = s.computeDueAt(subscription, now)

	transaction := transaction.Transaction{
		ID:          uuid.New(),
		Description: subscription.GetTransactionDescription(),
		Amount:      subscription.Fee,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.transactionManager.Execute(ctx, func(ctx context.Context) error {
		if err := s.subscriptionRepository.Save(ctx, subscription); err != nil {
			return err
		}

		if err := s.transactionRepository.Save(ctx, transaction); err != nil {
			return err
		}

		return nil
	}); err != nil {
		fmt.Println(err.Error())
		return NoSubscription, errors.ErrInternalServer
	}

	return subscription, nil
}

func NewSubscriptionService(subscriptionRepository SubscriptionRepository, transactionRepository transaction.TransactionRepository, transactionManager manager.TransactionManager) SubscriptionService {
	return &SubscriptionServiceImpl{
		subscriptionRepository: subscriptionRepository,
		transactionRepository:  transactionRepository,
		transactionManager:     transactionManager,
	}
}
