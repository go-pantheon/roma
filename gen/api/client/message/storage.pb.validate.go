// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: message/storage.proto

package climsg

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

// Validate checks the field values on UserStorageProto with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *UserStorageProto) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on UserStorageProto with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// UserStorageProtoMultiError, or nil if none found.
func (m *UserStorageProto) ValidateAll() error {
	return m.validate(true)
}

func (m *UserStorageProto) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Items

	// no validation rules for Packs

	{
		sorted_keys := make([]int64, len(m.GetRecoveryInfos()))
		i := 0
		for key := range m.GetRecoveryInfos() {
			sorted_keys[i] = key
			i++
		}
		sort.Slice(sorted_keys, func(i, j int) bool { return sorted_keys[i] < sorted_keys[j] })
		for _, key := range sorted_keys {
			val := m.GetRecoveryInfos()[key]
			_ = val

			// no validation rules for RecoveryInfos[key]

			if all {
				switch v := interface{}(val).(type) {
				case interface{ ValidateAll() error }:
					if err := v.ValidateAll(); err != nil {
						errors = append(errors, UserStorageProtoValidationError{
							field:  fmt.Sprintf("RecoveryInfos[%v]", key),
							reason: "embedded message failed validation",
							cause:  err,
						})
					}
				case interface{ Validate() error }:
					if err := v.Validate(); err != nil {
						errors = append(errors, UserStorageProtoValidationError{
							field:  fmt.Sprintf("RecoveryInfos[%v]", key),
							reason: "embedded message failed validation",
							cause:  err,
						})
					}
				}
			} else if v, ok := interface{}(val).(interface{ Validate() error }); ok {
				if err := v.Validate(); err != nil {
					return UserStorageProtoValidationError{
						field:  fmt.Sprintf("RecoveryInfos[%v]", key),
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}

		}
	}

	if len(errors) > 0 {
		return UserStorageProtoMultiError(errors)
	}

	return nil
}

// UserStorageProtoMultiError is an error wrapping multiple validation errors
// returned by UserStorageProto.ValidateAll() if the designated constraints
// aren't met.
type UserStorageProtoMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UserStorageProtoMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UserStorageProtoMultiError) AllErrors() []error { return m }

// UserStorageProtoValidationError is the validation error returned by
// UserStorageProto.Validate if the designated constraints aren't met.
type UserStorageProtoValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UserStorageProtoValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UserStorageProtoValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UserStorageProtoValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UserStorageProtoValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UserStorageProtoValidationError) ErrorName() string { return "UserStorageProtoValidationError" }

// Error satisfies the builtin error interface
func (e UserStorageProtoValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUserStorageProto.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UserStorageProtoValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UserStorageProtoValidationError{}

// Validate checks the field values on ItemRecoveryInfoProto with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *ItemRecoveryInfoProto) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ItemRecoveryInfoProto with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ItemRecoveryInfoProtoMultiError, or nil if none found.
func (m *ItemRecoveryInfoProto) ValidateAll() error {
	return m.validate(true)
}

func (m *ItemRecoveryInfoProto) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for DataId

	// no validation rules for Max

	// no validation rules for RecoverPerSec

	// no validation rules for UpdatedAt

	if len(errors) > 0 {
		return ItemRecoveryInfoProtoMultiError(errors)
	}

	return nil
}

// ItemRecoveryInfoProtoMultiError is an error wrapping multiple validation
// errors returned by ItemRecoveryInfoProto.ValidateAll() if the designated
// constraints aren't met.
type ItemRecoveryInfoProtoMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ItemRecoveryInfoProtoMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ItemRecoveryInfoProtoMultiError) AllErrors() []error { return m }

// ItemRecoveryInfoProtoValidationError is the validation error returned by
// ItemRecoveryInfoProto.Validate if the designated constraints aren't met.
type ItemRecoveryInfoProtoValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ItemRecoveryInfoProtoValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ItemRecoveryInfoProtoValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ItemRecoveryInfoProtoValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ItemRecoveryInfoProtoValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ItemRecoveryInfoProtoValidationError) ErrorName() string {
	return "ItemRecoveryInfoProtoValidationError"
}

// Error satisfies the builtin error interface
func (e ItemRecoveryInfoProtoValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sItemRecoveryInfoProto.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ItemRecoveryInfoProtoValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ItemRecoveryInfoProtoValidationError{}
