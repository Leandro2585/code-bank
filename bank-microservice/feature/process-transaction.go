package feature

import (
	"time"

	"github.com/leandro2585/codebank/domain"
	"github.com/leandro2585/codebank/dto"
)

type TransactionFeature struct {
	TransactionRepository domain.TransactionRepository
	KafkaProducer kafka.KafkaProducer
}

func NewTransactionFeature(transactionRepository domain.TransactionRepository) TransactionFeature {
	return TransactionFeature{TransactionRepository: transactionRepository}
}

func (f TransactionFeature) ProcessTransaction(transactionDTO dto.Transaction) (domain.Transaction, error) {
	creditCard := f.hydrateCreditCard(transactionDTO)
	ccBalanceLimit, err := f.TransactionRepository.GetCreditCard(*creditCard)
	if err != nil {
		return domain.Transaction{}, err
	}
	creditCard.ID = ccBalanceLimit.ID
	creditCard.Limit = ccBalanceLimit.Limit
	creditCard.Balance = ccBalanceLimit.Balance
	t := f.newTransaction(transactionDTO, ccBalanceLimit)
	t.ProcessAndValidate(creditCard)

	err = f.TransactionRepository.SaveTransaction(*t, *creditCard)
	if err != nil {
		return domain.Transaction{}, err
	}

	transactionDTO.ID = t.ID
	transactionDTO.CreatedAt = t.CreatedAt
	transactionJSON, err := json.Marshal(transactionDTO)
	if err != nil {
		return domain.Transaction{}, err
	}
	err = f.KafkaProducer.Publish(string(transactionJSON), "payments")
	if err != nil {
		return domain.Transaction{}, err
	}
	return *t, nil
}

func (f TransactionFeature) hydrateCreditCard(transactionDTO dto.Transaction) *domain.CreditCard {
	creditCard := domain.NewCreditCard()
	creditCard.Name = transactionDTO.Name
	creditCard.Number = transactionDTO.Number
	creditCard.ExpirationMonth = transactionDTO.ExpirationMonth
	creditCard.ExpirationYear = transactionDTO.ExpirationYear
	creditCard.CVV = transactionDTO.CVV
	creditCard.Balance = transactionDTO.Amount
	return creditCard
}

func (f TransactionFeature) newTransaction(transaction dto.Transaction, cc domain.CreditCard) *domain.Transaction {
	t := domain.NewTransaction()
	t.CreditCardId = cc.ID
	t.Amount = transaction.Amount
	t.Store = transaction.Store
	t.Description = transaction.Description
	t.CreatedAt = time.Now()
	return t
}
