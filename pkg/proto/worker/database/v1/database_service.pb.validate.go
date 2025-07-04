// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: worker/database/v1/database_service.proto

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

// Validate checks the field values on Database with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Database) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Database with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in DatabaseMultiError, or nil
// if none found.
func (m *Database) ValidateAll() error {
	return m.validate(true)
}

func (m *Database) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Name

	if l := utf8.RuneCountInString(m.GetDisplayName()); l < 1 || l > 255 {
		err := DatabaseValidationError{
			field:  "DisplayName",
			reason: "value length must be between 1 and 255 runes, inclusive",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if all {
		switch v := interface{}(m.GetCreatedAt()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, DatabaseValidationError{
					field:  "CreatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, DatabaseValidationError{
					field:  "CreatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetCreatedAt()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return DatabaseValidationError{
				field:  "CreatedAt",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetUpdatedAt()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, DatabaseValidationError{
					field:  "UpdatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, DatabaseValidationError{
					field:  "UpdatedAt",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetUpdatedAt()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return DatabaseValidationError{
				field:  "UpdatedAt",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return DatabaseMultiError(errors)
	}

	return nil
}

// DatabaseMultiError is an error wrapping multiple validation errors returned
// by Database.ValidateAll() if the designated constraints aren't met.
type DatabaseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DatabaseMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DatabaseMultiError) AllErrors() []error { return m }

// DatabaseValidationError is the validation error returned by
// Database.Validate if the designated constraints aren't met.
type DatabaseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DatabaseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DatabaseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DatabaseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DatabaseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DatabaseValidationError) ErrorName() string { return "DatabaseValidationError" }

// Error satisfies the builtin error interface
func (e DatabaseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDatabase.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DatabaseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DatabaseValidationError{}

// Validate checks the field values on ListDatabasesRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *ListDatabasesRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ListDatabasesRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ListDatabasesRequestMultiError, or nil if none found.
func (m *ListDatabasesRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *ListDatabasesRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if val := m.GetPageSize(); val < 0 || val > 1000 {
		err := ListDatabasesRequestValidationError{
			field:  "PageSize",
			reason: "value must be inside range [0, 1000]",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	// no validation rules for PageToken

	if len(errors) > 0 {
		return ListDatabasesRequestMultiError(errors)
	}

	return nil
}

// ListDatabasesRequestMultiError is an error wrapping multiple validation
// errors returned by ListDatabasesRequest.ValidateAll() if the designated
// constraints aren't met.
type ListDatabasesRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ListDatabasesRequestMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ListDatabasesRequestMultiError) AllErrors() []error { return m }

// ListDatabasesRequestValidationError is the validation error returned by
// ListDatabasesRequest.Validate if the designated constraints aren't met.
type ListDatabasesRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListDatabasesRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListDatabasesRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListDatabasesRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListDatabasesRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListDatabasesRequestValidationError) ErrorName() string {
	return "ListDatabasesRequestValidationError"
}

// Error satisfies the builtin error interface
func (e ListDatabasesRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListDatabasesRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListDatabasesRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListDatabasesRequestValidationError{}

// Validate checks the field values on ListDatabasesResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *ListDatabasesResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ListDatabasesResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ListDatabasesResponseMultiError, or nil if none found.
func (m *ListDatabasesResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *ListDatabasesResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetDatabases() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, ListDatabasesResponseValidationError{
						field:  fmt.Sprintf("Databases[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, ListDatabasesResponseValidationError{
						field:  fmt.Sprintf("Databases[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ListDatabasesResponseValidationError{
					field:  fmt.Sprintf("Databases[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	// no validation rules for NextPageToken

	if len(errors) > 0 {
		return ListDatabasesResponseMultiError(errors)
	}

	return nil
}

// ListDatabasesResponseMultiError is an error wrapping multiple validation
// errors returned by ListDatabasesResponse.ValidateAll() if the designated
// constraints aren't met.
type ListDatabasesResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ListDatabasesResponseMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ListDatabasesResponseMultiError) AllErrors() []error { return m }

// ListDatabasesResponseValidationError is the validation error returned by
// ListDatabasesResponse.Validate if the designated constraints aren't met.
type ListDatabasesResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListDatabasesResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListDatabasesResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListDatabasesResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListDatabasesResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListDatabasesResponseValidationError) ErrorName() string {
	return "ListDatabasesResponseValidationError"
}

// Error satisfies the builtin error interface
func (e ListDatabasesResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListDatabasesResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListDatabasesResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListDatabasesResponseValidationError{}

// Validate checks the field values on GetDatabaseRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetDatabaseRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetDatabaseRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetDatabaseRequestMultiError, or nil if none found.
func (m *GetDatabaseRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *GetDatabaseRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if !_GetDatabaseRequest_Name_Pattern.MatchString(m.GetName()) {
		err := GetDatabaseRequestValidationError{
			field:  "Name",
			reason: "value does not match regex pattern \"^databases\\\\/[^\\\\/]+$\"",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return GetDatabaseRequestMultiError(errors)
	}

	return nil
}

// GetDatabaseRequestMultiError is an error wrapping multiple validation errors
// returned by GetDatabaseRequest.ValidateAll() if the designated constraints
// aren't met.
type GetDatabaseRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetDatabaseRequestMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetDatabaseRequestMultiError) AllErrors() []error { return m }

// GetDatabaseRequestValidationError is the validation error returned by
// GetDatabaseRequest.Validate if the designated constraints aren't met.
type GetDatabaseRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetDatabaseRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetDatabaseRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetDatabaseRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetDatabaseRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetDatabaseRequestValidationError) ErrorName() string {
	return "GetDatabaseRequestValidationError"
}

// Error satisfies the builtin error interface
func (e GetDatabaseRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetDatabaseRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetDatabaseRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetDatabaseRequestValidationError{}

var _GetDatabaseRequest_Name_Pattern = regexp.MustCompile("^databases\\/[^\\/]+$")

// Validate checks the field values on CreateDatabaseRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *CreateDatabaseRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CreateDatabaseRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// CreateDatabaseRequestMultiError, or nil if none found.
func (m *CreateDatabaseRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *CreateDatabaseRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetDatabaseId() != "" {

		if l := utf8.RuneCountInString(m.GetDatabaseId()); l < 1 || l > 64 {
			err := CreateDatabaseRequestValidationError{
				field:  "DatabaseId",
				reason: "value length must be between 1 and 64 runes, inclusive",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

		if !_CreateDatabaseRequest_DatabaseId_Pattern.MatchString(m.GetDatabaseId()) {
			err := CreateDatabaseRequestValidationError{
				field:  "DatabaseId",
				reason: "value does not match regex pattern \"^[a-z0-9\\\\-_.]*$\"",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if m.GetDatabase() == nil {
		err := CreateDatabaseRequestValidationError{
			field:  "Database",
			reason: "value is required",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if all {
		switch v := interface{}(m.GetDatabase()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, CreateDatabaseRequestValidationError{
					field:  "Database",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, CreateDatabaseRequestValidationError{
					field:  "Database",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetDatabase()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return CreateDatabaseRequestValidationError{
				field:  "Database",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return CreateDatabaseRequestMultiError(errors)
	}

	return nil
}

// CreateDatabaseRequestMultiError is an error wrapping multiple validation
// errors returned by CreateDatabaseRequest.ValidateAll() if the designated
// constraints aren't met.
type CreateDatabaseRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CreateDatabaseRequestMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CreateDatabaseRequestMultiError) AllErrors() []error { return m }

// CreateDatabaseRequestValidationError is the validation error returned by
// CreateDatabaseRequest.Validate if the designated constraints aren't met.
type CreateDatabaseRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreateDatabaseRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreateDatabaseRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreateDatabaseRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreateDatabaseRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreateDatabaseRequestValidationError) ErrorName() string {
	return "CreateDatabaseRequestValidationError"
}

// Error satisfies the builtin error interface
func (e CreateDatabaseRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreateDatabaseRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreateDatabaseRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreateDatabaseRequestValidationError{}

var _CreateDatabaseRequest_DatabaseId_Pattern = regexp.MustCompile("^[a-z0-9\\-_.]*$")

// Validate checks the field values on UpdateDatabaseRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *UpdateDatabaseRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on UpdateDatabaseRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// UpdateDatabaseRequestMultiError, or nil if none found.
func (m *UpdateDatabaseRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *UpdateDatabaseRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetDatabase() == nil {
		err := UpdateDatabaseRequestValidationError{
			field:  "Database",
			reason: "value is required",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if all {
		switch v := interface{}(m.GetDatabase()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, UpdateDatabaseRequestValidationError{
					field:  "Database",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, UpdateDatabaseRequestValidationError{
					field:  "Database",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetDatabase()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return UpdateDatabaseRequestValidationError{
				field:  "Database",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetUpdateMask()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, UpdateDatabaseRequestValidationError{
					field:  "UpdateMask",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, UpdateDatabaseRequestValidationError{
					field:  "UpdateMask",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetUpdateMask()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return UpdateDatabaseRequestValidationError{
				field:  "UpdateMask",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return UpdateDatabaseRequestMultiError(errors)
	}

	return nil
}

// UpdateDatabaseRequestMultiError is an error wrapping multiple validation
// errors returned by UpdateDatabaseRequest.ValidateAll() if the designated
// constraints aren't met.
type UpdateDatabaseRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UpdateDatabaseRequestMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UpdateDatabaseRequestMultiError) AllErrors() []error { return m }

// UpdateDatabaseRequestValidationError is the validation error returned by
// UpdateDatabaseRequest.Validate if the designated constraints aren't met.
type UpdateDatabaseRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UpdateDatabaseRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UpdateDatabaseRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UpdateDatabaseRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UpdateDatabaseRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UpdateDatabaseRequestValidationError) ErrorName() string {
	return "UpdateDatabaseRequestValidationError"
}

// Error satisfies the builtin error interface
func (e UpdateDatabaseRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUpdateDatabaseRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UpdateDatabaseRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UpdateDatabaseRequestValidationError{}

// Validate checks the field values on DeleteDatabaseRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *DeleteDatabaseRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on DeleteDatabaseRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// DeleteDatabaseRequestMultiError, or nil if none found.
func (m *DeleteDatabaseRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *DeleteDatabaseRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if !_DeleteDatabaseRequest_Name_Pattern.MatchString(m.GetName()) {
		err := DeleteDatabaseRequestValidationError{
			field:  "Name",
			reason: "value does not match regex pattern \"^databases\\\\/[^\\\\/]+$\"",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return DeleteDatabaseRequestMultiError(errors)
	}

	return nil
}

// DeleteDatabaseRequestMultiError is an error wrapping multiple validation
// errors returned by DeleteDatabaseRequest.ValidateAll() if the designated
// constraints aren't met.
type DeleteDatabaseRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DeleteDatabaseRequestMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DeleteDatabaseRequestMultiError) AllErrors() []error { return m }

// DeleteDatabaseRequestValidationError is the validation error returned by
// DeleteDatabaseRequest.Validate if the designated constraints aren't met.
type DeleteDatabaseRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DeleteDatabaseRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DeleteDatabaseRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DeleteDatabaseRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DeleteDatabaseRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DeleteDatabaseRequestValidationError) ErrorName() string {
	return "DeleteDatabaseRequestValidationError"
}

// Error satisfies the builtin error interface
func (e DeleteDatabaseRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDeleteDatabaseRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DeleteDatabaseRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DeleteDatabaseRequestValidationError{}

var _DeleteDatabaseRequest_Name_Pattern = regexp.MustCompile("^databases\\/[^\\/]+$")
