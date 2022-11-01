package error

type ClientError struct {
	error
	errorCode string
	errorMsg  string
}

func NewClientError(errorCode string, errorMsg string, originalError error) ClientError {
	return ClientError{
		errorCode: errorCode,
		errorMsg:  errorMsg,
		error:     originalError,
	}
}
