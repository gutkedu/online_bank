package service

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/gutkedu/code_bank/dto"
	"github.com/gutkedu/code_bank/infra/grpc/pb"
	usecase "github.com/gutkedu/code_bank/useCase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TransactionService struct {
	ProcessTransactionUseCase usecase.UseCaseTransaction
	pb.UnimplementedPaymentServiceServer
}

func NewTransactionService() *TransactionService {
	return &TransactionService{}
}

func (t *TransactionService) Payment(ctx context.Context, in *pb.PaymentRequest) (*empty.Empty, error) {
	transactionDto := dto.Transaction{
		Name:            in.GetCreditCard().GetName(),
		Number:          in.CreditCard.GetNumber(),
		ExpirationMonth: in.GetCreditCard().ExpirationMonth,
		ExpirationYear:  in.CreditCard.GetExpirationYear(),
		CVV:             in.GetCreditCard().GetCvv(),
		Amount:          in.GetAmount(),
		Store:           in.GetStore(),
		Description:     in.GetDescription(),
	}
	transaction, err := t.ProcessTransactionUseCase.ProcessTransaction(transactionDto)
	if err != nil {
		return &empty.Empty{}, status.Error(codes.FailedPrecondition, err.Error())
	}
	if transaction.Status != "approved" {
		return &empty.Empty{}, status.Error(codes.FailedPrecondition, "transaction rejected by the bank")
	}
	return &empty.Empty{}, nil
}
