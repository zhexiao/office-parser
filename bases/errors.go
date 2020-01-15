package bases

type errorSignal int

const (
	NormalError errorSignal = 1
	LibError    errorSignal = 2
	SystemError errorSignal = 99
)

type OpError struct {
	ErrorCode    errorSignal `json:"error_code"`
	ErrorMessage string      `json:"error_message"`
}

func NewOpError(code errorSignal, msg string) *OpError {
	return &OpError{
		ErrorCode:    code,
		ErrorMessage: msg,
	}
}

func (e *OpError) Error() string {
	return e.ErrorMessage
}
