/**
 * Copyright (C) 2015 Red Hat, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package errors

import (
	"fmt"
	"io"
	"runtime/debug"
	"strings"

	"github.com/golang/glog"

	"github.com/openshift/origin/pkg/oc/util/prefixwriter"
)

type Error interface {
	error
	WithCause(error) Error
	WithSolution(string, ...interface{}) Error
	WithDetails(string) Error
}

func NewError(msg string, args ...interface{}) Error {
	return &internalError{
		msg: fmt.Sprintf(msg, args...),
	}
}

type internalError struct {
	msg      string
	cause    error
	solution string
	details  string
}

func (e *internalError) Error() string {
	return e.msg
}

func (e *internalError) Cause() error {
	return e.cause
}

func (e *internalError) Solution() string {
	return e.solution
}

func (e *internalError) Details() string {
	return e.details
}

func (e *internalError) WithCause(err error) Error {
	e.cause = err
	return e
}

func (e *internalError) WithDetails(details string) Error {
	e.details = details
	return e
}

func (e *internalError) WithSolution(solution string, args ...interface{}) Error {
	e.solution = fmt.Sprintf(solution, args...)
	return e
}

func LogError(err error) {
	if err == nil {
		return
	}
	glog.V(2).Infof("Unexpected error: %v", err)
	if glog.V(5) {
		debug.PrintStack()
	}
}

func PrintLog(out io.Writer, title string, content string) {
	fmt.Fprintf(out, "%s:\n", title)
	w := prefixwriter.New("  ", out)
	fmt.Fprintf(w, "%s", strings.TrimSpace(content))
	fmt.Fprintf(out, "\n")
}
