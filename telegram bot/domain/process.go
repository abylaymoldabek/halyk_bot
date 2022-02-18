package domain

type Process struct {
	// Value                interface{}
	Id                   string // ProcessID
	State                string // ACTIVE or COMPLETED
	ProcessDefinitionKey string // onboarding01, etc.
}
