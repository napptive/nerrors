package nerrors

import (
	"fmt"
	"github.com/napptive/grpc-common-go"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/runtime/protoiface"
	"reflect"
	"runtime"
	"strings"
)

// ExtendedError with an extended golang error
type ExtendedError struct {
	// Code with the type of the error (compatible with GRPC)
	Code ErrorCode
	// Msg with a textual description of the error.
	Msg string
	// From links with the parent error if any.
	From error
	// StackTrace related to where the error happened in the code base.
	StackTrace []string
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

// getDetails converts a Extended Message into a list of Proto Message
// The detail is Code: ... - Msg: ...
func (ee *ExtendedError) getDetails(list []protoiface.MessageV1) []protoiface.MessageV1 {

	debugInfo := &grpc_common_go.ErrorDetails{
		StackEntries: ee.StackTrace,
		Detail:       fmt.Sprintf("Code: %s - Msg: %s", ee.Code.String(), ee.Msg),
	}
	if ee.From == nil {
		return append(list, debugInfo)
	} else {
		return append(FromError(ee.From).getDetails(list), debugInfo)
	}
}


// TODO: in the next version, instead of use DebugInfo and compose a detail, we can implement our own protoiface.MessageV1
// ToGRPC converts an extended error to a GrpcError
func (ee *ExtendedError) ToGRPC() error {
	code, exists := ToGRPCCode[ee.Code]
	if !exists {
		fmt.Printf("Code (%d - %s) does not exist\n", ee.Code, ee.Code.String())
		return ee
	}

	st := status.New(code, ee.Msg)

	// we create as many details as errors we have in the chain. This is the way to convert a GPRC to Extended Error again
	details := make([]protoiface.MessageV1, 0)
	allDetails := ee.getDetails(details)

	complexSt, err := st.WithDetails(allDetails...)

	if err != nil {
		fmt.Printf("Error converting error to GRPC: %s\n", err.Error())
		return ee
	}
	return complexSt.Err()
}

// FromGRPC converts a GrpcError to an extended error
func FromGRPC(err error) *ExtendedError {
	st := status.Convert(err)
	code := st.Code()

	if len(st.Details()) == 0 {
		return &ExtendedError{
			Code:       FromGRPCCode[code],
			Msg:        st.Message(),
			From:       nil,
			StackTrace: getStackTrace(),
		}
	}

	extended := ExtendedErrorFromDetail(st.Details())
	extended.Code = FromGRPCCode[code]

	return extended

}

//
// getCodeFromGRPCMsg try to get the error code and the message if the details has the format belong
// Detail: fmt.Sprintf("Code: %s - Msg: %s", ee.Code.String(), ee.Msg),
func getCodeFromGRPCMsg(msg string) (string, ErrorCode) {

	var errCode ErrorCode
	var message string
	// if we recognize the msg as our conversion, we can get the error code and the message
	if ind := strings.Index(msg, "Code: "); ind >= 0 {
		if msgInd := strings.Index(msg, " - Msg:"); msgInd >= 0 {
			code := msg[ind+6 : msgInd]
			errCode = FromStringCode[code]

			message = msg[msgInd+8:]
			return message, errCode
		}
	}

	return msg, Unknown
}

// ExtendedErrorFromDetail create an extended error from the details of the grpc error
func ExtendedErrorFromDetail(details []interface{}) *ExtendedError {
	if len(details) == 0 {
		return nil
	}

	if len(details) == 1 {
		if e, ok := details[len(details)-1].(*grpc_common_go.ErrorDetails); ok {
			info := e
			stackTrace := info.StackEntries
			msg, code := getCodeFromGRPCMsg(info.Detail)
			return &ExtendedError{
				Msg:        msg,
				Code:       code,
				From:       nil,
				StackTrace: stackTrace,
			}
		}
	} else {
		if e, ok := details[len(details)-1].(*grpc_common_go.ErrorDetails); ok {
			info := e
			stackTrace := info.StackEntries
			msg, code := getCodeFromGRPCMsg(info.Detail)
			return &ExtendedError{
				Msg:        msg,
				Code:       code,
				From:       ExtendedErrorFromDetail(details[0 : len(details)-1]),
				StackTrace: stackTrace,
			}
		}
	}
	return nil

}

// FromError transforms a standard go error into an extended error
func FromError(err error) *ExtendedError {

	if e, ok := err.(*ExtendedError); ok {
		return e
	}

	return &ExtendedError{
		Code:       Unknown,
		Msg:        err.Error(),
		From:       nil,
		StackTrace: getStackTrace(),
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
