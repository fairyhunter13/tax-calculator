package repository

import (
	"database/sql"
	"sync"

	"github.com/afex/hystrix-go/hystrix"

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
	querySelectOne = `
		SELECT
			id
		FROM
			tax_object
		LIMIT 1
	`
	queryCreateTable = `
		CREATE TABLE tax_object (
			id serial PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			tax_code bigint NOT NULL,
			price double precision NOT NULL
		)
	`
)

const (
	commandGetAll = `command_taxobj_GetAll`
	commandCreate = `command_taxobj_Create`
)

var (
	timeoutCircuitBreaker = 2000
	errorPercentThreshold = 25
	once                  = new(sync.Once)
)

func init() {
	once.Do(func() {
		hystrix.ConfigureCommand(commandGetAll, hystrix.CommandConfig{
			Timeout:               timeoutCircuitBreaker,
			ErrorPercentThreshold: errorPercentThreshold,
		})
		hystrix.ConfigureCommand(commandCreate, hystrix.CommandConfig{
			Timeout:               timeoutCircuitBreaker,
			ErrorPercentThreshold: errorPercentThreshold,
		})
	})
}

//NewPqRepository creates the pq repository for tax object with postgre connection.
func NewPqRepository(pool *sql.DB) taxobj.Repository {
	return &PqRepository{
		pool:      pool,
		statement: statement{},
	}
}

//GetAll return all tax objects in postgre.
func (repo *PqRepository) GetAll() (taxObjects []taxobj.TaxObject, err error) {
	const funcName = `[taxobj] [pq_repository] [GetAll]`
	taxObjects = make([]taxobj.TaxObject, 0)
	output := make(chan []taxobj.TaxObject, 1)
	chanError := hystrix.Go(commandGetAll, func() (err error) {
		defer close(output)
		result, _ := repo.getAll()
		output <- result
		return
	}, nil)
	select {
	case err = <-chanError:
	case taxObjects = <-output:
	}
	return
}

func (repo *PqRepository) getAll() (taxObjects []taxobj.TaxObject, err error) {
	var (
		taxObject taxobj.TaxObject
	)
	taxObjects = make([]taxobj.TaxObject, 0)

	//Lazy init for preparing statement
	if repo.statement.selectAll == nil {
		stmt, err := repo.pool.Prepare(querySelectAll)
		if err != nil {
			return taxObjects, err
		}
		repo.statement.selectAll = stmt
	}

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
	const funcName = `[taxobj] [pq_repository] [Create]`
	output := make(chan bool, 1)
	chanError := hystrix.Go(commandGetAll, func() (err error) {
		defer close(output)
		repo.create(taxObj)
		output <- true
		return
	}, nil)
	select {
	case err = <-chanError:
	case <-output:
	}
	return
}

func (repo *PqRepository) create(taxObj *taxobj.TaxObject) (err error) {
	//Lazy init for preparing statement
	if repo.statement.insert == nil {
		stmt, err := repo.pool.Prepare(queryInsert)
		if err != nil {
			return err
		}
		repo.statement.insert = stmt
	}
	row := repo.statement.insert.QueryRow(taxObj.Name, taxObj.TaxCode, taxObj.Price)

	err = row.Scan(
		&taxObj.ID,
	)

	return
}

//Migrate create the table in the database if it doesn't exist.
func (repo *PqRepository) Migrate() (err error) {
	var id int64
	row := repo.pool.QueryRow(querySelectOne)
	err = row.Scan(&id)
	if err == nil || err == sql.ErrNoRows {
		return nil
	}
	_, err = repo.pool.Exec(queryCreateTable)
	return
}

//Close close all prepared statements in this repository.
func (repo *PqRepository) Close() {
	if repo.statement.insert != nil {
		repo.statement.insert.Close()
	}
	if repo.statement.selectAll != nil {
		repo.statement.selectAll.Close()
	}
}
