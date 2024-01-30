package subscription

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type SubscriptionSpecification interface {
	Call(subscription Subscription) bool
}

type SubscriptionSpecifications []SubscriptionSpecification

type NameLikeSpecification struct {
	Substring string
}

func (spec NameLikeSpecification) Call(subscription Subscription) bool {
	return strings.Contains(subscription.Name, spec.Substring)
}

func NameLike(value string) SubscriptionSpecification {
	return &NameLikeSpecification{
		Substring: value,
	}
}

type TypeIsSpecification struct {
	Type Type
}

func (spec TypeIsSpecification) Call(subscription Subscription) bool {
	return subscription.Type == spec.Type
}

func TypeIs(value Type) SubscriptionSpecification {
	return &TypeIsSpecification{
		Type: value,
	}
}

type DueBetweenSpecification struct {
	Start time.Time
	End   time.Time
}

func (spec DueBetweenSpecification) Call(subscription Subscription) bool {
	return subscription.DueAt.After(spec.Start) && subscription.DueAt.Before(spec.End)
}

func DueBetween(start time.Time, end time.Time) SubscriptionSpecification {
	return &DueBetweenSpecification{
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

func (spec DueBeforeSpecification) Call(subscription Subscription) bool {
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

func (spec CreatedBetweenSpecification) Call(subscription Subscription) bool {
	return subscription.CreatedAt.After(spec.Start) && subscription.CreatedAt.Before(spec.End)
}

func CreatedBetween(start time.Time, end time.Time) SubscriptionSpecification {
	return CreatedBetweenSpecification{
		Start: start,
		End:   end,
	}
}

type StartedBetweenSpecification struct {
	Start time.Time
	End   time.Time
}

func (spec StartedBetweenSpecification) Call(subscription Subscription) bool {
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

func (spec EndedBetweenSpecification) Call(subscription Subscription) bool {
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

func (spec NotEndedSpecification) Call(subscription Subscription) bool {
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

func (spec WithIDSpecification) Call(subscription Subscription) bool {
	return spec.ID == subscription.ID
}

func WithID(id uuid.UUID) SubscriptionSpecification {
	return WithIDSpecification{
		ID: id,
	}
}
