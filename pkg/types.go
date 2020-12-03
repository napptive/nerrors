package nerrors

import (
	"google.golang.org/grpc/codes"
)

type ErrorCode int

// https://github.com/grpc/grpc/blob/master/doc/statuscodes.md
// ErrorCode compatible with gRPC codes
const (
	OK                 ErrorCode = iota // Not an error
	Canceled                            // The operation was cancelled
	Unknown                             // Unknown error
	InvalidArgument                     // invalid argument
	DeadlineExceeded                    // The deadline expired before the operation could complete.
	NotFound                            // Some requested entity was not found.
	AlreadyExists                       // The entity to create already exists.
	PermissionDenied                    // The caller does not have permission to execute the specified operation
	ResourceExhausted                   // Some resource has been exhausted
	FailedPrecondition                  // The system is not in a state required for the operation's execution.
	Aborted                             // Aborted
	OutOfRange                          // The operation was attempted past the valid range.
	Unimplemented                       // The operation is not implemented or is not supported
	Internal                            // Internal error
	Unavailable                         // Currently unavailable
	DataLoss                            // Unrecoverable data loss or corruption
	Unauthenticated                     // The request does not have valid authentication credentials for the operation
)

func (ec ErrorCode) String() string {
	return [...]string{"OK", "Canceled", "Unknown", "InvalidArgument", "DeadlineExceeded",
		"NotFound", "AlreadyExists", "PermissionDenied", "ResourceExhausted",
		"FailedPrecondition", "Aborted", "OutOfRange", "Unimplemented", "Internal",
		"Unavailable", "DataLoss", "Unauthenticated"}[ec]
}

var ToGRPCCode = map[ErrorCode]codes.Code{
	OK:                 codes.OK,
	Canceled:           codes.Canceled,
	Unknown:            codes.Unknown,
	InvalidArgument:    codes.InvalidArgument,
	DeadlineExceeded:   codes.DeadlineExceeded,
	NotFound:           codes.NotFound,
	AlreadyExists:      codes.AlreadyExists,
	PermissionDenied:   codes.PermissionDenied,
	ResourceExhausted:  codes.ResourceExhausted,
	FailedPrecondition: codes.FailedPrecondition,
	Aborted:            codes.Aborted,
	OutOfRange:         codes.OutOfRange,
	Unimplemented:      codes.Unimplemented,
	Internal:           codes.Internal,
	Unavailable:        codes.Unavailable,
	DataLoss:           codes.DataLoss,
	Unauthenticated:    codes.Unauthenticated,
}

var FromGRPCCode = map[codes.Code]ErrorCode{
	codes.OK:                 OK,
	codes.Canceled:           Canceled,
	codes.Unknown:            Unknown,
	codes.InvalidArgument:    InvalidArgument,
	codes.DeadlineExceeded:   DeadlineExceeded,
	codes.NotFound:           NotFound,
	codes.AlreadyExists:      AlreadyExists,
	codes.PermissionDenied:   PermissionDenied,
	codes.ResourceExhausted:  ResourceExhausted,
	codes.FailedPrecondition: FailedPrecondition,
	codes.Aborted:            Aborted,
	codes.OutOfRange:         OutOfRange,
	codes.Unimplemented:      Unimplemented,
	codes.Internal:           Internal,
	codes.Unavailable:        Unavailable,
	codes.DataLoss:           DataLoss,
	codes.Unauthenticated:    Unauthenticated,
}
