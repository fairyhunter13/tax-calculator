package taxobj

//Repository define the required behavior of data management in the tax object.
type Repository interface {
	GetAll() ([]TaxObject, error)
	Create(*TaxObject) error
	Close()
}
