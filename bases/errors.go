package bases

type errorSignal int

const (
	NormalError errorSignal = 1
	LibError    errorSignal = 2
	SystemError errorSignal = 99
)

type OpError struct {
	ErrorCode errorSignal `json:"error_code"`
	Message   string      `json:"message"`
}

func NewOpError(code errorSignal, msg string) *OpError {
	return &OpError{
		ErrorCode: code,
		Message:   msg,
	}
}

func (e *OpError) Error() string {
	return e.Message
}
