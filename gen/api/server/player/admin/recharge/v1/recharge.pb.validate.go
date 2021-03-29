// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: player/admin/recharge/v1/recharge.proto

package adminv1

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on GetOrderListRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetOrderListRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetOrderListRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetOrderListRequestMultiError, or nil if none found.
func (m *GetOrderListRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *GetOrderListRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Page

	// no validation rules for PageSize

	if all {
		switch v := interface{}(m.GetCond()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, GetOrderListRequestValidationError{
					field:  "Cond",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, GetOrderListRequestValidationError{
					field:  "Cond",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetCond()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return GetOrderListRequestValidationError{
				field:  "Cond",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return GetOrderListRequestMultiError(errors)
	}

	return nil
}

// GetOrderListRequestMultiError is an error wrapping multiple validation
// errors returned by GetOrderListRequest.ValidateAll() if the designated
// constraints aren't met.
type GetOrderListRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetOrderListRequestMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetOrderListRequestMultiError) AllErrors() []error { return m }

// GetOrderListRequestValidationError is the validation error returned by
// GetOrderListRequest.Validate if the designated constraints aren't met.
type GetOrderListRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetOrderListRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetOrderListRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetOrderListRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetOrderListRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetOrderListRequestValidationError) ErrorName() string {
	return "GetOrderListRequestValidationError"
}

// Error satisfies the builtin error interface
func (e GetOrderListRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetOrderListRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetOrderListRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetOrderListRequestValidationError{}

// Validate checks the field values on GetOrderListResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetOrderListResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetOrderListResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetOrderListResponseMultiError, or nil if none found.
func (m *GetOrderListResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *GetOrderListResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Code

	// no validation rules for Message

	for idx, item := range m.GetOrders() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, GetOrderListResponseValidationError{
						field:  fmt.Sprintf("Orders[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, GetOrderListResponseValidationError{
						field:  fmt.Sprintf("Orders[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return GetOrderListResponseValidationError{
					field:  fmt.Sprintf("Orders[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	// no validation rules for Total

	if len(errors) > 0 {
		return GetOrderListResponseMultiError(errors)
	}

	return nil
}

// GetOrderListResponseMultiError is an error wrapping multiple validation
// errors returned by GetOrderListResponse.ValidateAll() if the designated
// constraints aren't met.
type GetOrderListResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetOrderListResponseMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetOrderListResponseMultiError) AllErrors() []error { return m }

// GetOrderListResponseValidationError is the validation error returned by
// GetOrderListResponse.Validate if the designated constraints aren't met.
type GetOrderListResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetOrderListResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetOrderListResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetOrderListResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetOrderListResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetOrderListResponseValidationError) ErrorName() string {
	return "GetOrderListResponseValidationError"
}

// Error satisfies the builtin error interface
func (e GetOrderListResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetOrderListResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetOrderListResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetOrderListResponseValidationError{}

// Validate checks the field values on GetOrderListCond with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *GetOrderListCond) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetOrderListCond with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetOrderListCondMultiError, or nil if none found.
func (m *GetOrderListCond) ValidateAll() error {
	return m.validate(true)
}

func (m *GetOrderListCond) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Store

	// no validation rules for TransId

	// no validation rules for Uid

	// no validation rules for Ack

	if len(errors) > 0 {
		return GetOrderListCondMultiError(errors)
	}

	return nil
}

// GetOrderListCondMultiError is an error wrapping multiple validation errors
// returned by GetOrderListCond.ValidateAll() if the designated constraints
// aren't met.
type GetOrderListCondMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetOrderListCondMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetOrderListCondMultiError) AllErrors() []error { return m }

// GetOrderListCondValidationError is the validation error returned by
// GetOrderListCond.Validate if the designated constraints aren't met.
type GetOrderListCondValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetOrderListCondValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetOrderListCondValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetOrderListCondValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetOrderListCondValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetOrderListCondValidationError) ErrorName() string { return "GetOrderListCondValidationError" }

// Error satisfies the builtin error interface
func (e GetOrderListCondValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetOrderListCond.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetOrderListCondValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetOrderListCondValidationError{}

// Validate checks the field values on GetOrderByIdRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetOrderByIdRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetOrderByIdRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetOrderByIdRequestMultiError, or nil if none found.
func (m *GetOrderByIdRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *GetOrderByIdRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Store

	// no validation rules for TransId

	if len(errors) > 0 {
		return GetOrderByIdRequestMultiError(errors)
	}

	return nil
}

// GetOrderByIdRequestMultiError is an error wrapping multiple validation
// errors returned by GetOrderByIdRequest.ValidateAll() if the designated
// constraints aren't met.
type GetOrderByIdRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetOrderByIdRequestMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetOrderByIdRequestMultiError) AllErrors() []error { return m }

// GetOrderByIdRequestValidationError is the validation error returned by
// GetOrderByIdRequest.Validate if the designated constraints aren't met.
type GetOrderByIdRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetOrderByIdRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetOrderByIdRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetOrderByIdRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetOrderByIdRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetOrderByIdRequestValidationError) ErrorName() string {
	return "GetOrderByIdRequestValidationError"
}

// Error satisfies the builtin error interface
func (e GetOrderByIdRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetOrderByIdRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetOrderByIdRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetOrderByIdRequestValidationError{}

// Validate checks the field values on GetOrderByIdResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetOrderByIdResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetOrderByIdResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetOrderByIdResponseMultiError, or nil if none found.
func (m *GetOrderByIdResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *GetOrderByIdResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Code

	// no validation rules for Message

	if all {
		switch v := interface{}(m.GetOrder()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, GetOrderByIdResponseValidationError{
					field:  "Order",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, GetOrderByIdResponseValidationError{
					field:  "Order",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetOrder()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return GetOrderByIdResponseValidationError{
				field:  "Order",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return GetOrderByIdResponseMultiError(errors)
	}

	return nil
}

// GetOrderByIdResponseMultiError is an error wrapping multiple validation
// errors returned by GetOrderByIdResponse.ValidateAll() if the designated
// constraints aren't met.
type GetOrderByIdResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetOrderByIdResponseMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetOrderByIdResponseMultiError) AllErrors() []error { return m }

// GetOrderByIdResponseValidationError is the validation error returned by
// GetOrderByIdResponse.Validate if the designated constraints aren't met.
type GetOrderByIdResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetOrderByIdResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetOrderByIdResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetOrderByIdResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetOrderByIdResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetOrderByIdResponseValidationError) ErrorName() string {
	return "GetOrderByIdResponseValidationError"
}

// Error satisfies the builtin error interface
func (e GetOrderByIdResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetOrderByIdResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetOrderByIdResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetOrderByIdResponseValidationError{}

// Validate checks the field values on UpdateOrderAckStateByIdRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *UpdateOrderAckStateByIdRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on UpdateOrderAckStateByIdRequest with
// the rules defined in the proto definition for this message. If any rules
// are violated, the result is a list of violation errors wrapped in
// UpdateOrderAckStateByIdRequestMultiError, or nil if none found.
func (m *UpdateOrderAckStateByIdRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *UpdateOrderAckStateByIdRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Store

	// no validation rules for TransId

	// no validation rules for Ack

	if len(errors) > 0 {
		return UpdateOrderAckStateByIdRequestMultiError(errors)
	}

	return nil
}

// UpdateOrderAckStateByIdRequestMultiError is an error wrapping multiple
// validation errors returned by UpdateOrderAckStateByIdRequest.ValidateAll()
// if the designated constraints aren't met.
type UpdateOrderAckStateByIdRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UpdateOrderAckStateByIdRequestMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UpdateOrderAckStateByIdRequestMultiError) AllErrors() []error { return m }

// UpdateOrderAckStateByIdRequestValidationError is the validation error
// returned by UpdateOrderAckStateByIdRequest.Validate if the designated
// constraints aren't met.
type UpdateOrderAckStateByIdRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UpdateOrderAckStateByIdRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UpdateOrderAckStateByIdRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UpdateOrderAckStateByIdRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UpdateOrderAckStateByIdRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UpdateOrderAckStateByIdRequestValidationError) ErrorName() string {
	return "UpdateOrderAckStateByIdRequestValidationError"
}

// Error satisfies the builtin error interface
func (e UpdateOrderAckStateByIdRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUpdateOrderAckStateByIdRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UpdateOrderAckStateByIdRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UpdateOrderAckStateByIdRequestValidationError{}

// Validate checks the field values on UpdateOrderAckStateByIdResponse with the
// rules defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *UpdateOrderAckStateByIdResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on UpdateOrderAckStateByIdResponse with
// the rules defined in the proto definition for this message. If any rules
// are violated, the result is a list of violation errors wrapped in
// UpdateOrderAckStateByIdResponseMultiError, or nil if none found.
func (m *UpdateOrderAckStateByIdResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *UpdateOrderAckStateByIdResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Code

	// no validation rules for Message

	if len(errors) > 0 {
		return UpdateOrderAckStateByIdResponseMultiError(errors)
	}

	return nil
}

// UpdateOrderAckStateByIdResponseMultiError is an error wrapping multiple
// validation errors returned by UpdateOrderAckStateByIdResponse.ValidateAll()
// if the designated constraints aren't met.
type UpdateOrderAckStateByIdResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UpdateOrderAckStateByIdResponseMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UpdateOrderAckStateByIdResponseMultiError) AllErrors() []error { return m }

// UpdateOrderAckStateByIdResponseValidationError is the validation error
// returned by UpdateOrderAckStateByIdResponse.Validate if the designated
// constraints aren't met.
type UpdateOrderAckStateByIdResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UpdateOrderAckStateByIdResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UpdateOrderAckStateByIdResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UpdateOrderAckStateByIdResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UpdateOrderAckStateByIdResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UpdateOrderAckStateByIdResponseValidationError) ErrorName() string {
	return "UpdateOrderAckStateByIdResponseValidationError"
}

// Error satisfies the builtin error interface
func (e UpdateOrderAckStateByIdResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUpdateOrderAckStateByIdResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UpdateOrderAckStateByIdResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UpdateOrderAckStateByIdResponseValidationError{}

// Validate checks the field values on OrderProto with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *OrderProto) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on OrderProto with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in OrderProtoMultiError, or
// nil if none found.
func (m *OrderProto) ValidateAll() error {
	return m.validate(true)
}

func (m *OrderProto) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Store

	// no validation rules for TransId

	// no validation rules for Ack

	// no validation rules for Uid

	// no validation rules for ProductId

	// no validation rules for PurchasedAt

	// no validation rules for AckAt

	// no validation rules for Detail

	if len(errors) > 0 {
		return OrderProtoMultiError(errors)
	}

	return nil
}

// OrderProtoMultiError is an error wrapping multiple validation errors
// returned by OrderProto.ValidateAll() if the designated constraints aren't met.
type OrderProtoMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m OrderProtoMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m OrderProtoMultiError) AllErrors() []error { return m }

// OrderProtoValidationError is the validation error returned by
// OrderProto.Validate if the designated constraints aren't met.
type OrderProtoValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e OrderProtoValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e OrderProtoValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e OrderProtoValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e OrderProtoValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e OrderProtoValidationError) ErrorName() string { return "OrderProtoValidationError" }

// Error satisfies the builtin error interface
func (e OrderProtoValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sOrderProto.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = OrderProtoValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = OrderProtoValidationError{}
