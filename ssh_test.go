package fury

import (
	"os/exec"
	"testing"
)

func TestShellEscape(t *testing.T) {
	cases := map[string]string{
		`foo`:    `'foo'`,
		`f'tang`: `'f'\''tang'`,
		`'bar'`:  `''\''bar'\'''`,
	}

	for in, want := range cases {
		got := shellEscape(in)
		if got != want {
			t.Errorf("incorrect escaping\n  got : %s\n  want: %s", got, want)
		}
	}
}

func TestCommandArgv(t *testing.T) {
	cases := map[string]Command{
		"foo bar r'lyeh\n": Command{
			Path: "/usr/bin/echo",
			Args: []string{"foo", "bar", "r'lyeh"},
		},
		"testing\n": Command{
			Path: "/usr/bin/printenv",
			Args: []string{"UNIVERSE"},
			Env: map[string]string{
				"UNIVERSE": "testing",
			},
		},
		"/\n": Command{
			Path: "/usr/bin/pwd",
			Dir:  "/",
		},
		"/tmp\nr'lyeh\n": Command{
			Path: "/bin/sh",
			Args: []string{"-c", "pwd && printenv UNIVERSE"},
			Env: map[string]string{
				"UNIVERSE": "r'lyeh",
			},
			Dir: "/tmp",
		},
	}

	for want, in := range cases {
		argv, err := commandArgv(&in)
		if err != nil {
			t.Fatalf("constructing commandline: %s", err)
		}
		out, err := exec.Command("/bin/sh", "-c", argv).CombinedOutput()
		if err != nil {
			t.Fatalf("executing command: %s (output %q)", err, string(out))
		}
		got := string(out)
		if got != want {
			t.Errorf("unexpected commandArgv result, got %q, want %q", got, want)
		}
	}
}

// TODO: run an ssh server and test what shell commands we receive.
