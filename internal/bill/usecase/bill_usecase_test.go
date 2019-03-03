// +build unit

package usecase

import (
	"errors"
	"testing"

	"github.com/fairyhunter13/tax-calculator/internal/bill"
	mocksBill "github.com/fairyhunter13/tax-calculator/internal/bill/mocks"
	"github.com/fairyhunter13/tax-calculator/internal/taxobj"
	mocksTax "github.com/fairyhunter13/tax-calculator/internal/taxobj/mocks"
	"github.com/stretchr/testify/assert"
)

var (
	errDatabaseRepo = errors.New("Error in connecting to the database")
)

func TestBillUsecase_LoadData(t *testing.T) {
	t.Parallel()
	type fields struct {
		billRepo bill.Repository
		taxRepo  taxobj.Repository
	}
	tests := []struct {
		name    string
		fields  func() (bill.Repository, taxobj.Repository)
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Positive Case",
			fields: func() (bill.Repository, taxobj.Repository) {
				taxObject := taxobj.TaxObject{
					Name:    "MACD",
					TaxCode: 1,
					Price:   20000,
				}
				taxRepo := &mocksTax.Repository{}
				taxRepo.On("GetAll").Return([]taxobj.TaxObject{
					taxObject,
				}, nil)
				billRepo := &mocksBill.Repository{}
				billRepo.On("Add", taxObject)
				return billRepo, taxRepo
			},
			wantErr: false,
		},
		{
			name: "Tax Repo Database Error",
			fields: func() (bill.Repository, taxobj.Repository) {
				taxRepo := &mocksTax.Repository{}
				taxRepo.On("GetAll").Return([]taxobj.TaxObject{}, errDatabaseRepo)
				billRepo := &mocksBill.Repository{}
				return billRepo, taxRepo
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			billRepo, taxRepo := tt.fields()
			ucase := &BillUsecase{
				billRepo: billRepo,
				taxRepo:  taxRepo,
			}
			if err := ucase.LoadData(); (err != nil) != tt.wantErr {
				t.Errorf("BillUsecase.LoadData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBillUsecase_GetBill(t *testing.T) {
	t.Parallel()
	type fields struct {
		billRepo bill.Repository
		taxRepo  taxobj.Repository
	}
	tests := []struct {
		name   string
		fields func() (bill.Repository, taxobj.Repository)
		want   []bill.Bill
		want1  bill.Total
	}{
		// TODO: Add test cases.
		{
			name: "Empty Data",
			fields: func() (bill.Repository, taxobj.Repository) {
				taxRepo := &mocksTax.Repository{}
				billRepo := &mocksBill.Repository{}
				billRepo.On("GetAll").Return([]bill.Bill{}, bill.Total{})
				return billRepo, taxRepo
			},
			want:  []bill.Bill{},
			want1: bill.Total{},
		},
		{
			name: "A Data Stored in The Cache",
			fields: func() (bill.Repository, taxobj.Repository) {
				taxRepo := &mocksTax.Repository{}
				billRepo := &mocksBill.Repository{}
				aBill := []bill.Bill{
					bill.Bill{
						Name:       "MACD",
						TaxCode:    1,
						Price:      20000,
						Tax:        2000,
						Type:       "Food & Beverage",
						Refundable: "Yes",
						Amount:     22000,
					},
				}
				totalBill := bill.Total{
					PriceSubtotal: 20000,
					TaxSubtotal:   2000,
					GrandTotal:    22000,
				}
				billRepo.On("GetAll").Return(aBill, totalBill)
				return billRepo, taxRepo
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
			billRepo, taxRepo := tt.fields()
			ucase := &BillUsecase{
				billRepo: billRepo,
				taxRepo:  taxRepo,
			}
			got, got1 := ucase.GetBill()
			assert.EqualValues(t, got, tt.want)
			assert.EqualValues(t, got1, tt.want1)
		})
	}
}

func TestNewBillUsecase(t *testing.T) {
	type args struct {
		billRepo bill.Repository
		taxRepo  taxobj.Repository
	}
	billRepo := new(mocksBill.Repository)
	taxRepo := new(mocksTax.Repository)
	tests := []struct {
		name string
		args args
		want bill.Usecase
	}{
		// TODO: Add test cases.
		{
			name: "Init Bill Usecase",
			args: args{
				billRepo: billRepo,
				taxRepo:  taxRepo,
			},
			want: &BillUsecase{
				billRepo: billRepo,
				taxRepo:  taxRepo,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.EqualValues(t, NewBillUsecase(tt.args.billRepo, tt.args.taxRepo), tt.want)
		})
	}
}
