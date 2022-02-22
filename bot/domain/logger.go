package domain

type Logger interface {
	Error(error)
	Info(string)
	Debug(string)
}
