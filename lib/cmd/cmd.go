package cmd

import (
	"errors"
	"os"
	"os/exec"
	"sync"
)

var (
	ErrCommandFailed error = errors.New("Command failed to be executed")
	ErrNoSuchCommand error = errors.New("No such command in entry list")
)

type ShellCmd struct {
	innerCmd []*exec.Cmd
}

func New() ShellCmd {
	return ShellCmd{
		innerCmd: make([]*exec.Cmd, 8),
	}
}

func (sh *ShellCmd) Append(name string, args ...string) {
	if sh != nil {
		sh.innerCmd = append(sh.innerCmd, exec.Command(name, args...))
	}
}

func (sh *ShellCmd) Run() error {
	if sh == nil {
		return ErrCommandFailed
	}

	waiter := sync.WaitGroup{}
	waitable := len(sh.innerCmd)

	if waitable == 0 {
		return ErrNoSuchCommand
	}

	waiter.Add(waitable)

	for _, command := range sh.innerCmd {
		go func() {
			stdout, err := command.Output()
			defer waiter.Done()

			if err != nil {
				return
			}

			if _, err = os.Stdout.Write(stdout); err == nil {
				os.Stdout.Sync()
			}
		}()
	}

	waiter.Wait()
	return nil
}
