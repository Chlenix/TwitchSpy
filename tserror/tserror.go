package tserror

// TODO: HashMap of errors and their consecutive count,
// if an error type exceeeds allowed # of warning types,
// promote to suspension and suspend a module

const (
	Critical = 2 // Terminate App
	Suspend  = 1 // Stop Module
	Warning  = 0 // Proceed
)

func New(err error, level int) *AppError {
	return &AppError{
		e:     err,
		Level: level,
	}
}

type AppError struct {
	e     error
	Level int
}

func (e AppError) Error() string {
	return e.Error()
}
