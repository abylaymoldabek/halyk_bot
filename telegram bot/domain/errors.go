package domain

import "errors"

var (
	ErrIncidentURLNotFound   = errors.New("failed to get URL for getting incidents")
	ErrNoDataFound           = errors.New("failed to get any processes on given client")
	ErrNoIncidentFound       = errors.New("failed to get incident(s)")
	ErrProcessIDNotFound     = errors.New("failed to get process ID")
	ErrProcessNotFound       = errors.New("failed to get requested processes for given client")
	ErrNoVarsFound           = errors.New("failed to get process variables")
	ErrProcessStatusNotFound = errors.New("failed to get process status")
	ErrProcessURLNotFound    = errors.New("failed to get process URL")
	ErrProcessesURLNotFound  = errors.New("failed to get processes URL")
	ErrRetriesURLNotFound    = errors.New("failed to get retries URL")
	ErrTokenNotFound         = errors.New("failed to get token")
	ErrTokenURLNotFound      = errors.New("failed to get token URL")
	ErrUnauthorized          = errors.New("failed to authorize")
	ErrUnknownIncident       = errors.New("couldn't recognize incident type")
)
