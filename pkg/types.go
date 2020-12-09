package nerrors

import (
	"google.golang.org/grpc/codes"
)

type ErrorCode int

// https://github.com/grpc/grpc/blob/master/doc/statuscodes.md
// ErrorCode compatible with gRPC codes
const (
	// Ok indicates that this is not an error. This value is useful as enum use zero to represent the first element, and that information may not be sent through gRPC.
	OK                 ErrorCode = iota
	// Canceled indicates an operation was canceled and will no longer be executed.
	Canceled
	// Unknown indicates that while an error happened, its type is not known.
	Unknown
	// InvalidArgument indicates an error was detected due to the use of an invalid argument (e.g., invalid value).
	InvalidArgument
	// DeadlineExceeded indicates that the deadline expired before the operation could complete.
	DeadlineExceeded
	// NotFound indicates that some information about the requested entity was not found.
	NotFound
	// AlreadyExists indicates that an operation failed to create an entity because it already exists.
	AlreadyExists
	// PermissionDenied indicates that the caller does not have permission to execute the specified operation.
	PermissionDenied
	// ResourceExhausted indicates that the requested resource has been exhaused (e.g., no more pages exists).
	ResourceExhausted
	// FailedPrecondition indicates that an operation failed to satisfy a required precondition for its execution (e.g., base namespace was not created).
	FailedPrecondition
	// Aborted indicates that an operation was aborted due to some triggering factor.
	Aborted
	// OutOfRange indicates that an operation tried to access an element past the valid range.
	OutOfRange
	// Unimplemented indicates that the requested operation while declared is not implemented.
	Unimplemented
	// Internal indicates that an internal error in some component or system prevented the operation to be executed or failed during its execution.
	Internal
	// Unavailable indicates that the requested entity or operation is not available at this point.
	Unavailable
	// DataLoss indicates that an error related to unrecoverable data loss or corruption has happened.
	DataLoss
	// Unauthenticated indicates that the request does not have valid authentication credentials for the operation.
	Unauthenticated
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
