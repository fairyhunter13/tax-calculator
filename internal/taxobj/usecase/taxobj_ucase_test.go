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

func TestTaxObjectUsecase_CreateTaxObject(t *testing.T) {
	t.Parallel()
	type args struct {
		taxObject *taxobj.TaxObject
	}
	tests := []struct {
		name    string
		ucase   func() *TaxObjectUsecase
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Positive Case",
			ucase: func() *TaxObjectUsecase {
				taxObj := &taxobj.TaxObject{
					Name:    "MACD",
					TaxCode: 1,
					Price:   20000,
				}
				taxRepo := &mocksTax.Repository{}
				taxRepo.On("Create", taxObj).Return(nil)
				billRepo := &mocksBill.Repository{}
				billRepo.On("Add", *taxObj).Return()
				ucase := NewTaxObjectUsecase(taxRepo, billRepo)
				return ucase.(*TaxObjectUsecase)
			},
			args: args{
				taxObject: &taxobj.TaxObject{
					Name:    "MACD",
					TaxCode: 1,
					Price:   20000,
				},
			},
			wantErr: false,
		},
		{
			name: "Error in storing to the database for the tax object",
			ucase: func() *TaxObjectUsecase {
				taxObj := &taxobj.TaxObject{
					Name:    "MACD",
					TaxCode: 1,
					Price:   20000,
				}
				taxRepo := &mocksTax.Repository{}
				taxRepo.On("Create", taxObj).Return(errors.New("Error in storing to the database"))
				billRepo := &mocksBill.Repository{}
				ucase := NewTaxObjectUsecase(taxRepo, billRepo)
				return ucase.(*TaxObjectUsecase)
			},
			args: args{
				taxObject: &taxobj.TaxObject{
					Name:    "MACD",
					TaxCode: 1,
					Price:   20000,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.ucase().CreateTaxObject(tt.args.taxObject); (err != nil) != tt.wantErr {
				t.Errorf("TaxObjectUsecase.CreateTaxObject() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewTaxObjectUsecase(t *testing.T) {
	type args struct {
		taxObjRepo taxobj.Repository
		billRepo   bill.Repository
	}
	taxObjRepo := &mocksTax.Repository{}
	billRepo := &mocksBill.Repository{}
	tests := []struct {
		name string
		args args
		want taxobj.Usecase
	}{
		// TODO: Add test cases.
		{
			name: "Init Tax Object Usecase",
			args: args{
				taxObjRepo: taxObjRepo,
				billRepo:   billRepo,
			},
			want: &TaxObjectUsecase{
				taxObjRepo: taxObjRepo,
				billRepo:   billRepo,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.EqualValues(t, NewTaxObjectUsecase(taxObjRepo, billRepo), tt.want)
		})
	}
}
