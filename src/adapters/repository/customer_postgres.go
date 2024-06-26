package repository

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/realtobi999/GO_BankDemoApi/src/core/domain"
)


func (p *Postgres) GetCustomer(id uuid.UUID) (domain.Customer, error) {
    query := `SELECT * FROM customers WHERE id = $1 LIMIT 1`

    var customer domain.Customer

    err := p.DB.QueryRow(query, id).Scan(&customer.ID, &customer.FirstName, &customer.LastName, &customer.Birthday, &customer.Email, &customer.Phone, &customer.State, &customer.Address, &customer.CreatedAt, &customer.Token)
    if err != nil {
        return domain.Customer{}, err
    }

    return customer, nil
}

func (p *Postgres) GetAllCustomers(limit int, offset int) ([]domain.Customer, error) {
    query := `SELECT * FROM customers ORDER BY created_at LIMIT $1 OFFSET $2`

    rows, err := p.DB.Query(query, limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var customers []domain.Customer

    for rows.Next() {
        var customer domain.Customer

        if err := rows.Scan(&customer.ID, &customer.FirstName, &customer.LastName, &customer.Birthday, &customer.Email, &customer.Phone, &customer.State, &customer.Address, &customer.CreatedAt, &customer.Token); err != nil {
            return nil, err
        }

        customers = append(customers, customer)
    }
    
    if err := rows.Err(); err != nil{
        return nil, err
    }

    if len(customers) == 0 {
        return nil, sql.ErrNoRows
    }

    return customers, nil
}

func (p *Postgres) CreateCustomer(customer domain.Customer) (int64, error) {
    query := `
    INSERT INTO customers 
    (id, first_name, last_name, birthday, email, phone, state, address, created_at, token) 
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

    _, err := p.DB.Exec(query, customer.ID.String(), customer.FirstName, customer.LastName, customer.Birthday, customer.Email, customer.Phone, customer.State, customer.Address, customer.CreatedAt, customer.Token)
    if err != nil {
        return 0, err
    }

    return 1, nil
}

func (p *Postgres) UpdateCustomer(customer domain.Customer) (int64, error) {
    query := `
    UPDATE customers
    SET first_name = $1, last_name = $2, birthday = $3, email = $4, phone = $5, state = $6, address = $7
    WHERE id = $8`

    result, err := p.DB.Exec(query, customer.FirstName, customer.LastName, customer.Birthday, customer.Email, customer.Phone, customer.State, customer.Address, customer.ID)
    if err != nil {
        return 0, err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return 0, err
    }

    return rowsAffected, nil
}

func (p *Postgres) DeleteCustomer(customerID uuid.UUID) (int64, error) {
    query := `DELETE FROM customers WHERE id = $1`

    result, err := p.DB.Exec(query, customerID)
    if err != nil {
        return 0, err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return 0, err
    }

    return rowsAffected, nil
}

func (p *Postgres) AuthCustomer(customerID uuid.UUID, token string) (bool, error) {
    query := `SELECT EXISTS(SELECT 1 FROM customers WHERE id = $1 AND token = $2)`

    var exists bool
    err := p.DB.QueryRow(query, customerID, token).Scan(&exists)
    switch {
    case err == sql.ErrNoRows:
        return false, nil
    case err != nil:
        return false, err
    default:
        return exists, nil
    }
}