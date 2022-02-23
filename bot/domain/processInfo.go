package domain

import "context"

type ResponseStruct struct {
	Status, Message string
}

type ProcessUsecase interface {
	ProcessRequest(ctx context.Context, searchCriteria Criteria) (*ResponseStruct, error)
}

type ProcessRepository interface {
	GetProcess(ctx context.Context, searchCriteria Criteria) (*Process, error)
	GetProcessStatus(ctx context.Context, processID string) (string, error)
	RetryJobOrTask(ctx context.Context, processID string) error
}
