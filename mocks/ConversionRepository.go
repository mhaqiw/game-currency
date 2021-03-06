// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/mhaqiw/game-currency/domain"
	mock "github.com/stretchr/testify/mock"

	time "time"
)

// ConversionRepository is an autogenerated mock type for the ConversionRepository type
type ConversionRepository struct {
	mock.Mock
}

// GetFromToID provides a mock function with given fields: ctx, fromID, toID
func (_m *ConversionRepository) GetFromToID(ctx context.Context, fromID int64, toID int64) (domain.Conversion, error) {
	ret := _m.Called(ctx, fromID, toID)

	var r0 domain.Conversion
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64) domain.Conversion); ok {
		r0 = rf(ctx, fromID, toID)
	} else {
		r0 = ret.Get(0).(domain.Conversion)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64, int64) error); ok {
		r1 = rf(ctx, fromID, toID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Post provides a mock function with given fields: ctx, conversion
func (_m *ConversionRepository) Post(ctx context.Context, conversion domain.Conversion) (int64, int64, time.Time, error) {
	ret := _m.Called(ctx, conversion)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, domain.Conversion) int64); ok {
		r0 = rf(ctx, conversion)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 int64
	if rf, ok := ret.Get(1).(func(context.Context, domain.Conversion) int64); ok {
		r1 = rf(ctx, conversion)
	} else {
		r1 = ret.Get(1).(int64)
	}

	var r2 time.Time
	if rf, ok := ret.Get(2).(func(context.Context, domain.Conversion) time.Time); ok {
		r2 = rf(ctx, conversion)
	} else {
		r2 = ret.Get(2).(time.Time)
	}

	var r3 error
	if rf, ok := ret.Get(3).(func(context.Context, domain.Conversion) error); ok {
		r3 = rf(ctx, conversion)
	} else {
		r3 = ret.Error(3)
	}

	return r0, r1, r2, r3
}
