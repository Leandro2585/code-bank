package service

import (
	"github.com/leandro2585/code-bank/dto"
	"github.com/leandro2585/code-bank/infra/grpc/pb"
	"github.com/leandro2585/code-bank/feature"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TransactionService struct {
	ProcessTransactionFeature feature.TransactionFeature
	pb.UnimplementedPaymentServiceServer
}

func NewTransactionService() *TransactionService {
	return &TransactionService{}
}


func (t *TransactionService) Payment(ctx  context.Context, in *pb.PaymentRequest) (*empty.Empty, error) {
	transactionDTO := dto.Transaction{
		Name: in.GetCreditCard().GetName(),
		Number: in.CreditCard.GetNumber(),
		ExpirationMonth: in.GetCreditCard().GetExpirationMonth(),
		ExpirationYear: in.GetCreditCard().GetExpirationYear(),
		CVV: in.GetCreditCard().GetCvv(),
		Amount: in.GetAmount(),
		Store: in.GetStore(),
		Description: in.GetDescription()
	}
	transaction, err := t.ProcessTransactionFeature.ProcessTransaction(transactionDTO)
	if err != nil {
		return &empty.Empty{}, status.Error(codes.FailedPrecondition, err.Error())
	}
	if transaction.Status != "approved" {
		return &empty.Empty{}, status.Error(codes.FailedPrecondition, "transaction rejected by the bank")
	}
	return &empty.Empty{}, nil 
}