// Code generated by protoc-gen-go-errors. DO NOT EDIT.

package servicev1

import (
	fmt "fmt"
	errors "github.com/go-kratos/kratos/errors"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
const _ = errors.SupportPackageIsVersion1

func IsRoomServiceErrorReasonUnspecified(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == RoomServiceErrorReason_ROOM_SERVICE_ERROR_REASON_UNSPECIFIED.String() && e.Code == 500
}

func ErrorRoomServiceErrorReasonUnspecified(format string, args ...interface{}) *errors.Error {
	return errors.New(500, RoomServiceErrorReason_ROOM_SERVICE_ERROR_REASON_UNSPECIFIED.String(), fmt.Sprintf(format, args...))
}

func IsRoomServiceErrorReasonServer(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == RoomServiceErrorReason_ROOM_SERVICE_ERROR_REASON_SERVER.String() && e.Code == 500
}

func ErrorRoomServiceErrorReasonServer(format string, args ...interface{}) *errors.Error {
	return errors.New(500, RoomServiceErrorReason_ROOM_SERVICE_ERROR_REASON_SERVER.String(), fmt.Sprintf(format, args...))
}

func IsRoomServiceErrorReasonId(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == RoomServiceErrorReason_ROOM_SERVICE_ERROR_REASON_ID.String() && e.Code == 401
}

func ErrorRoomServiceErrorReasonId(format string, args ...interface{}) *errors.Error {
	return errors.New(401, RoomServiceErrorReason_ROOM_SERVICE_ERROR_REASON_ID.String(), fmt.Sprintf(format, args...))
}
