package usecase

import (
	"context"
	"log"
	"time"
	"v/domain"
)

const (
	UVK       = "Подпроцесс УВК"
	Obrabotka = "Обработка перевода при обращении клиента"
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

// ProcessRequest returns status for completed processes and retries incidents, if any
func (p *processUsecase) ProcessRequest(c context.Context, searchCriteria domain.Criteria) (*domain.ResponseStruct, error) {
	ctx, cancel := context.WithTimeout(c, p.contextTimeout)
	defer cancel()
	// Check manager's role (basic, needs updating)
	if err := p.processRepo.GetRole(ctx, searchCriteria.Tab); err != nil {
		return nil, err
	}
	// Process search
	process, err := p.processRepo.GetProcess(ctx, searchCriteria)
	if err != nil {
		return nil, err
	}
	log.Println("ProcessID:", process.Id)
	// Search for incidents if process is active
	if process.State == "ACTIVE" {
		err := p.processRepo.RetryJobOrTask(ctx, process.Id)
		if err != nil {
			if err != domain.ErrNoIncidentFound {
				return nil, err
			}
		} else {
			return &domain.ResponseStruct{Message: "You may check now"}, nil
		}
	}
	// Fetch and return process status
	status, err := p.processRepo.GetProcessStatus(ctx, process.Id)
	if err != nil {
		return nil, err
	}
	return &domain.ResponseStruct{Status: status}, nil
}

// ProcessTransfer reattempts UVK or updates branchSapCode and initRole and then repeats Obrabotka...
func (p *processUsecase) ProcessTransfer(c context.Context, searchCriteria domain.Criteria) (*domain.ResponseStruct, error) {
	ctx, cancel := context.WithTimeout(c, p.contextTimeout)
	defer cancel()
	process, err := p.processRepo.GetProcess(ctx, searchCriteria)
	if err != nil {
		return nil, err
	}
	log.Println("ProcessID:", process.Id)
	activityID, err := p.processRepo.GetActivityID(ctx, process.Id, searchCriteria.ActivityName)
	if err != nil {
		return nil, err
	}
	// in a separate goroutine?
	if searchCriteria.ActivityName == Obrabotka {
		if err := p.processRepo.UpdateBranch(ctx, process.Id, searchCriteria.BranchCode); err != nil {
			return nil, err
		}
	}
	if err := p.processRepo.Redo(ctx, process.Id, activityID); err != nil {
		return nil, err
	}
	return &domain.ResponseStruct{Message: "You may check now"}, nil // TODO!!!
}
