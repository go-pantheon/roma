// Code generated by protoc-gen-go-errors. DO NOT EDIT.

package adminv1

import (
	fmt "fmt"
	errors "github.com/go-kratos/kratos/v2/errors"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
const _ = errors.SupportPackageIsVersion1

func IsUserAdminErrorReasonUnspecified(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == UserAdminErrorReason_USER_ADMIN_ERROR_REASON_UNSPECIFIED.String() && e.Code == 500
}

func ErrorUserAdminErrorReasonUnspecified(format string, args ...interface{}) *errors.Error {
	return errors.New(500, UserAdminErrorReason_USER_ADMIN_ERROR_REASON_UNSPECIFIED.String(), fmt.Sprintf(format, args...))
}
