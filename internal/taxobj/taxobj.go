package taxobj

//TaxObject define the model for tax object.
//This is the data that the user will input.
//Tax objects are also used to calculate bills.
type TaxObject struct {
	ID      int64   `json:"id"`
	Name    string  `json:"name" validate:"required"`
	TaxCode int64   `json:"tax_code" validate:"required,gte=1,lte=3"`
	Price   float64 `json:"price" validate:"required"`
}
