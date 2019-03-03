package usecase

import (
	"github.com/fairyhunter13/tax-calculator/internal/bill"
	"github.com/fairyhunter13/tax-calculator/internal/taxobj"
)

//BillUsecase define the business logic for bill.
type BillUsecase struct {
	billRepo bill.Repository
	taxRepo  taxobj.Repository
}

//NewBillUsecase creates the new BillUsacase concrete implementation.
func NewBillUsecase(billRepo bill.Repository, taxRepo taxobj.Repository) bill.Usecase {
	return &BillUsecase{
		billRepo,
		taxRepo,
	}
}

//LoadData init the cache for the first time start of application.
//It is useful to store already calculated value in the fly,
//so the performance is good when fetching all tax object in bill.
func (ucase *BillUsecase) LoadData() (err error) {
	taxObjects, err := ucase.taxRepo.GetAll()
	if err != nil {
		return
	}
	for _, taxObject := range taxObjects {
		ucase.billRepo.Add(taxObject)
	}
	return
}

//GetBill get the bill and total data.
func (ucase *BillUsecase) GetBill() ([]bill.Bill, bill.Total) {
	return ucase.billRepo.GetAll()
}
