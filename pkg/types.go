package nerrors

import (
	"google.golang.org/grpc/codes"
)

type ErrorCode int

const (
	OK ErrorCode = iota
	Canceled
	Unknown
	InvalidArgument
	DeadlineExceeded
	NotFound
	AlreadyExists
	PermissionDenied
	ResourceExhausted
	FailedPrecondition
	Aborted
	OutOfRange
	Unimplemented
	Internal
	Unavailable
	DataLoss
	Unauthenticated
)

func (ec ErrorCode) String() string {
	return [...]string{"OK", "Canceled", "Unknown", "InvalidArgument", "DeadlineExceeded",
		"NotFound", "AlreadyExists", "PermissionDenied", "ResourceExhausted",
		"FailedPrecondition", "Aborted", "OutOfRange", "Unimplemented", "Internal",
		"Unavailable", "DataLoss", "Unauthenticated"}[ec]
}

func (ec ErrorCode) String2() string {
	switch ec {
	case OK:
		return "OK"
	case Canceled:
		return "Canceled"
	case Unknown:
		return "Unknown"
	case InvalidArgument:
		return "InvalidArgument"
	case DeadlineExceeded:
		return "DeadlineExceeded"
	case NotFound:
		return "NotFound"
	case AlreadyExists:
		return "AlreadyExists"
	case PermissionDenied:
		return "PermissionDenied"
	case ResourceExhausted:
		return "ResourceExhausted"
	case FailedPrecondition:
		return "FailedPrecondition"
	case Aborted:
		return "Aborted"
	case OutOfRange:
		return "OutOfRange"
	case Unimplemented:
		return "Unimplemented"
	case Internal:
		return "Internal"
	case Unavailable:
		return "Unavailable"
	case DataLoss:
		return "DataLoss"
	case Unauthenticated:
		return "Unauthenticated"
	default:
		return "Invalid Code"
	}
	return ""
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
