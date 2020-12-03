package nerrors

import (
	"fmt"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"reflect"
	"runtime"
	"strings"

	"google.golang.org/grpc/status"
)

// ExtendedError with an extended golang error
type ExtendedError struct {
	Code       ErrorCode // Code with the type of the error (compatible with GRPC)
	Msg        string    // Error Msg
	From       error     // Parent error
	StackTrace []string  // Stack trace related to where the error happened in the code base
}

// NewExtendedError generic method to create an extended error
func NewExtendedError(code ErrorCode, format string, a ...interface{}) *ExtendedError {
	return &ExtendedError{
		Code:       code,
		Msg:        formatMsg(format, a...),
		StackTrace: getStackTrace(),
	}
}
// NewExtendedError From generic method to create an extended error from another caused by another one
func NewExtendedErrorFrom(code ErrorCode, err error,
	format string, a ...interface{}) *ExtendedError {
	return &ExtendedError{
		Code:       code,
		Msg:        formatMsg(format, a...),
		From:       err,
		StackTrace: getStackTrace(),
	}
}
// Error method to implement error interface
func (ee *ExtendedError) Error() string {
	return ee.String()
}

func (ee *ExtendedError) String() string {
	msg := ee.ShortString()
	if ee.From != nil {
		return fmt.Sprintf("%s caused by %s", msg, ee.From.Error())
	}
	return msg
}

func (ee *ExtendedError) ShortString() string {
	return fmt.Sprintf("[%s] %s", ee.Code.String(), ee.Msg)
}

// Unwrap method to implement Wrapper interface (provides context around another error)
func (ee *ExtendedError) Unwrap() error {
	return ee.From
}

// StackTraceToString loops through error chain showing stack trace
func (ee *ExtendedError) StackTraceToString() string {
	if ee == nil {
		return ""
	}
	traces := ee.ShortString() + "\n" + strings.Join(ee.StackTrace, "")
	if ee.From != nil {
		traces += "Caused by "
		var pp *ExtendedError
		if reflect.TypeOf(ee.From) == reflect.TypeOf(pp) {
			traces += ee.From.(*ExtendedError).StackTraceToString()
		} else {
			traces += ee.From.Error() + "\n" + " <stack trace no available>"
		}
	}
	return traces

}

// getStackTrace get the stack trace when an error occurs
func getStackTrace() []string {
	buf := make([]uintptr, 32)
	callers := runtime.Callers(2, buf)
	stackTrace := make([]string, callers)
	for i := 0; i < callers; i++ {
		frames := runtime.FuncForPC(buf[i])
		filePath, line := frames.FileLine(buf[i])
		stackTrace[i] = fmt.Sprintf("%s:%d - %s\n", filePath, line, frames.Name())
	}
	return stackTrace
}

// ---------------
// ToGRPC converts an extended error to a GrpcError
func (ee *ExtendedError) ToGRPC() error {
	code, exists := ToGRPCCode[ee.Code]
	if !exists {
		fmt.Printf("Code (%d - %s) does not exist\n", ee.Code, ee.Code.String())
		return ee
	}

	status := status.New(code, ee.Msg)

	complexSt, err := status.WithDetails(
		&errdetails.DebugInfo{
			StackEntries: ee.StackTrace,
			Detail:       ee.Error(),
		})
	if err != nil {
		fmt.Printf("Error converting error to GRPC: %s\n", err.Error())
		return ee
	}
	return complexSt.Err()
}

// ToGRPC converts a GrpcError to an extended error
func ToGRPC(err error) *ExtendedError {
	status := status.Convert(err)

	status.Details()

	return &ExtendedError{
		Code:       FromGRPCCode[status.Code()],
		Msg:        status.Message(),
		From:       nil,
		StackTrace: nil,
	}
}

// ---------------
func formatMsg(format string, a ...interface{}) string {
	return fmt.Sprintf(format, a...)
}
func NewCanceledError(format string, a ...interface{}) *ExtendedError {
	return NewExtendedError(Canceled, format, a...)
}
func NewCanceledErrorFrom(err error, format string, a ...interface{}) *ExtendedError {
	return NewExtendedErrorFrom(Canceled, err, format, a...)
}
func NewUnknownError(format string, a ...interface{}) *ExtendedError {
	return NewExtendedError(Unknown, format, a...)
}
func NewUnknownErrorFrom(err error, format string, a ...interface{}) *ExtendedError {
	return NewExtendedErrorFrom(Unknown, err, format, a...)
}

