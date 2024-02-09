package subscription_specification

import (
	"strings"
	"time"

	"github.com/google/uuid"

	subscription_entity "github.com/fikrirnurhidayat/banda-lumaksa/internal/domain/subscription/entity"
	subscription_types "github.com/fikrirnurhidayat/banda-lumaksa/internal/domain/subscription/types"
)

type SubscriptionSpecification interface {
	Call(subscription subscription_entity.Subscription) bool
}

type SubscriptionSpecifications []SubscriptionSpecification

type NameLikeSpecification struct {
	Substring string
}

func (spec NameLikeSpecification) Call(subscription subscription_entity.Subscription) bool {
	return strings.Contains(subscription.Name, spec.Substring)
}

func NameLike(value string) SubscriptionSpecification {
	return NameLikeSpecification{
		Substring: value,
	}
}

type NameIsSpecification struct {
	Name string
}

func (spec NameIsSpecification) Call(subscription subscription_entity.Subscription) bool {
	return subscription.Name == spec.Name
}

func NameIs(value string) SubscriptionSpecification {
	return NameIsSpecification{
		Name: value,
	}
}

type TypeIsSpecification struct {
	Type subscription_types.Type
}

func (spec TypeIsSpecification) Call(subscription subscription_entity.Subscription) bool {
	return subscription.Type == spec.Type
}

func TypeIs(value subscription_types.Type) SubscriptionSpecification {
	return TypeIsSpecification{
		Type: value,
	}
}

type DueBetweenSpecification struct {
	Start time.Time
	End   time.Time
}

func (spec DueBetweenSpecification) Call(subscription subscription_entity.Subscription) bool {
	return subscription.DueAt.After(spec.Start) && subscription.DueAt.Before(spec.End)
}

func DueBetween(start time.Time, end time.Time) SubscriptionSpecification {
	return DueBetweenSpecification{
		Start: start,
		End:   end,
	}
}

func DueIn(now time.Time) SubscriptionSpecification {
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
	return DueBetween(startOfDay, endOfDay)
}

type DueBeforeSpecification struct {
	Now time.Time
}

func (spec DueBeforeSpecification) Call(subscription subscription_entity.Subscription) bool {
	return subscription.DueAt.Before(spec.Now)
}

func DueBefore(now time.Time) SubscriptionSpecification {
	return DueBeforeSpecification{
		Now: now,
	}
}

type CreatedBetweenSpecification struct {
	Start time.Time
	End   time.Time
}

func (spec CreatedBetweenSpecification) Call(subscription subscription_entity.Subscription) bool {
	return subscription.CreatedAt.After(spec.Start) && subscription.CreatedAt.Before(spec.End)
}

func CreatedBetween(start time.Time, end time.Time) SubscriptionSpecification {
	return CreatedBetweenSpecification{
		Start: start,
		End:   end,
	}
}

type UpdatedBetweenSpecification struct {
	Start time.Time
	End   time.Time
}

func (spec UpdatedBetweenSpecification) Call(subscription subscription_entity.Subscription) bool {
	return subscription.UpdatedAt.After(spec.Start) && subscription.UpdatedAt.Before(spec.End)
}

func UpdatedBetween(start time.Time, end time.Time) SubscriptionSpecification {
	return UpdatedBetweenSpecification{
		Start: start,
		End:   end,
	}
}

type StartedBetweenSpecification struct {
	Start time.Time
	End   time.Time
}

func (spec StartedBetweenSpecification) Call(subscription subscription_entity.Subscription) bool {
	return subscription.StartedAt.After(spec.Start) && subscription.StartedAt.Before(spec.End)
}

func StartedBetween(start time.Time, end time.Time) SubscriptionSpecification {
	return StartedBetweenSpecification{
		Start: start,
		End:   end,
	}
}

type EndedBetweenSpecification struct {
	Start time.Time
	End   time.Time
}

func (spec EndedBetweenSpecification) Call(subscription subscription_entity.Subscription) bool {
	return subscription.EndedAt.After(spec.Start) && subscription.EndedAt.Before(spec.End)
}

func EndedBetween(start time.Time, end time.Time) SubscriptionSpecification {
	return EndedBetweenSpecification{
		Start: start,
		End:   end,
	}
}

type NotEndedSpecification struct {
	Now time.Time
}

func (spec NotEndedSpecification) Call(subscription subscription_entity.Subscription) bool {
	return subscription.EndedAt.Before(spec.Now)
}

func NotEnded(now time.Time) SubscriptionSpecification {
	return NotEndedSpecification{
		Now: now,
	}
}

type WithIDSpecification struct {
	ID uuid.UUID
}

func (spec WithIDSpecification) Call(subscription subscription_entity.Subscription) bool {
	return spec.ID == subscription.ID
}

func WithID(id uuid.UUID) SubscriptionSpecification {
	return WithIDSpecification{
		ID: id,
	}
}
