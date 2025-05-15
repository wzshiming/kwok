/*
Copyright 2025 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package server

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"golang.org/x/term"
	remotecommandclient "k8s.io/client-go/tools/remotecommand"
	"sigs.k8s.io/kwok/pkg/apis/internalversion"

	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

func runAll(ctx context.Context, conf *internalversion.ExecSandbox, cmd []string, in io.Reader, out, errOut io.Writer, tty bool, resize <-chan remotecommandclient.TerminalSize) error {

	if tty {
		stream := struct {
			io.Reader
			io.Writer
		}{
			in,
			out,
		}

		t := term.NewTerminal(stream, "")

		errOut = t
		out = t

		pr, pw := io.Pipe()
		in = pr
		go func() {
			for {
				line, err := t.ReadLine()
				if err != nil {
					return
				}
				pw.Write([]byte(line))
				pw.Write([]byte("\n"))
			}
		}()
		if resize != nil {
			go func() {
				for s := range resize {
					_ = t.SetSize(int(s.Width), int(s.Height))
				}
			}()
		}
	}

	return runInteractive(ctx, in, out, errOut)
}

var testBuiltinsMap = map[string]func(interp.HandlerContext, []string) error{
	"sleep": func(hc interp.HandlerContext, args []string) error {
		if len(args) < 1 {
			_, _ = io.WriteString(hc.Stderr, "sleep: missing operand\r\n")
			return nil
		}

		for _, arg := range args {
			duration, err := time.ParseDuration(arg)
			if err != nil {
				i, err := strconv.ParseInt(arg, 0, 0)
				if err != nil {
					_, _ = io.WriteString(hc.Stderr, fmt.Sprintf("sleep: invalid time interval '%s'\r\n", arg))
					return nil
				}
				duration = time.Duration(i) * time.Second
			}

			time.Sleep(duration)
		}
		return nil
	},
	"date": func(hc interp.HandlerContext, s []string) error {
		_, _ = io.WriteString(hc.Stdout, time.Now().UTC().Format(time.UnixDate)+"\r\n")
		return nil
	},
}

func testExecHandler(next interp.ExecHandlerFunc) interp.ExecHandlerFunc {
	return func(ctx context.Context, args []string) error {
		if fn := testBuiltinsMap[args[0]]; fn != nil {
			return fn(interp.HandlerCtx(ctx), args[1:])
		}
		return nil
		// return next(ctx, args)
	}
}

func runInteractive(ctx context.Context, in io.Reader, out, errOut io.Writer) error {
	r, err := interp.New(
		interp.StdIO(nil, out, errOut),
		interp.Interactive(true),
		interp.ExecHandlers(testExecHandler),
		interp.OpenHandler(func(ctx context.Context, path string, flag int, perm os.FileMode) (io.ReadWriteCloser, error) {
			return nil, os.ErrPermission
		}),
		interp.ReadDirHandler2(func(ctx context.Context, path string) ([]fs.DirEntry, error) {
			return nil, os.ErrPermission
		}),
		interp.StatHandler(func(ctx context.Context, name string, followSymlinks bool) (fs.FileInfo, error) {
			return nil, os.ErrPermission
		}),
	)
	if err != nil {
		return err
	}

	parser := syntax.NewParser()
	fmt.Fprintf(out, "$ ")
	var runErr error
	fn := func(stmts []*syntax.Stmt) bool {
		if parser.Incomplete() {
			fmt.Fprintf(out, "> ")
			return true
		}
		for _, stmt := range stmts {
			runErr = r.Run(ctx, stmt)
			if r.Exited() {
				return false
			}
		}
		fmt.Fprintf(out, "$ ")
		return true
	}

	if err := parser.Interactive(in, fn); err != nil {
		return err
	}
	return runErr
}

// execInSandbox simulates container environments for Kubernetes exec debugging.
// designed for development/testing scenarios requiring predictable container behavior.
func execInSandbox(ctx context.Context, conf *internalversion.ExecSandbox, cmd []string, in io.Reader, out, errOut io.Writer, tty bool, resize <-chan remotecommandclient.TerminalSize) error {
	if resize != nil {
		go func() {
			for {
				select {
				case <-resize:
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	if tty {
		errOut = out
	}

	if cmd[0] != "sh" && cmd[0] != "/bin/sh" {
		return execInSandboxSingleCommand(ctx, conf, cmd, in, out, errOut)
	}

	if in == nil {
		return nil
	}

	if !tty {
		// TODO: Implement non-TTY (non-interactive) shell command support
		return fmt.Errorf("unsupport yet")
	}

	t := newSandboxTerminal(in, out, conf)

	return t.Run(ctx)
}

var (
	sandboxCommand = map[string]func(t *sandboxContext, args []string) bool{
		"sleep": func(t *sandboxContext, args []string) bool {
			if len(args) < 1 {
				_, _ = io.WriteString(t.errOut, "sleep: missing operand\r\n")
				return true
			}
			duration, err := time.ParseDuration(args[0])
			if err != nil {
				i, err := strconv.ParseInt(args[0], 0, 0)
				if err != nil {
					_, _ = io.WriteString(t.errOut, fmt.Sprintf("sleep: invalid time interval '%s'\r\n", args[0]))
					return true
				}
				duration = time.Duration(i) * time.Second
			}
			time.Sleep(duration)
			return true
		},
		"pwd": func(t *sandboxContext, args []string) bool {
			_, _ = io.WriteString(t.out, t.currentDir+"\r\n")
			return true
		},
		"echo": func(t *sandboxContext, args []string) bool {
			_, _ = io.WriteString(t.out, strings.Join(args, " ")+"\r\n")
			return true
		},
		"exit": func(t *sandboxContext, args []string) bool {
			_, _ = io.WriteString(t.errOut, "exit\r\n")
			if len(t.stack) > 0 {
				*t = t.stack[len(t.stack)-1]
				return true
			}
			return false
		},
		"ls": func(t *sandboxContext, args []string) bool {
			var dir string
			if len(args) > 0 {
				dir = args[0]
			} else {
				dir = t.currentDir
			}

			targetDir := path.Clean(dir)
			contents, exists := sandboxDirTree[targetDir]
			if !exists {
				return true
			}

			_, _ = io.WriteString(t.out, strings.Join(contents, "  ")+"\r\n")
			return true
		},
		"cd": func(t *sandboxContext, args []string) bool {
			var dir string
			if len(args) > 0 {
				dir = args[0]
			} else {
				dir = "/root"
			}

			var targetDir string
			if strings.HasPrefix(dir, "/") {
				targetDir = path.Clean(dir)
			} else {
				targetDir = path.Clean(path.Join(t.currentDir, dir))
			}

			if _, exists := sandboxDirTree[targetDir]; exists {
				t.currentDir = targetDir
				return true
			}
			_, _ = io.WriteString(t.errOut, fmt.Sprintf("cd: no such file or directory: %s\r\n", dir))
			return true
		},

		"date": func(t *sandboxContext, args []string) bool {
			_, _ = io.WriteString(t.out, time.Now().UTC().Format(time.UnixDate)+"\r\n")
			return true
		},
		"clear": func(t *sandboxContext, args []string) bool {
			_, _ = io.WriteString(t.errOut, "\033[H\033[2J")
			return true
		},
		"sh": func(t *sandboxContext, args []string) bool {
			t.stack = append(t.stack, *t)
			return true
		},
	}
	sandboxDirTree = map[string][]string{
		"/":      {"bin", "dev", "etc", "home", "lib", "lib64", "proc", "root", "sys", "tmp", "usr", "var"},
		"/bin":   {"clear", "date", "sleep", "sh"},
		"/dev":   {},
		"/etc":   {},
		"/home":  {},
		"/lib":   {},
		"/lib64": {},
		"/proc":  {},
		"/root":  {},
		"/sys":   {},
		"/sbin":  {},
		"/tmp":   {},
		"/usr":   {},
		"/var":   {},
	}
)

type sandboxContext struct {
	currentDir string
	in         io.Reader
	out        io.Writer
	errOut     io.Writer
	stack      []sandboxContext
}

type sandboxTerminal struct {
	terminal *term.Terminal
	conf     *internalversion.ExecSandbox
	sandboxContext
}

func execInSandboxSingleCommand(ctx context.Context, conf *internalversion.ExecSandbox, cmd []string, in io.Reader, out, errOut io.Writer) error {
	sctx := sandboxContext{
		currentDir: "/",
		in:         in,
		out:        out,
		errOut:     errOut,
	}
	sc := getExecSandboxExtraCommand(conf, cmd[0])
	if sc != nil {
		return runSandboxCommand(ctx, sc, &sctx)
	}

	fun, ok := sandboxCommand[cmd[0]]
	if !ok {
		return fmt.Errorf("exec failed: unable to start container process: exec: %q: executable file not found in $PATH", cmd[0])
	}

	_ = fun(&sctx, cmd[1:])

	return nil
}

func newSandboxTerminal(in io.Reader, out io.Writer, conf *internalversion.ExecSandbox) *sandboxTerminal {
	stream := struct {
		io.Reader
		io.Writer
	}{
		in,
		out,
	}
	t := term.NewTerminal(stream, "# ")

	st := &sandboxTerminal{
		terminal: t,
		conf:     conf,
		sandboxContext: sandboxContext{
			currentDir: "/",
			out:        t,
			errOut:     t,
		},
	}
	return st
}

func (t *sandboxTerminal) handleCommand(ctx context.Context, line string) bool {
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return true
	}

	cmd := parts[0]
	args := parts[1:]

	sc := getExecSandboxExtraCommand(t.conf, cmd)
	if sc != nil {
		return runSandboxCommand(ctx, sc, &t.sandboxContext) == nil
	}
	if handler, exists := sandboxCommand[cmd]; exists {
		return handler(&t.sandboxContext, args)
	}

	_, _ = io.WriteString(t.sandboxContext.errOut, fmt.Sprintf("%s: command not found\n", cmd))
	return true
}

func (t *sandboxTerminal) Run(ctx context.Context) error {
	for {
		line, err := t.terminal.ReadLine()
		if err != nil {
			return err
		}

		if !t.handleCommand(ctx, line) {
			return nil
		}
	}
}

func getExecSandboxExtraCommand(conf *internalversion.ExecSandbox, cmd string) *internalversion.ExecSandboxExtraCommand {
	for _, c := range conf.ExtraCommands {
		if c.Command == cmd {
			return &c
		}
	}
	return nil
}

func runSandboxCommand(ctx context.Context, cmd *internalversion.ExecSandboxExtraCommand, t *sandboxContext) error {
	for _, step := range cmd.Steps {
		if step.DelayMilliseconds > 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(time.Duration(step.DelayMilliseconds) * time.Millisecond):
			}
		}

		if step.Stdout != "" {
			stdout := step.Stdout
			_, err := io.WriteString(t.out, stdout)
			if err != nil {
				return err
			}
		}

		if step.Stderr != "" {
			stderr := step.Stderr
			_, err := io.WriteString(t.errOut, stderr)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
