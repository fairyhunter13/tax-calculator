package taxobj

//Usecase defines the required behavior for business logic in the tax object.
type Usecase interface {
	CreateTaxObject(*TaxObject) error
}
