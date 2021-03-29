// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: message/hero_service.proto

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

// Validate checks the field values on SCPushHeroUnlock with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *SCPushHeroUnlock) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on SCPushHeroUnlock with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// SCPushHeroUnlockMultiError, or nil if none found.
func (m *SCPushHeroUnlock) ValidateAll() error {
	return m.validate(true)
}

func (m *SCPushHeroUnlock) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetHeroes() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, SCPushHeroUnlockValidationError{
						field:  fmt.Sprintf("Heroes[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, SCPushHeroUnlockValidationError{
						field:  fmt.Sprintf("Heroes[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return SCPushHeroUnlockValidationError{
					field:  fmt.Sprintf("Heroes[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return SCPushHeroUnlockMultiError(errors)
	}

	return nil
}

// SCPushHeroUnlockMultiError is an error wrapping multiple validation errors
// returned by SCPushHeroUnlock.ValidateAll() if the designated constraints
// aren't met.
type SCPushHeroUnlockMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m SCPushHeroUnlockMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m SCPushHeroUnlockMultiError) AllErrors() []error { return m }

// SCPushHeroUnlockValidationError is the validation error returned by
// SCPushHeroUnlock.Validate if the designated constraints aren't met.
type SCPushHeroUnlockValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e SCPushHeroUnlockValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e SCPushHeroUnlockValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e SCPushHeroUnlockValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e SCPushHeroUnlockValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e SCPushHeroUnlockValidationError) ErrorName() string { return "SCPushHeroUnlockValidationError" }

// Error satisfies the builtin error interface
func (e SCPushHeroUnlockValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sSCPushHeroUnlock.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = SCPushHeroUnlockValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = SCPushHeroUnlockValidationError{}

// Validate checks the field values on CSHeroLevelUpgrade with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *CSHeroLevelUpgrade) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CSHeroLevelUpgrade with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// CSHeroLevelUpgradeMultiError, or nil if none found.
func (m *CSHeroLevelUpgrade) ValidateAll() error {
	return m.validate(true)
}

func (m *CSHeroLevelUpgrade) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for HeroId

	if len(errors) > 0 {
		return CSHeroLevelUpgradeMultiError(errors)
	}

	return nil
}

// CSHeroLevelUpgradeMultiError is an error wrapping multiple validation errors
// returned by CSHeroLevelUpgrade.ValidateAll() if the designated constraints
// aren't met.
type CSHeroLevelUpgradeMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CSHeroLevelUpgradeMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CSHeroLevelUpgradeMultiError) AllErrors() []error { return m }

// CSHeroLevelUpgradeValidationError is the validation error returned by
// CSHeroLevelUpgrade.Validate if the designated constraints aren't met.
type CSHeroLevelUpgradeValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CSHeroLevelUpgradeValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CSHeroLevelUpgradeValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CSHeroLevelUpgradeValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CSHeroLevelUpgradeValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CSHeroLevelUpgradeValidationError) ErrorName() string {
	return "CSHeroLevelUpgradeValidationError"
}

// Error satisfies the builtin error interface
func (e CSHeroLevelUpgradeValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCSHeroLevelUpgrade.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CSHeroLevelUpgradeValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CSHeroLevelUpgradeValidationError{}

// Validate checks the field values on SCHeroLevelUpgrade with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *SCHeroLevelUpgrade) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on SCHeroLevelUpgrade with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// SCHeroLevelUpgradeMultiError, or nil if none found.
func (m *SCHeroLevelUpgrade) ValidateAll() error {
	return m.validate(true)
}

func (m *SCHeroLevelUpgrade) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Code

	if all {
		switch v := interface{}(m.GetHero()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, SCHeroLevelUpgradeValidationError{
					field:  "Hero",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, SCHeroLevelUpgradeValidationError{
					field:  "Hero",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetHero()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return SCHeroLevelUpgradeValidationError{
				field:  "Hero",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for Costs

	if len(errors) > 0 {
		return SCHeroLevelUpgradeMultiError(errors)
	}

	return nil
}

// SCHeroLevelUpgradeMultiError is an error wrapping multiple validation errors
// returned by SCHeroLevelUpgrade.ValidateAll() if the designated constraints
// aren't met.
type SCHeroLevelUpgradeMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m SCHeroLevelUpgradeMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m SCHeroLevelUpgradeMultiError) AllErrors() []error { return m }

// SCHeroLevelUpgradeValidationError is the validation error returned by
// SCHeroLevelUpgrade.Validate if the designated constraints aren't met.
type SCHeroLevelUpgradeValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e SCHeroLevelUpgradeValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e SCHeroLevelUpgradeValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e SCHeroLevelUpgradeValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e SCHeroLevelUpgradeValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e SCHeroLevelUpgradeValidationError) ErrorName() string {
	return "SCHeroLevelUpgradeValidationError"
}

// Error satisfies the builtin error interface
func (e SCHeroLevelUpgradeValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sSCHeroLevelUpgrade.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = SCHeroLevelUpgradeValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = SCHeroLevelUpgradeValidationError{}
