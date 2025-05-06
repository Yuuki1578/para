package cmd

import (
	"errors"
	"os/exec"
	"sync"
)

var (
	ErrNilReciever error = errors.New("Reciever is nil")
)

type ShellCmd struct {
	innerCmd []*exec.Cmd
}

func New() *ShellCmd {
	return &ShellCmd{
		innerCmd: make([]*exec.Cmd, 8),
	}
}

func (sh *ShellCmd) Append(name string, args ...string) *ShellCmd {
	if sh == nil {
		return New().Append(name, args...)
	}

	sh.innerCmd = append(sh.innerCmd, exec.Command(name, args...))
	return sh
}

func (sh *ShellCmd) AppendRawCmd(cmd *exec.Cmd) *ShellCmd {
	if sh == nil {
		return New().AppendRawCmd(cmd)
	}

	sh.innerCmd = append(sh.innerCmd, cmd)
	return sh
}

func (sh *ShellCmd) TotalTask() (int, error) {
	if sh == nil {
		return 0, ErrNilReciever
	}

	total := 0

	for _, task := range sh.innerCmd {
		if task != nil {
			total += 1
		}
	}

	return total, nil
}

func (sh *ShellCmd) FixErrorCommand() (*ShellCmd, error) {
	if sh == nil {
		return nil, ErrNilReciever
	}

	clearTask := New()

	for _, cmd := range sh.innerCmd {
		if cmd != nil {
			clearTask.AppendRawCmd(cmd)
		}
	}

	return sh, nil
}

func (sh *ShellCmd) Run(fn func(command *exec.Cmd)) error {
	if sh == nil {
		return ErrNilReciever
	}

	if ok, err := sh.FixErrorCommand(); err != nil {
		return err
	} else {
		sh = ok
	}

	waiter := sync.WaitGroup{}

	switch taskLen, err := sh.TotalTask(); err {
	case nil:
		waiter.Add(taskLen)

	default:
		return err
	}

	for _, command := range sh.innerCmd {
		if command == nil {
			continue
		}

		go func() {
			defer waiter.Done()

			if fn != nil {
				fn(command)
			}
		}()
	}

	waiter.Wait()
	return nil
}
