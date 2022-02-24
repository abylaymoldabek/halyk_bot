package domain

import "context"

type ResponseStruct struct {
	Status, Message string
}

type ProcessUsecase interface {
	ProcessRequest(ctx context.Context, searchCriteria Criteria) (*ResponseStruct, error)
}

type ProcessRepository interface {
	GetActivityID(ctx context.Context, processID string, activityName string) (string, error)
	GetProcess(ctx context.Context, searchCriteria Criteria) (*Process, error)
	GetProcessStatus(ctx context.Context, processID string) (string, error)
	GetRole(ctx context.Context, tab string) error
	RetryJobOrTask(ctx context.Context, processID string) error
	Redo(ctx context.Context, processID, activityID string) error
	UpdateBranch(ctx context.Context, processID, branchCode string) error
}
