package connectkit

import (
	"errors"

	"connectrpc.com/connect"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func NewUnexpected() *connect.Error {
	err := connect.NewError(
		connect.CodeInternal,
		errors.New("An unexpected error occurred"))
	err.Meta().Set("x-error-code", "INTERNAL")
	return err
}

func NewUnauthorized() *connect.Error {
	err := connect.NewError(
		connect.CodePermissionDenied,
		errors.New("Unauthorized to run operation"))
	err.Meta().Set("x-error-code", "UNAUTHORIZED")
	return err
}

func NewInvalidArgument(subject string) *connect.Error {
	err := connect.NewError(
		connect.CodeInvalidArgument,
		errors.New("Invalid "+subject))
	err.Meta().Set("x-error-code", "INVALID_ARGUMENT")
	return err
}

func NewNotFound() *connect.Error {
	err := connect.NewError(
		connect.CodeNotFound,
		errors.New("Not found"))
	err.Meta().Set("x-error-code", "NOT_FOUND")
	return err
}

func NewResourceExhausted(reason string) *connect.Error {
	err := connect.NewError(
		connect.CodeUnauthenticated,
		errors.New(reason))
	if detail, detailErr := connect.NewErrorDetail(&errdetails.ErrorInfo{
		Reason: reason,
	}); detailErr == nil {
		err.AddDetail(detail)
	}
	err.Meta().Set("x-error-code", "RESOURCE_EXHAUSTED")
	err.Meta().Set("x-error-reason", reason)
	return err
}

func NewUnauthenticated(reason string) *connect.Error {
	err := connect.NewError(
		connect.CodeUnauthenticated,
		errors.New(reason))
	if detail, detailErr := connect.NewErrorDetail(&errdetails.ErrorInfo{
		Reason: reason,
	}); detailErr == nil {
		err.AddDetail(detail)
	}
	err.Meta().Set("x-error-code", "UNAUTHENTICATED")
	err.Meta().Set("x-error-reason", reason)
	return err
}

func NewConflict(reason string) *connect.Error {
	err := connect.NewError(
		connect.CodeAborted,
		errors.New(reason))
	if detail, detailErr := connect.NewErrorDetail(&errdetails.ErrorInfo{
		Reason: reason,
	}); detailErr == nil {
		err.AddDetail(detail)
	}
	err.Meta().Set("x-error-code", "CONFLICT")
	err.Meta().Set("x-error-reason", reason)
	return err
}
