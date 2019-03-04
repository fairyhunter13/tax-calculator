// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import bill "github.com/fairyhunter13/tax-calculator/internal/bill"
import mock "github.com/stretchr/testify/mock"
import taxobj "github.com/fairyhunter13/tax-calculator/internal/taxobj"

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// Add provides a mock function with given fields: _a0
func (_m *Repository) Add(_a0 taxobj.TaxObject) {
	_m.Called(_a0)
}

// GetAll provides a mock function with given fields:
func (_m *Repository) GetAll() ([]bill.Bill, bill.Total) {
	ret := _m.Called()

	var r0 []bill.Bill
	if rf, ok := ret.Get(0).(func() []bill.Bill); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]bill.Bill)
		}
	}

	var r1 bill.Total
	if rf, ok := ret.Get(1).(func() bill.Total); ok {
		r1 = rf()
	} else {
		r1 = ret.Get(1).(bill.Total)
	}

	return r0, r1
}