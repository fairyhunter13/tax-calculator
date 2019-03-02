package repository

import (
	"sync"

	"github.com/fairyhunter13/tax-calculator/internal/taxobj"

	"github.com/fairyhunter13/tax-calculator/internal/bill"
)

const (
	//Refundable defines the text to show if it's refundable.
	Refundable = "Yes"
	//NotRefundable defines the text to show if it's not refundable.
	NotRefundable = "No"
)

//CacheRepository defines the data management for the bill.
type CacheRepository struct {
	mutex *sync.Mutex
	//mutex here protected the following fileds.
	bills []bill.Bill
	total bill.Total
}

var (
	typeMap = map[int64]string{
		1: "Food & Beverage",
		2: "Tobacco",
		3: "Entertainment",
	}
)

//NewCacheRepository return the concrete implementation of repository using cache.
func NewCacheRepository() bill.Repository {
	cacheRepo := &CacheRepository{
		mutex: new(sync.Mutex),
		total: bill.Total{},
		bills: make([]bill.Bill, 0),
	}
	return cacheRepo
}

//Add add tax object to the bill list.
func (repo *CacheRepository) Add(taxObject taxobj.TaxObject) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	//Adding bill to bill cache
	billObject := bill.Bill{
		Name:       taxObject.Name,
		Price:      taxObject.Price,
		TaxCode:    taxObject.TaxCode,
		Refundable: repo.getRefundable(taxObject.TaxCode),
		Type:       repo.getType(taxObject.TaxCode),
		Tax:        repo.getTax(taxObject.TaxCode, taxObject.Price),
	}
	billObject.Amount = billObject.Tax + billObject.Price
	repo.bills = append(repo.bills, billObject)
	//Calculating total cache
	repo.total.PriceSubtotal += taxObject.Price
	repo.total.TaxSubtotal += billObject.Tax
	repo.total.GrandTotal += billObject.Amount
}

//GetAll return the bill list.
func (repo *CacheRepository) GetAll() ([]bill.Bill, bill.Total) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	return repo.bills, repo.total
}

//getRefundable return the refundable text to display based on the tax code.
func (repo *CacheRepository) getRefundable(taxCode int64) (refundable string) {
	switch taxCode {
	case 1:
		refundable = Refundable
	case 2:
		refundable = NotRefundable
	case 3:
		refundable = NotRefundable
	}
	return
}

//getType return the type text for the given tax code.
func (repo *CacheRepository) getType(taxCode int64) string {
	return typeMap[taxCode]
}

//getTax return the calculated tax for the given tax code and price.
func (repo *CacheRepository) getTax(taxCode int64, price float64) (tax float64) {
	if price <= 0 {
		return
	}
	switch taxCode {
	case 1:
		tax = float64(10) / float64(100) * price
	case 2:
		tax = float64(10) + (float64(2) / float64(100) * price)
	case 3:
		if price < 100 {
			return
		}
		tax = float64(1) / float64(100) * (price - float64(100))
	}
	return
}
