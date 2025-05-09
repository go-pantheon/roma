// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: player/v1/hero.proto

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

// Validate checks the field values on UserHeroListProto with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *UserHeroListProto) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on UserHeroListProto with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// UserHeroListProtoMultiError, or nil if none found.
func (m *UserHeroListProto) ValidateAll() error {
	return m.validate(true)
}

func (m *UserHeroListProto) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	{
		sorted_keys := make([]int64, len(m.GetHeroes()))
		i := 0
		for key := range m.GetHeroes() {
			sorted_keys[i] = key
			i++
		}
		sort.Slice(sorted_keys, func(i, j int) bool { return sorted_keys[i] < sorted_keys[j] })
		for _, key := range sorted_keys {
			val := m.GetHeroes()[key]
			_ = val

			// no validation rules for Heroes[key]

			if all {
				switch v := interface{}(val).(type) {
				case interface{ ValidateAll() error }:
					if err := v.ValidateAll(); err != nil {
						errors = append(errors, UserHeroListProtoValidationError{
							field:  fmt.Sprintf("Heroes[%v]", key),
							reason: "embedded message failed validation",
							cause:  err,
						})
					}
				case interface{ Validate() error }:
					if err := v.Validate(); err != nil {
						errors = append(errors, UserHeroListProtoValidationError{
							field:  fmt.Sprintf("Heroes[%v]", key),
							reason: "embedded message failed validation",
							cause:  err,
						})
					}
				}
			} else if v, ok := interface{}(val).(interface{ Validate() error }); ok {
				if err := v.Validate(); err != nil {
					return UserHeroListProtoValidationError{
						field:  fmt.Sprintf("Heroes[%v]", key),
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}

		}
	}

	if len(errors) > 0 {
		return UserHeroListProtoMultiError(errors)
	}

	return nil
}

// UserHeroListProtoMultiError is an error wrapping multiple validation errors
// returned by UserHeroListProto.ValidateAll() if the designated constraints
// aren't met.
type UserHeroListProtoMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UserHeroListProtoMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UserHeroListProtoMultiError) AllErrors() []error { return m }

// UserHeroListProtoValidationError is the validation error returned by
// UserHeroListProto.Validate if the designated constraints aren't met.
type UserHeroListProtoValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UserHeroListProtoValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UserHeroListProtoValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UserHeroListProtoValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UserHeroListProtoValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UserHeroListProtoValidationError) ErrorName() string {
	return "UserHeroListProtoValidationError"
}

// Error satisfies the builtin error interface
func (e UserHeroListProtoValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUserHeroListProto.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UserHeroListProtoValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UserHeroListProtoValidationError{}

// Validate checks the field values on UserHeroProto with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *UserHeroProto) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on UserHeroProto with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in UserHeroProtoMultiError, or
// nil if none found.
func (m *UserHeroProto) ValidateAll() error {
	return m.validate(true)
}

func (m *UserHeroProto) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	// no validation rules for DataId

	// no validation rules for Level

	// no validation rules for Skills

	if len(errors) > 0 {
		return UserHeroProtoMultiError(errors)
	}

	return nil
}

// UserHeroProtoMultiError is an error wrapping multiple validation errors
// returned by UserHeroProto.ValidateAll() if the designated constraints
// aren't met.
type UserHeroProtoMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UserHeroProtoMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UserHeroProtoMultiError) AllErrors() []error { return m }

// UserHeroProtoValidationError is the validation error returned by
// UserHeroProto.Validate if the designated constraints aren't met.
type UserHeroProtoValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UserHeroProtoValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UserHeroProtoValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UserHeroProtoValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UserHeroProtoValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UserHeroProtoValidationError) ErrorName() string { return "UserHeroProtoValidationError" }

// Error satisfies the builtin error interface
func (e UserHeroProtoValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUserHeroProto.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UserHeroProtoValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UserHeroProtoValidationError{}
