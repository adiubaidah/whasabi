package exception

type NotFoundError struct {
	Error string
}

func NewNotFoundError(error string) NotFoundError {
	return NotFoundError{Error: error}
}

type BadRequestError struct {
	Error string
}

func NewBadRequestError(error string) BadRequestError {
	return BadRequestError{Error: error}
}

type ForbiddenError struct {
	Error string
}

func NewForbiddenError(error string) ForbiddenError {
	return ForbiddenError{Error: error}
}

type UnauthorizedError struct {
	Error string
}

func NewUnauthorizedError(error string) UnauthorizedError {
	return UnauthorizedError{Error: error}
}
