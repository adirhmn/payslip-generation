package utils

import (
	"path/filepath"
	"runtime"
	"strings"
)

// GetCallerFunctionName returns the name of the function that called it
func GetCallerFunctionName(skip int) string {
	pc := make([]uintptr, 10)
	runtime.Callers(2+skip, pc)
	f := runtime.FuncForPC(pc[0])
	if f == nil {
		return ""
	}
	fullFuncName := f.Name()
	functionName := filepath.Ext(fullFuncName)
	return strings.TrimPrefix(functionName, ".")
}

// GetCallerBaseFunctionName returns the base name of the function that called it
func GetCallerBaseFunctionName(skip int) string {
	pc := make([]uintptr, 10)
	runtime.Callers(2+skip, pc)
	f := runtime.FuncForPC(pc[0])
	if f == nil {
		return ""
	}
	fullFuncName := f.Name()
	functionName := filepath.Base(fullFuncName)
	return functionName
}