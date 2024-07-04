package processor

import (
	"cashapp/models"

	"cashapp/repository"
	"fmt"
	"log"

	"gorm.io/gorm"
)

type Processor struct {
	Repo repository.Repo
}

func New(r repository.Repo) Processor {
	return Processor{
		Repo: r,
	}
}

func (p *Processor) ProcessTransaction(fromTrans models.Transaction) error {
	switch fromTrans.Purpose {
	case models.Deposit:
		if err := p.DepositMoneyToProvider(fromTrans); err != nil {
			return fmt.Errorf("money deposit failed. %v", err)
		}
	case models.Withdrawal:
		if err := p.WithdrawMoneyFromProvider(fromTrans); err != nil {
			return fmt.Errorf("money withdrawal failed. %v", err)
		}
	default:
		log.Println("no handler for purpose. purpose=", fromTrans.Purpose)
	}
	return nil
}

func (p *Processor) SuccessCallback(fromTrans, toTrans *models.Transaction) error {
	fromTrans.Status = models.Success
	toTrans.Status = models.Success

	return p.Repo.Transactions.SQLTransaction(func(tx *gorm.DB) error {
		return p.Repo.Transactions.Updates(tx, fromTrans, toTrans)
	})
}

func (p *Processor) FailureCallback(fromTrans, toTrans *models.Transaction, err error) error {
	fromTrans.Status = models.Failed
	toTrans.Status = models.Failed
	fromTrans.FailureReason = err.Error()
	toTrans.FailureReason = err.Error()

	return p.Repo.Transactions.SQLTransaction(func(tx *gorm.DB) error {
		return p.Repo.Transactions.Updates(tx, fromTrans, toTrans)
	})
}
