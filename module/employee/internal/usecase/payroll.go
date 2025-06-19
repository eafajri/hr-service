package usecase

import (
	"log"
	"time"

	"github.com/eafajri/hr-service.git/module/employee/internal/entity"
	"go.uber.org/zap"
)

//go:generate mockery --name PayrollUseCase --output ./mocks
type PayrollUseCase interface {
	ClosePayrollPeriod(periodID int64) error
	GetPayslip(userID int64, periodID int64) (entity.PayrollPayslip, error)
	GetPayslips(periodID int64) ([]entity.PayrollPayslip, error)
}

type PayrollUseCaseImpl struct {
	payrollRepository PayrollRepository
	userUseCase       UserUseCase
}

func NewPayrollUseCase(
	payrollRepository PayrollRepository,
	userUseCase UserUseCase,
) *PayrollUseCaseImpl {
	return &PayrollUseCaseImpl{
		payrollRepository: payrollRepository,
		userUseCase:       userUseCase,
	}
}

func (p *PayrollUseCaseImpl) ClosePayrollPeriod(periodID int64) error {
	err := p.payrollRepository.ClosePayrollPeriod(periodID)
	if err != nil {
		log.Println(
			"error when GetPayslip",
			zap.String("method", "PayrollUseCaseImpl.ClosePayrollPeriod"),
			zap.Int64("period_id", periodID),
			zap.Error(err),
		)
		return err
	}

	return nil
}

func (p *PayrollUseCaseImpl) GetPayslip(userID int64, periodID int64) (entity.PayrollPayslip, error) {
	if periodID == 0 {
		period, err := p.payrollRepository.GetPeriodByEntityDate(time.Now())
		if err != nil {
			log.Println(
				"error when GetPeriodByEntityDate",
				zap.String("method", "PayrollUseCaseImpl.GetPayslip"),
				zap.Int64("user_id", userID),
				zap.Error(err),
			)
			return entity.PayrollPayslip{}, err
		}

		periodID = period.ID
	}

	payslip, err := p.payrollRepository.GetPayslip(userID, periodID)
	if err != nil {
		log.Println(
			"error when GetPayslip",
			zap.String("method", "PayrollUseCaseImpl.GetPayslip"),
			zap.Int64("user_id", userID),
			zap.Int64("period_id", periodID),
			zap.Error(err),
		)
		return entity.PayrollPayslip{}, err
	}

	return payslip, nil
}

func (p *PayrollUseCaseImpl) GetPayslips(periodID int64) ([]entity.PayrollPayslip, error) {
	if periodID == 0 {
		period, err := p.payrollRepository.GetPeriodByEntityDate(time.Now())
		if err != nil {
			log.Println(
				"error when GetPeriodByEntityDate",
				zap.String("method", "PayrollUseCaseImpl.GetPayslips"),
				zap.Int64("period_id", periodID),
				zap.Error(err),
			)
			return nil, err
		}

		periodID = period.ID
	}

	payslips, err := p.payrollRepository.GetPayslips(periodID)
	if err != nil {
		log.Println(
			"error when GetPayslips",
			zap.String("method", "PayrollUseCaseImpl.GetPayslips"),
			zap.Int64("period_id", periodID),
			zap.Error(err),
		)
		return nil, err
	}

	return payslips, nil
}
