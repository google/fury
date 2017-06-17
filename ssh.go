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

package fury

import (
	"fmt"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/crypto/ssh/knownhosts"
)

// SSH runs commands over an SSH connection.
type SSH struct {
	client *ssh.Client
}

// NewRawSSH builds an SSH from an open ssh.Client connection.
func NewRawSSH(ssh *ssh.Client) *SSH {
	return &SSH{ssh}
}

// NewSSH builds an SSH connected to the given address.
//
// The connection is made as root, and authenticated with
// the resident SSH authentication agent.
func NewSSH(addr string) (*SSH, error) {
	authSock := os.Getenv("SSH_AUTH_SOCK")
	if authSock == "" {
		return nil, fmt.Errorf("no SSH agent found, SSH_AUTH_SOCK not defined")
	}
	authConn, err := net.Dial("unix", authSock)
	if err != nil {
		return nil, fmt.Errorf("dialing SSH agent: %s", err)
	}
	defer authConn.Close()
	agent := agent.NewClient(authConn)
	signers, err := agent.Signers()
	if err != nil {
		return nil, fmt.Errorf("getting signers from SSH agent: %s", err)
	}

	checker, err := knownHosts()
	if err != nil {
		return nil, fmt.Errorf("loading known hosts: %s", err)
	}

	client, err := ssh.Dial("tcp", addr, &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signers...),
		},
		HostKeyCallback: checker,
	})

	return &SSH{client}, err
}

// Close closes the session.
func (s *SSH) Close() error {
	return s.client.Close()
}

// Run executes cmd on the remote host.
func (s *SSH) Run(cmd *Command) error {
	sess, err := s.client.NewSession()
	if err != nil {
		return err
	}
	defer sess.Close()

	sess.Stdin = cmd.Stdin
	sess.Stdout = cmd.Stdout
	sess.Stderr = cmd.Stderr
	command, err := commandArgv(cmd)
	if err != nil {
		return fmt.Errorf("construct cmdline: %s", err)
	}

	return sess.Run(command)
}

// commandArgv turns cmd into a command string that can be provided to
// a shell for execution.
func commandArgv(cmd *Command) (string, error) {
	parts := []string{}

	if cmd.Dir != "" {
		parts = append(parts, "cd", shellEscape(cmd.Dir), "&&")
	}

	parts = append(parts, shellEscape("/usr/bin/env"), shellEscape("-"))

	var env []string
	for e := range cmd.Env {
		env = append(env, e)
	}
	sort.Strings(env)
	for _, e := range env {
		if err := validEnvVar(e); err != nil {
			return "", err
		}
		parts = append(parts, fmt.Sprintf("%s=%s", e, shellEscape(cmd.Env[e])))
	}

	parts = append(parts, shellEscape(cmd.Path))
	for _, arg := range cmd.Args {
		parts = append(parts, shellEscape(arg))
	}

	return strings.Join([]string{
		shellEscape("/bin/sh"),
		shellEscape("-c"),
		shellEscape(strings.Join(parts, " ")),
	}, " "), nil
}

// validEnvVar checks that e is a valid environment variable name.
func validEnvVar(e string) error {
	for i, c := range e {
		switch {
		case c >= 'a' && c <= 'z':
		case c >= 'A' && c <= 'Z':
		case c >= '0' && c <= '9':
			if i == 0 {
				return fmt.Errorf("environment variable %q starts with a digit", e)
			}
		case c == '_':
		default:
			return fmt.Errorf("illegal character %q in environment variable name %q", c, e)
		}
	}
	return nil
}

// shellEscape escapes a value such that /bin/sh will interpret it
// completely literally, with no expansions at all.
func shellEscape(val string) string {
	return fmt.Sprintf("'%s'", strings.Replace(val, "'", `'\''`, -1))
}

func knownHosts() (ssh.HostKeyCallback, error) {
	u, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("get current user: %s", err)
	}

	if u.HomeDir == "" {
		// No homedir, no known hosts.
		return nil, nil
	}

	path := filepath.Join(u.HomeDir, ".ssh/known_hosts")
	if _, err := os.Stat(path); err != nil && os.IsNotExist(err) {
		// No known_hosts file, no known hosts.
		return nil, nil
	}

	checker, err := knownhosts.New(path)
	if err != nil {
		return nil, fmt.Errorf("load %q: %s", path, err)
	}

	return checker, nil
}
