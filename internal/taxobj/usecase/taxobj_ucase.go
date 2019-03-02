package usecase

import (
	"github.com/fairyhunter13/tax-calculator/internal/taxobj"
)

//TaxObjectUsecase defines all the business logic for the tax object.
type TaxObjectUsecase struct {
	taxObjRepo taxobj.Repository
}

//NewTaxObjectUsecase return the tax object usecase.
func NewTaxObjectUsecase(taxObjRepo taxobj.Repository) taxobj.Usecase {
	return &TaxObjectUsecase{
		taxObjRepo,
	}
}

//CreateTaxObject create a new tax object and store it into the database.
func (ucase *TaxObjectUsecase) CreateTaxObject(taxObject *taxobj.TaxObject) (err error) {
	err = ucase.taxObjRepo.Create(taxObject)
	return
}
