// +build unit

package repository

import (
	"reflect"
	"sync"
	"testing"

	"github.com/fairyhunter13/tax-calculator/internal/bill"
	"github.com/fairyhunter13/tax-calculator/internal/taxobj"
	"github.com/stretchr/testify/assert"
)

func TestCacheRepository_Add(t *testing.T) {
	t.Parallel()
	type fields struct {
		mutex *sync.Mutex
		bills []bill.Bill
		total bill.Total
	}
	type args struct {
		taxObject taxobj.TaxObject
	}
	type expectedState struct {
		bills []bill.Bill
		total bill.Total
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		expectedState expectedState
	}{
		// TODO: Add test cases.
		{
			name: "Positive Case",
			fields: fields{
				mutex: new(sync.Mutex),
				bills: []bill.Bill{},
				total: bill.Total{},
			},
			args: args{
				taxObject: taxobj.TaxObject{
					Name:    "MACD",
					TaxCode: 1,
					Price:   20000,
				},
			},
			expectedState: expectedState{
				bills: []bill.Bill{
					bill.Bill{
						Name:       "MACD",
						TaxCode:    1,
						Price:      20000,
						Tax:        2000,
						Type:       "Food & Beverage",
						Refundable: "Yes",
						Amount:     22000,
					},
				},
				total: bill.Total{
					PriceSubtotal: 20000,
					TaxSubtotal:   2000,
					GrandTotal:    22000,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &CacheRepository{
				mutex: tt.fields.mutex,
				bills: tt.fields.bills,
				total: tt.fields.total,
			}
			repo.Add(tt.args.taxObject)
			assert.Equal(t, repo.bills, tt.expectedState.bills)
			assert.Equal(t, repo.total, tt.expectedState.total)
		})
	}
}

func TestCacheRepository_GetAll(t *testing.T) {
	t.Parallel()
	type fields struct {
		mutex *sync.Mutex
		bills []bill.Bill
		total bill.Total
	}
	tests := []struct {
		name   string
		fields fields
		want   []bill.Bill
		want1  bill.Total
	}{
		// TODO: Add test cases.
		{
			name: "Positive Case Empty Data",
			fields: fields{
				mutex: new(sync.Mutex),
				bills: []bill.Bill{},
				total: bill.Total{},
			},
			want:  []bill.Bill{},
			want1: bill.Total{},
		},
		{
			name: "Positive Case A Data",
			fields: fields{
				mutex: new(sync.Mutex),
				bills: []bill.Bill{
					bill.Bill{
						Name:       "MACD",
						TaxCode:    1,
						Price:      20000,
						Tax:        2000,
						Type:       "Food & Beverage",
						Refundable: "Yes",
						Amount:     22000,
					},
				},
				total: bill.Total{
					PriceSubtotal: 20000,
					TaxSubtotal:   2000,
					GrandTotal:    22000,
				},
			},
			want: []bill.Bill{
				bill.Bill{
					Name:       "MACD",
					TaxCode:    1,
					Price:      20000,
					Tax:        2000,
					Type:       "Food & Beverage",
					Refundable: "Yes",
					Amount:     22000,
				},
			},
			want1: bill.Total{
				PriceSubtotal: 20000,
				TaxSubtotal:   2000,
				GrandTotal:    22000,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &CacheRepository{
				mutex: tt.fields.mutex,
				bills: tt.fields.bills,
				total: tt.fields.total,
			}
			got, got1 := repo.GetAll()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CacheRepository.GetAll() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("CacheRepository.GetAll() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestCacheRepository_getRefundable(t *testing.T) {
	t.Parallel()
	type fields struct {
		mutex *sync.Mutex
		bills []bill.Bill
		total bill.Total
	}
	type args struct {
		taxCode int64
	}
	defaultFields := fields{
		mutex: new(sync.Mutex),
		bills: []bill.Bill{},
		total: bill.Total{},
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantRefundable string
	}{
		// TODO: Add test cases.
		{
			name:   "Food & Beverage",
			fields: defaultFields,
			args: args{
				taxCode: 1,
			},
			wantRefundable: "Yes",
		},
		{
			name:   "Tobacco",
			fields: defaultFields,
			args: args{
				taxCode: 2,
			},
			wantRefundable: "No",
		},
		{
			name:   "Entertainment",
			fields: defaultFields,
			args: args{
				taxCode: 3,
			},
			wantRefundable: "No",
		},
		{
			name:   "Invalid Tax Code",
			fields: defaultFields,
			args: args{
				taxCode: 0,
			},
			wantRefundable: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &CacheRepository{
				mutex: tt.fields.mutex,
				bills: tt.fields.bills,
				total: tt.fields.total,
			}
			if gotRefundable := repo.getRefundable(tt.args.taxCode); gotRefundable != tt.wantRefundable {
				t.Errorf("CacheRepository.getRefundable() = %v, want %v", gotRefundable, tt.wantRefundable)
			}
		})
	}
}

func TestCacheRepository_getType(t *testing.T) {
	t.Parallel()
	type fields struct {
		mutex *sync.Mutex
		bills []bill.Bill
		total bill.Total
	}
	type args struct {
		taxCode int64
	}
	defaultFields := fields{
		mutex: new(sync.Mutex),
		bills: []bill.Bill{},
		total: bill.Total{},
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		// TODO: Add test cases.
		{
			name:   "Food & Beverage",
			fields: defaultFields,
			args: args{
				taxCode: 1,
			},
			want: "Food & Beverage",
		},
		{
			name:   "Tobacco",
			fields: defaultFields,
			args: args{
				taxCode: 2,
			},
			want: "Tobacco",
		},
		{
			name:   "Entertainment",
			fields: defaultFields,
			args: args{
				taxCode: 3,
			},
			want: "Entertainment",
		},
		{
			name:   "Invalid Tax Code",
			fields: defaultFields,
			args: args{
				taxCode: 0,
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &CacheRepository{
				mutex: tt.fields.mutex,
				bills: tt.fields.bills,
				total: tt.fields.total,
			}
			if got := repo.getType(tt.args.taxCode); got != tt.want {
				t.Errorf("CacheRepository.getType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCacheRepository_getTax(t *testing.T) {
	t.Parallel()
	type fields struct {
		mutex *sync.Mutex
		bills []bill.Bill
		total bill.Total
	}
	type args struct {
		taxCode int64
		price   float64
	}
	defaultFields := fields{
		mutex: new(sync.Mutex),
		bills: []bill.Bill{},
		total: bill.Total{},
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantTax float64
	}{
		// TODO: Add test cases.
		{
			name:   "Food & Beverage",
			fields: defaultFields,
			args: args{
				taxCode: 1,
				price:   10000,
			},
			wantTax: 1000,
		},
		{
			name:   "Tobacco",
			fields: defaultFields,
			args: args{
				taxCode: 2,
				price:   1000,
			},
			wantTax: 30,
		},
		{
			name:   "Entertainment Above 100",
			fields: defaultFields,
			args: args{
				taxCode: 3,
				price:   120,
			},
			wantTax: 0.2,
		},
		{
			name:   "Entertainment Below 100",
			fields: defaultFields,
			args: args{
				taxCode: 3,
				price:   50,
			},
			wantTax: 0,
		},
		{
			name:   "Invalid Tax Code",
			fields: defaultFields,
			args: args{
				taxCode: 0,
			},
			wantTax: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &CacheRepository{
				mutex: tt.fields.mutex,
				bills: tt.fields.bills,
				total: tt.fields.total,
			}
			if gotTax := repo.getTax(tt.args.taxCode, tt.args.price); gotTax != tt.wantTax {
				t.Errorf("CacheRepository.getTax() = %v, want %v", gotTax, tt.wantTax)
			}
		})
	}
}

func TestNewCacheRepository(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		want bill.Repository
	}{
		// TODO: Add test cases.
		{
			name: "Init Bill Cache Repository",
			want: &CacheRepository{
				mutex: new(sync.Mutex),
				bills: []bill.Bill{},
				total: bill.Total{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewCacheRepository()
			assert.EqualValues(t, got, tt.want)
		})
	}
}
