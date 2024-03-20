package repositories

import (
	"crypto-challenge/entities"
	"database/sql"
	"log"
)

type TransactionMySqlRepository struct {
	db *sql.DB
}

func NewTransactionMySqlRepository(db *sql.DB) *TransactionMySqlRepository {
	return &TransactionMySqlRepository{db}
}

func (r *TransactionMySqlRepository) Create(newTransaction *entities.Transaction) error {
	query := "INSERT INTO transactions (id, user_document, credit_card_token, `value`) VALUES (?, ?, ?, ?)"

	_, err := r.db.Exec(query, newTransaction.ID, newTransaction.UserDocument,
		newTransaction.CreditCardToken, newTransaction.Value)
	if err != nil {
		log.Println(err)
	}

	return err
}

func (r *TransactionMySqlRepository) FindByID(idToSearch string) (*entities.Transaction, error) {
	query := "SELECT id, user_document, credit_card_token, `value` FROM transactions WHERE id = ?"

	var (
		id, userDocument, creditCardToken string
		value                             float64
	)

	err := r.db.QueryRow(query, idToSearch).Scan(&id, &userDocument, &creditCardToken, &value)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &entities.Transaction{
		ID:              id,
		UserDocument:    userDocument,
		CreditCardToken: creditCardToken,
		Value:           value,
	}, nil
}

func (r *TransactionMySqlRepository) FindAll() ([]*entities.Transaction, error) {
	query := "SELECT id, user_document, credit_card_token, `value` FROM transactions"

	rows, err := r.db.Query(query)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer rows.Close()

	foundTransactions := make([]*entities.Transaction, 0, 5)

	for rows.Next() {
		var (
			id, userDocument, creditCardToken string
			value                             float64
		)

		err = rows.Scan(&id, &userDocument, &creditCardToken, &value)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		foundTransactions = append(foundTransactions,
			&entities.Transaction{ID: id, UserDocument: userDocument, CreditCardToken: creditCardToken, Value: value})
	}

	rows.Close()

	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return foundTransactions, nil
}

func (r *TransactionMySqlRepository) UpdateByID(updatedTransaction *entities.Transaction) error {
	query := "UPDATE transactions SET user_document = ?, credit_card_token = ?, `value` = ?  WHERE id = ?"

	_, err := r.db.Exec(query, updatedTransaction.UserDocument, updatedTransaction.CreditCardToken, updatedTransaction.Value, updatedTransaction.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *TransactionMySqlRepository) DeleteByID(idToDelete string) error {
	query := "DELETE FROM transactions WHERE id = ?"

	_, err := r.db.Exec(query, idToDelete)
	if err != nil {
		return err
	}

	return nil
}
