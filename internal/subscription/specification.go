package subscription

import (
	"strings"
	"time"
)

type Specification interface {
	Call(subscription Subscription) bool
}

type NameLikeSpecification struct {
	Substring string
}

func (spec NameLikeSpecification) Call(subscription Subscription) bool {
	return strings.Contains(subscription.Name, spec.Substring)
}

func NameLike(value string) Specification {
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

func TypeIs(value Type) Specification {
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

func DueBetween(start time.Time, end time.Time) Specification {
	return &DueBetweenSpecification{
		Start: start,
		End:   end,
	}
}

type CreatedBetweenSpecification struct {
	Start time.Time
	End   time.Time
}

func (spec CreatedBetweenSpecification) Call(subscription Subscription) bool {
	return subscription.CreatedAt.After(spec.Start) && subscription.CreatedAt.Before(spec.End)
}

func CreatedBetween(start time.Time, end time.Time) Specification {
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

func StartedBetween(start time.Time, end time.Time) Specification {
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

func EndedBetween(start time.Time, end time.Time) Specification {
	return EndedBetweenSpecification{
		Start: start,
		End:   end,
	}
}

type LimitSpecification struct {
	Limit uint32
}

func (LimitSpecification) Call(subscription Subscription) bool {
	return true
}

func Limit(limit uint32) Specification {
	return LimitSpecification{
		Limit: limit,
	}
}

type OffsetSpecification struct {
	Offset uint32
}

func (OffsetSpecification) Call(subscription Subscription) bool {
	return true
}

func Offset(offset uint32) Specification {
	return OffsetSpecification{
		Offset: offset,
	}
}

type SortSpecification struct {
	Args []SortArg
}

func (SortSpecification) Call(subscription Subscription) bool {
	return true
}

type SortArg struct {
	Key       string
	Direction string
}

func Sort(args ...SortArg) Specification {
	return SortSpecification{
		Args: args,
	}
}
