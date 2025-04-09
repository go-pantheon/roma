// Code generated by protoc-gen-go-errors. DO NOT EDIT.

package adminv1

import (
	fmt "fmt"
	errors "github.com/go-kratos/kratos/v2/errors"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
const _ = errors.SupportPackageIsVersion1

func IsStorageAdminErrorReasonUnspecified(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == StorageAdminErrorReason_STORAGE_ADMIN_ERROR_REASON_UNSPECIFIED.String() && e.Code == 500
}

func ErrorStorageAdminErrorReasonUnspecified(format string, args ...interface{}) *errors.Error {
	return errors.New(500, StorageAdminErrorReason_STORAGE_ADMIN_ERROR_REASON_UNSPECIFIED.String(), fmt.Sprintf(format, args...))
}

func IsStorageAdminErrorReasonServer(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == StorageAdminErrorReason_STORAGE_ADMIN_ERROR_REASON_SERVER.String() && e.Code == 500
}

func ErrorStorageAdminErrorReasonServer(format string, args ...interface{}) *errors.Error {
	return errors.New(500, StorageAdminErrorReason_STORAGE_ADMIN_ERROR_REASON_SERVER.String(), fmt.Sprintf(format, args...))
}

func IsStorageAdminErrorReasonUid(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == StorageAdminErrorReason_STORAGE_ADMIN_ERROR_REASON_UID.String() && e.Code == 401
}

func ErrorStorageAdminErrorReasonUid(format string, args ...interface{}) *errors.Error {
	return errors.New(401, StorageAdminErrorReason_STORAGE_ADMIN_ERROR_REASON_UID.String(), fmt.Sprintf(format, args...))
}
