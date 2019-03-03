// +build unit

package delivery

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fairyhunter13/tax-calculator/internal/taxobj"
	"github.com/fairyhunter13/tax-calculator/internal/taxobj/mocks"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

const (
	validJSON = `
		{
			"name": "MACD",
			"tax_code": 1,
			"price": 20000
		}
	`
	invalidJSON = `
		{
			"hello": "hai",
		}
	`
	invalidInput = `
		{
			"name": "MACD",
			"tax_code": 0,
			"price": -20000
		}
	`
)

var (
	errDatabaseNotOnline = errors.New("Database is not online")
)

type errorResp struct {
	Message string `json:"message"`
}

func TestHTTPTaxObjectHandler_CreateTaxObject_Positive(t *testing.T) {
	t.Parallel()
	const logFail = `[TestHTTPTaxObjectHandler_CreateTaxObject_Positive] %s: %s`
	expectedResp := taxobj.TaxObject{
		ID:      1,
		Name:    "MACD",
		TaxCode: 1,
		Price:   20000,
	}
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/tax", strings.NewReader(validJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	taxUcase := &usecase{}
	h := &HTTPTaxObjectHandler{
		taxObjUcase: taxUcase,
	}

	// Assertions using testify framework
	if assert.NoError(t, h.CreateTaxObject(ctx)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		actualResp := taxobj.TaxObject{}
		err := json.Unmarshal(rec.Body.Bytes(), &actualResp)
		if err != nil {
			t.Errorf(logFail, "Error in unmarshaling json", err)
		}
		assert.Equal(t, expectedResp, actualResp)
	}
}

func TestHTTPTaxObjectHandler_CreateTaxObject_InvalidJSONSyntax(t *testing.T) {
	t.Parallel()
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/tax", strings.NewReader(invalidJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	taxUcase := &mocks.Usecase{}
	taxUcase.On("CreateTaxObject", new(taxobj.TaxObject)).Return(nil)
	h := &HTTPTaxObjectHandler{
		taxObjUcase: taxUcase,
	}

	// Assertions using testify framework
	err := h.CreateTaxObject(ctx)
	if assert.Error(t, err) {
		//Assertion using the string of error.
		//If we use the recorder, we don't get valid output.
		assert.Equal(t, fmt.Sprintf("%s", ErrInvalidInput), fmt.Sprintf("%s", err))
	}
}

func TestHTTPTaxObjectHandler_CreateTaxObject_InvalidInput(t *testing.T) {
	t.Parallel()
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/tax", strings.NewReader(invalidInput))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	taxUcase := &mocks.Usecase{}
	taxUcase.On("CreateTaxObject", new(taxobj.TaxObject)).Return(nil)
	h := &HTTPTaxObjectHandler{
		taxObjUcase: taxUcase,
	}

	// Assertions using testify framework
	err := h.CreateTaxObject(ctx)
	if assert.Error(t, err) {
		//Assertion using the string of error.
		//If we use the recorder, we don't get valid output.
		assert.Equal(t, fmt.Sprintf("%s", ErrInvalidInput), fmt.Sprintf("%s", err))
	}
}

func TestHTTPTaxObjectHandler_CreateTaxObject_InternalServerError(t *testing.T) {
	t.Parallel()
	e := echo.New()
	arg := &taxobj.TaxObject{
		Name:    "MACD",
		TaxCode: 1,
		Price:   20000,
	}
	req := httptest.NewRequest(http.MethodPost, "/tax", strings.NewReader(validJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	taxUcase := &mocks.Usecase{}
	taxUcase.On("CreateTaxObject", arg).Return(errDatabaseNotOnline)
	h := &HTTPTaxObjectHandler{
		taxObjUcase: taxUcase,
	}

	// Assertions using testify framework
	err := h.CreateTaxObject(ctx)
	if assert.Error(t, err) {
		//Assertion using the string of error.
		//If we use the recorder, we don't get valid output.
		assert.Equal(t, errDatabaseNotOnline, err)
	}
}

func TestNewTaxObjectHandler(t *testing.T) {
	t.Parallel()
	type args struct {
		e           *echo.Echo
		taxObjUcase taxobj.Usecase
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "Init Tax Object HTTPHandler",
			args: args{
				e:           echo.New(),
				taxObjUcase: &mocks.Usecase{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NewTaxObjectHandler(tt.args.e, tt.args.taxObjUcase)
		})
	}
}
