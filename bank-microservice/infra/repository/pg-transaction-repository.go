package repository

import (
	"database/sql"
	"errors"

	"github.com/leandro2585/codebank/domain"
)

type PgTransactionRepository struct {
	db *sql.DB
}

func NewPgTransactionRepository(db *sql.DB) *PgTransactionRepository {
	return &PgTransactionRepository{db: db}
}

func (t *PgTransactionRepository) SaveTransaction(transaction domain.Transaction, creditCard domain.CreditCard) error {
	stmt, err := t.db.Prepare(`INSERT INTO transactions(id, credit_card_id, amount, status, description, store, created_at)
							values ($1,$2,$3,$4,$5,$6,$7)`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(transaction.ID, creditCard.ID, transaction.Amount, transaction.Status, transaction.Description, transaction.Store, transaction.CreatedAt)
	if err != nil {
		return err
	}
	if transaction.Status == "approved" {
	}
	err = stmt.Close()
	if err != nil {
		return err
	}
	return nil
}

func (t *PgTransactionRepository) updateBalance(creditCard domain.CreditCard) error {
	_, err := t.db.Exec("update credit_cards set balance = $1 where id = $2",
		creditCard.Balance, creditCard.ID)
	if err != nil {
		return err
	}
	return nil
}

func (t *PgTransactionRepository) CreateCreditCard(creditCard domain.CreditCard) error {
	stmt, err := t.db.Prepare(`insert into credit_cards(id, name, number, expiration_month,expiration_year, CVV,balance, balance_limit) 
								values($1,$2,$3,$4,$5,$6,$7,$8)`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(
		creditCard.ID,
		creditCard.Name,
		creditCard.Number,
		creditCard.ExpirationMonth,
		creditCard.ExpirationYear,
		creditCard.CVV,
		creditCard.Balance,
		creditCard.Limit,
	)
	if err != nil {
		return err
	}
	err = stmt.Close()
	if err != nil {
		return err
	}
	return nil
}

func (t *PgTransactionRepository) GetCreditCard(creditCard domain.CreditCard) (domain.CreditCard, error) {
	var c domain.CreditCard
	stmt, err := t.db.Prepare(`SELECT id, balance, balance_limit from credit_cards where id = $1`)
	if err != nil {
		return c, err
	}
	if err = stmt.QueryRow(creditCard.Number).Scan(&c.ID, &c.Balance, &c.Limit); err != nil {
		return c, errors.New("credit card does not exists")
	}
	return c, nil
}