func NewInvalidArgumentError(format string, a ...interface{}) *ExtendedError {
	return NewExtendedError(InvalidArgument, format, a...)
}
func NewInvalidArgumentErrorFrom(err error, format string, a ...interface{}) *ExtendedError {
	return NewExtendedErrorFrom(InvalidArgument, err, format, a...)
}

func NewDeadlineExceededError(format string, a ...interface{}) *ExtendedError {
	return NewExtendedError(DeadlineExceeded, format, a...)
}
func NewDeadlineExceededErrorFrom(err error, format string, a ...interface{}) *ExtendedError {
	return NewExtendedErrorFrom(DeadlineExceeded, err, format, a...)
}

func NewNotFoundError(format string, a ...interface{}) *ExtendedError {
	return NewExtendedError(NotFound, format, a...)
}
func NewNotFoundErrorFrom(err error, format string, a ...interface{}) *ExtendedError {
	return NewExtendedErrorFrom(NotFound, err, format, a...)
}

func NewAlreadyExistsError(format string, a ...interface{}) *ExtendedError {
	return NewExtendedError(AlreadyExists, format, a...)
}
func NewAlreadyExistsErrorFrom(err error, format string, a ...interface{}) *ExtendedError {
	return NewExtendedErrorFrom(AlreadyExists, err, format, a...)
}

func NewPermissionDeniedError(format string, a ...interface{}) *ExtendedError {
	return NewExtendedError(PermissionDenied, format, a...)
}
func NewPermissionDeniedErrorFrom(err error, format string, a ...interface{}) *ExtendedError {
	return NewExtendedErrorFrom(PermissionDenied, err, format, a...)
}

func NewResourceExhaustedError(format string, a ...interface{}) *ExtendedError {
	return NewExtendedError(ResourceExhausted, format, a...)
}
func NewResourceExhaustedErrorFrom(err error, format string, a ...interface{}) *ExtendedError {
	return NewExtendedErrorFrom(ResourceExhausted, err, format, a...)
}

func NewFailedPreconditionError(format string, a ...interface{}) *ExtendedError {
	return NewExtendedError(FailedPrecondition, format, a...)
}
func NewFailedPreconditionErrorFrom(err error, format string, a ...interface{}) *ExtendedError {
	return NewExtendedErrorFrom(FailedPrecondition, err, format, a...)
}

func NewAbortedError(format string, a ...interface{}) *ExtendedError {
	return NewExtendedError(Aborted, format, a...)
}
func NewAbortedErrorFrom(err error, format string, a ...interface{}) *ExtendedError {
	return NewExtendedErrorFrom(Aborted, err, format, a...)
}

func NewOutOfRangeError(format string, a ...interface{}) *ExtendedError {
	return NewExtendedError(OutOfRange, format, a...)
}
func NewOutOfRangeErrorFrom(err error, format string, a ...interface{}) *ExtendedError {
	return NewExtendedErrorFrom(OutOfRange, err, format, a...)
}

func NewUnimplementedError(format string, a ...interface{}) *ExtendedError {
	return NewExtendedError(Unimplemented, format, a...)
}
func NewUnimplementedErrorFrom(err error, format string, a ...interface{}) *ExtendedError {
	return NewExtendedErrorFrom(Unimplemented, err, format, a...)
}

func NewInternalError(format string, a ...interface{}) *ExtendedError {
	return NewExtendedError(Internal, format, a...)
}
func NewInternalErrorFrom(err error, format string, a ...interface{}) *ExtendedError {
	return NewExtendedErrorFrom(Internal, err, format, a...)
}

func NewUnavailableError(format string, a ...interface{}) *ExtendedError {
	return NewExtendedError(Unavailable, format, a...)
}
func NewUnavailableErrorFrom(err error, format string, a ...interface{}) *ExtendedError {
	return NewExtendedErrorFrom(Unavailable, err, format, a...)
}

func NewDataLossError(format string, a ...interface{}) *ExtendedError {
	return NewExtendedError(DataLoss, format, a...)
}
func NewDataLossErrorFrom(err error, format string, a ...interface{}) *ExtendedError {
	return NewExtendedErrorFrom(DataLoss, err, format, a...)
}

func NewUnauthenticatedError(format string, a ...interface{}) *ExtendedError {
	return NewExtendedError(Unauthenticated, format, a...)
}
func NewUnauthenticatedErrorFrom(err error, format string, a ...interface{}) *ExtendedError {
	return NewExtendedErrorFrom(Unauthenticated, err, format, a...)
}
