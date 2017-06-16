// Copyright 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package fury implements a minimalist machine configuration
// management system.
package fury

import (
	"fmt"
	"io"
)

// A Command represents an external program to run.
type Command struct {
	// Path is the path of the command to run.
	//
	// This is the only field that must be set to a non-zero
	// value. If Path is relative, it is evaluated relative
	// to Dir.
	Path string

	// Args holds command line arguments.
	Args []string

	// Env specifies the environment of the process.
	Env map[string]string

	// Dir specifies the working directory of the command.
	// If Dir is the empty string, the working directory is /.
	Dir string

	// Stdin specifies the process's standard input.
	// If Stdin is nil, the process reads from the null device (os.DevNull).
	Stdin io.Reader

	// Stdout and Stderr specify the process's standard output and error.
	//
	// If either is nil, Run connects the corresponding file descriptor
	// to the null device (os.DevNull).
	//
	// If Stdout and Stderr are the same writer, at most one
	// goroutine at a time will call Write.
	Stdout io.Writer
	Stderr io.Writer
}

// ExitStatus stores information about an executed Command.
type ExitStatus struct {
	Status int
}

func (e ExitStatus) Error() string {
	return fmt.Sprintf("command exited with status %d", e.Status)
}

// A Runner executes external programs.
type Runner interface {
	// Run executes cmd.
	// The returned error is nil if the command runs, has no problems
	// copying stdin, stdout, and stderr, and exits with a zero exit
	// status.
	Run(cmd *Command) error
}

// An Applier applies its state to a Runner.
type Applier interface {
	Apply(Runner) error
}
