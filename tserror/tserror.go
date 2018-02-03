package tserror

import "fmt"

// TODO: HashMap of errors and their consecutive count,
// if an error type exceeeds allowed # of warning types,
// promote to suspension and suspend a module


type AppError struct {
	ErrorObject error
	Level       int
}

type ErrorTable map[error]int

const (
	Critical = 2  // Terminate App
	Suspend  = 1  // Stop Module
	Warning  = 0  // Proceed
	Ignore   = -1 // Ignore and proceed

	MaxWarnings = 5
)

func New(err error, level int) *AppError {
	return &AppError{
		ErrorObject: err,
		Level:       level,
	}
}

func (errorTable ErrorTable) Handle(e AppError) {
	switch e.Level {
	case Warning:
		errorTable[e.ErrorObject]++
		if errorTable[e.ErrorObject] >= MaxWarnings {
			fmt.Errorf("%s\n", e.Error())
		}
		break
	case Critical:
		panic(e.Error())
	}
}

func (e AppError) Error() string {
	return e.ErrorObject.Error()
}
