package bill

//Bill define the data model for bill.
//Bill list all the calculated data from the tax objects.
//This data that will be seen by user.
type Bill struct {
	Name       string  `json:"name"`
	TaxCode    int64   `json:"tax_code"`
	Type       string  `json:"type"`
	Refundable string  `json:"refundable"`
	Price      float64 `json:"price"`
	Tax        float64 `json:"tax"`
	Amount     float64 `json:"amount"`
}

//Total define the total calculation for each price, tax, and amount.
type Total struct {
	PriceSubtotal float64 `json:"price_subtotal"`
	TaxSubtotal   float64 `json:"tax_subtotal"`
	GrandTotal    float64 `json:"grand_total"`
}
