// +build unit

package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fairyhunter13/tax-calculator/internal/bill"
	"github.com/fairyhunter13/tax-calculator/internal/bill/mocks"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestHTTPBillHandler_GetBill_EmptyData(t *testing.T) {
	t.Parallel()
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/bill", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	billUcase := &mocks.Usecase{}
	actualResponse := BillResponse{
		Bill:  []bill.Bill{},
		Total: bill.Total{},
	}
	billUcase.On("GetBill").Return(actualResponse.Bill, actualResponse.Total)
	h := &HTTPBillHandler{
		billUcase: billUcase,
	}

	// Assertions using testify framework
	err := h.GetBill(ctx)
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		billResp := BillResponse{}
		err = json.Unmarshal(rec.Body.Bytes(), &billResp)
		if err != nil {
			assert.Errorf(t, err, "Error unmarshaling bill response: %s", err)
		}
		assert.EqualValues(t, billResp, actualResponse)
	}
}

func TestHTTPBillHandler_GetBill_AvailableData(t *testing.T) {
	t.Parallel()
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/bill", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	billUcase := &mocks.Usecase{}
	actualResponse := BillResponse{
		Bill: []bill.Bill{
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
		Total: bill.Total{
			PriceSubtotal: 20000,
			TaxSubtotal:   2000,
			GrandTotal:    22000,
		},
	}
	billUcase.On("GetBill").Return(actualResponse.Bill, actualResponse.Total)
	h := &HTTPBillHandler{
		billUcase: billUcase,
	}

	// Assertions using testify framework
	err := h.GetBill(ctx)
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		billResp := BillResponse{}
		err = json.Unmarshal(rec.Body.Bytes(), &billResp)
		if err != nil {
			assert.Errorf(t, err, "Error unmarshaling bill response: %s", err)
		}
		assert.EqualValues(t, billResp, actualResponse)
	}
}

func TestNewHTTPBillHandler(t *testing.T) {
	t.Parallel()
	type args struct {
		e         *echo.Echo
		billUcase bill.Usecase
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "Init Bill HTTP Handler",
			args: args{
				e:         echo.New(),
				billUcase: new(mocks.Usecase),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NewHTTPBillHandler(tt.args.e, tt.args.billUcase)
		})
	}
}
