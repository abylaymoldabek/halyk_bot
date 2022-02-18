package usecase

import (
	"context"
	"support/domain"
	"time"
)

type processUsecase struct {
	processRepo    domain.ProcessRepository
	contextTimeout time.Duration
}

// NewProcessUsecase will create new processUsecase object representation of domain.ProcessUsecase interface
func NewProcessUsecase(p domain.ProcessRepository, timeout time.Duration) domain.ProcessUsecase {
	return &processUsecase{
		processRepo:    p,
		contextTimeout: timeout,
	}
}

func (p *processUsecase) MainLogic(c context.Context, searchCriteria domain.Criteria) (*domain.ResponseStruct, error) {
	ctx, cancel := context.WithTimeout(c, p.contextTimeout)
	defer cancel()
	process, err := p.processRepo.GetProcess(ctx, searchCriteria)
	if err != nil {
		return nil, err
	}
	if process.State == "ACTIVE" {
		if err := p.processRepo.RetryJobOrTask(ctx, process.Id); err != nil {
			return nil, err
		}
		return &domain.ResponseStruct{Message: "You may check now"}, nil
	}
	status, err := p.processRepo.GetProcessStatus(ctx, process.Id)
	if err != nil {
		return nil, err
	}
	return &domain.ResponseStruct{Status: status}, nil
}

func (c *Client) getToken() string {
	c.token.mx.RLock() 
	defer c.token.mx.RUnlock()
	token := c.token.tokenString
	return token
}