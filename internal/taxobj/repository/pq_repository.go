package repository

import (
	"database/sql"

	"github.com/fairyhunter13/tax-calculator/internal/taxobj"
)

//PqRepository is the repository for managing the data using postgre.
type PqRepository struct {
	pool      *sql.DB
	statement statement
}

type statement struct {
	insert    *sql.Stmt
	selectAll *sql.Stmt
}

const (
	queryInsert = `
		INSERT INTO tax_object
			(id, name, tax_code, price)
		VALUES
			(DEFAULT, $1, $2, $3)
		RETURNING id
	`
	querySelectAll = `
		SELECT 
			*
		FROM
			tax_object
	`
)

//NewPqRepository creates the pq repository for tax object with postgre connection.
func NewPqRepository(pool *sql.DB) (taxobj.Repository, error) {
	pqRepository := &PqRepository{
		pool: pool,
	}
	stmt, err := pool.Prepare(queryInsert)
	if err != nil {
		return nil, err
	}
	pqRepository.statement.insert = stmt
	stmt, err = pool.Prepare(querySelectAll)
	if err != nil {
		return nil, err
	}
	pqRepository.statement.selectAll = stmt

	return pqRepository, nil
}

//GetAll return all tax objects in postgre.
func (repo *PqRepository) GetAll() (taxObjects []taxobj.TaxObject, err error) {
	var taxObject taxobj.TaxObject
	taxObjects = make([]taxobj.TaxObject, 0)

	rows, err := repo.statement.selectAll.Query()
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		taxObject = taxobj.TaxObject{}
		err = rows.Scan(
			&taxObject.ID,
			&taxObject.Name,
			&taxObject.TaxCode,
			&taxObject.Price,
		)
		if err != nil {
			return
		}
		taxObjects = append(taxObjects, taxObject)
	}

	err = rows.Err()

	return
}

//Create create a new tax object in the database.
func (repo *PqRepository) Create(taxObj *taxobj.TaxObject) (err error) {
	row := repo.statement.insert.QueryRow(taxObj.Name, taxObj.TaxCode, taxObj.Price)

	err = row.Scan(
		&taxObj.ID,
	)

	return
}

//Close close all prepared statements in this repository.
func (repo *PqRepository) Close() {
	repo.statement.insert.Close()
	repo.statement.selectAll.Close()
}
