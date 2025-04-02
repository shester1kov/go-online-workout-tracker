package apperrors

type AppError struct {
	Code    int
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

const (
	PgErrUniqueViolation     = "23505"
	PgErrForeignKeyViolation = "23503"
)
