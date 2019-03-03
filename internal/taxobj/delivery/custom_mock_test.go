// +build unit

package delivery

import (
	taxobj "github.com/fairyhunter13/tax-calculator/internal/taxobj"
)

// usecase is an autogenerated mock type for the Usecase type
type usecase struct {
}

// CreateTaxObject provides a mock function with given fields: _a0
func (ucase *usecase) CreateTaxObject(taxObj *taxobj.TaxObject) error {
	taxObj.ID = 1
	return nil
}
