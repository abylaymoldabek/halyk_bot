package domain

type Logger interface {
	Error(error, string)
	Info(string)
	Debug(string)
}
