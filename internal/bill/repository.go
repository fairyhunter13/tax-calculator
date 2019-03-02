package bill

import (
	"github.com/fairyhunter13/tax-calculator/internal/taxobj"
)

//Repository define the required behavior of data management in the bill.
type Repository interface {
	Add(taxobj.TaxObject)
	GetAll() ([]Bill, Total)
}
