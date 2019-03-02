package bill

//Usecase defines the required behavior for business logic in the bill.
type Usecase interface {
	LoadData() error
	GetBill() ([]Bill, Total)
}
