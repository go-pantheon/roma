// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: player/v1/system.proto

package dbv1

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

// Validate checks the field values on UserSystemProto with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *UserSystemProto) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on UserSystemProto with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// UserSystemProtoMultiError, or nil if none found.
func (m *UserSystemProto) ValidateAll() error {
	return m.validate(true)
}

func (m *UserSystemProto) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for CurrentGenId

	for idx, item := range m.GetEvents() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, UserSystemProtoValidationError{
						field:  fmt.Sprintf("Events[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, UserSystemProtoValidationError{
						field:  fmt.Sprintf("Events[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return UserSystemProtoValidationError{
					field:  fmt.Sprintf("Events[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return UserSystemProtoMultiError(errors)
	}

	return nil
}

// UserSystemProtoMultiError is an error wrapping multiple validation errors
// returned by UserSystemProto.ValidateAll() if the designated constraints
// aren't met.
type UserSystemProtoMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UserSystemProtoMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UserSystemProtoMultiError) AllErrors() []error { return m }

// UserSystemProtoValidationError is the validation error returned by
// UserSystemProto.Validate if the designated constraints aren't met.
type UserSystemProtoValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UserSystemProtoValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UserSystemProtoValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UserSystemProtoValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UserSystemProtoValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UserSystemProtoValidationError) ErrorName() string { return "UserSystemProtoValidationError" }

// Error satisfies the builtin error interface
func (e UserSystemProtoValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUserSystemProto.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UserSystemProtoValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UserSystemProtoValidationError{}

// Validate checks the field values on WorkerEvent with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *WorkerEvent) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on WorkerEvent with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in WorkerEventMultiError, or
// nil if none found.
func (m *WorkerEvent) ValidateAll() error {
	return m.validate(true)
}

func (m *WorkerEvent) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	// no validation rules for Type

	if len(errors) > 0 {
		return WorkerEventMultiError(errors)
	}

	return nil
}

// WorkerEventMultiError is an error wrapping multiple validation errors
// returned by WorkerEvent.ValidateAll() if the designated constraints aren't met.
type WorkerEventMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m WorkerEventMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m WorkerEventMultiError) AllErrors() []error { return m }

// WorkerEventValidationError is the validation error returned by
// WorkerEvent.Validate if the designated constraints aren't met.
type WorkerEventValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e WorkerEventValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e WorkerEventValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e WorkerEventValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e WorkerEventValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e WorkerEventValidationError) ErrorName() string { return "WorkerEventValidationError" }

// Error satisfies the builtin error interface
func (e WorkerEventValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sWorkerEvent.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = WorkerEventValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = WorkerEventValidationError{}
