package handler

import (
	"net/http"

	"github.com/fairyhunter13/tax-calculator/internal/bill"
	"github.com/labstack/echo"
)

//HTTPBillHandler define the http delivery layer for the bill.
type HTTPBillHandler struct {
	billUcase bill.Usecase
}

//BillResponse define the default json response for the bill.
type BillResponse struct {
	Bill  []bill.Bill `json:"bill"`
	Total bill.Total  `json:"total"`
}

var (
	httpHandler *HTTPBillHandler
)

//NewHTTPBillHandler define the routing for HTTPBillHandler.
func NewHTTPBillHandler(e *echo.Echo, billUcase bill.Usecase) {
	httpHandler = &HTTPBillHandler{
		billUcase,
	}
	e.GET("/bill", httpHandler.GetBill)
}

//GetBill get the bill list that has been calculated.
func (handler *HTTPBillHandler) GetBill(c echo.Context) (err error) {
	bills, total := handler.billUcase.GetBill()
	billResp := &BillResponse{
		Bill:  bills,
		Total: total,
	}
	c.JSON(http.StatusOK, billResp)
	return
}
