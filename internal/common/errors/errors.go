package errors

import (
	"errors"
	"reflect"
	"slices"

	"payslip-generation-system/internal/utils"

	jsoniter "github.com/json-iterator/go"
)

type Error interface {
	Error() string
}

type err struct {
	Title  string   `json:"title"`
	Detail string   `json:"detail,omitempty"`
	Traces []string `json:"traces,omitempty"`
	ErrVal error    `json:"err_val,omitempty"`
}

// New creates a new (custom) error with the given arguments
func New(arg interface{}) *err {
	var newErr err

	switch v := arg.(type) {
	case string:
		newErr = err{
			Title: v,
			Traces: []string{
				utils.GetFileAndLoC(1),
			},
		}
	case err:
		errVal := arg.(err)
		errVal.Traces = append(
			[]string{utils.GetFileAndLoC(1)},
			errVal.Traces...,
		)
		newErr = errVal
	case *err:
		errVal := arg.(*err)
		errVal.Traces = append(
			[]string{utils.GetFileAndLoC(1)},
			errVal.Traces...,
		)
		newErr = *errVal
	case error:
		newErr = err{
			Title: v.Error(),
			Traces: []string{
				utils.GetFileAndLoC(1),
			},
			ErrVal: v,
		}
	}

	return &newErr
}

// Error returns the error as a string
func (e *err) Error() string {
	// reversing the traces to show the latest trace first
	eCopy := *e // copying the error to avoid modifying the original error (i.e: upon calling Error for logs)
	slices.Reverse(eCopy.Traces)
	err, _ := jsoniter.MarshalToString(eCopy)
	return err
}


// Is checks if the error is equal to the error being compared
func (e *err) IsEqual(errCompared error) bool {
	if reflect.TypeOf(errCompared) == reflect.TypeOf(e) {
		errCast := errCompared.(*err)
		return errors.Is(e.ErrVal, errCast.ErrVal)
	}

	if e.ErrVal != nil {
		return errors.Is(e.ErrVal, errCompared)
	}

	return errors.Is(errCompared, e)
}

// Is checks if the two errors are equal
func Is(err1, err2 error) bool {
	if err1 == nil || err2 == nil {
		return err1 == err2
	}

	err1Cast, ok1 := err1.(*err)
	if ok1 {
		return err1Cast.IsEqual(err2)
	}

	err2Cast, ok2 := err2.(*err)
	if ok2 {
		return err2Cast.IsEqual(err1)
	}

	return errors.Is(err1, err2)
}