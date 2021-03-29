// Code generated by protoc-gen-go-errors. DO NOT EDIT.

package adminv1

import (
	fmt "fmt"
	errors "github.com/go-kratos/kratos/errors"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
const _ = errors.SupportPackageIsVersion1

func IsGuildAdminErrorReasonUnspecified(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == GuildAdminErrorReason_GUILD_ADMIN_ERROR_REASON_UNSPECIFIED.String() && e.Code == 500
}

func ErrorGuildAdminErrorReasonUnspecified(format string, args ...interface{}) *errors.Error {
	return errors.New(500, GuildAdminErrorReason_GUILD_ADMIN_ERROR_REASON_UNSPECIFIED.String(), fmt.Sprintf(format, args...))
}

func IsGuildAdminErrorReasonServer(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == GuildAdminErrorReason_GUILD_ADMIN_ERROR_REASON_SERVER.String() && e.Code == 500
}

func ErrorGuildAdminErrorReasonServer(format string, args ...interface{}) *errors.Error {
	return errors.New(500, GuildAdminErrorReason_GUILD_ADMIN_ERROR_REASON_SERVER.String(), fmt.Sprintf(format, args...))
}

func IsGuildAdminErrorReasonServerId(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == GuildAdminErrorReason_GUILD_ADMIN_ERROR_REASON_SERVER_ID.String() && e.Code == 401
}

func ErrorGuildAdminErrorReasonServerId(format string, args ...interface{}) *errors.Error {
	return errors.New(401, GuildAdminErrorReason_GUILD_ADMIN_ERROR_REASON_SERVER_ID.String(), fmt.Sprintf(format, args...))
}
