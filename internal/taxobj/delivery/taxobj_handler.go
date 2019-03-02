package delivery

import (
	"net/http"
	"sync"

	"github.com/fairyhunter13/tax-calculator/internal/customhttp"

	"github.com/fairyhunter13/tax-calculator/internal/taxobj"
	"github.com/labstack/echo"
	"github.com/microcosm-cc/bluemonday"
	validator "gopkg.in/go-playground/validator.v9"
)

var (
	httpHandler      *HTTPTaxObjectHandler
	once             sync.Once
	requestValidator *validator.Validate
	sanitizer        *bluemonday.Policy
)

func init() {
	//Init once sanitizer and request validator.
	once.Do(func() {
		requestValidator = validator.New()
		sanitizer = bluemonday.UGCPolicy()
	})
}

//HTTPTaxObjectHandler defines the http delivery layer for the tax object.
type HTTPTaxObjectHandler struct {
	taxObjUcase taxobj.Usecase
}

//NewTaxObjectHandler create the HTTPTaxObjectHandler with customed routing for echo.
func NewTaxObjectHandler(e *echo.Echo, taxObjUcase taxobj.Usecase) {
	httpHandler = &HTTPTaxObjectHandler{
		taxObjUcase,
	}
	e.POST("/tax", httpHandler.CreateTaxObject)
}

//CreateTaxObject handle request for creating the tax object.
func (handler *HTTPTaxObjectHandler) CreateTaxObject(c echo.Context) (err error) {
	taxObject := taxobj.TaxObject{}
	if err = c.Bind(&taxObject); err != nil {
		err = customhttp.ErrInvalidInput
		return
	}
	if err = requestValidator.Struct(&taxObject); err != nil {
		err = customhttp.ErrInvalidInput
		return
	}
	taxObject.Name = sanitizer.Sanitize(taxObject.Name)
	if err = handler.taxObjUcase.CreateTaxObject(&taxObject); err != nil {
		return
	}
	err = c.JSON(http.StatusCreated, &taxObject)
	return
}
