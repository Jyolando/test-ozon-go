package entities

const (
	HTTP400 = "invalid link"
	HTTP404 = "element not found"
	HTTP500 = "server error"
)

type RunError struct {
	error
}

type MissingStorageTypeError struct {
	error
}

type IncorrectPsqlStorage struct {
	error
}
