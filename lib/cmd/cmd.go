package cmd

import (
	"errors"
	"github.com/Yuuki1578/para/lib/env"
	"os/exec"
	"sync"
)

var (
	ErrNilReciever error = errors.New("Reciever is nil")
	ErrNoSuchTask  error = errors.New("No such command")
)

type Command struct {
	backing []string
	command *exec.Cmd
	count   uint64
}

func newCommand(count uint64, name string, args ...string) Command {
	return Command{
		backing: append([]string{name}, args...),
		command: exec.Command(name, args...),
		count:   count,
	}
}

type CommandGroup struct {
	commandList []Command
}

func New() *CommandGroup {
	return &CommandGroup{
		commandList: make([]Command, 0, 8),
	}
}

func (group *CommandGroup) Append(count uint64, name string, args ...string) *CommandGroup {
	command := newCommand(count, name, args...)
	group.commandList = append(group.commandList, command)
	return group
}

func (group *CommandGroup) Fix() {
	temporary := make([]Command, 0, cap(group.commandList))
	for _, instance := range group.commandList {
		if instance.command != nil {
			temporary = append(temporary, instance)
		}
	}

	group.commandList = temporary
}

func (group CommandGroup) Total() int {
	total := 0
	for _, instance := range group.commandList {
		if instance.command != nil {
			total++
		}
	}

	return total
}

func (group CommandGroup) Capacity() int {
	return cap(group.commandList)
}

func singleCommand(
	waiter *sync.WaitGroup,
	instanceOf *Command,
	fn func(command *exec.Cmd),
) {
	defer waiter.Done()
	fn(instanceOf.command)
}

func multiCommand(
	waiter *sync.WaitGroup,
	instanceOf *Command,
	fn func(command *exec.Cmd),
) {
	count := instanceOf.count
	name := instanceOf.backing
	length := len(name)
	localWait := sync.WaitGroup{}
	defer waiter.Done()

	if usingEnv := env.GetEnv(); usingEnv != 0 {
		count = usingEnv
	}

	fnOnce := func() {
		fn(exec.Command(name[0]))
		localWait.Done()
	}

	fnMul := func() {
		fn(exec.Command(name[0], name...))
		localWait.Done()
	}

	if length == 0 {
		return
	}

	localWait.Add(int(count))
	for range count {
		switch length {
		case 1:
			go fnOnce()

		default:
			go fnMul()
		}
	}

	localWait.Wait()
}

func (group *CommandGroup) Run(fn func(command *exec.Cmd)) error {
	group.Fix()

	waiter := sync.WaitGroup{}
	total := group.Total()
	if total <= 0 {
		return ErrNoSuchTask
	}

	waiter.Add(total)
	for _, instance := range group.commandList {
		count := instance.count
		if instance.command == nil {
			continue
		}

		switch count {
		case 0, 1:
			go singleCommand(&waiter, &instance, fn)

		default:
			go multiCommand(&waiter, &instance, fn)
		}
	}

	waiter.Wait()
	return nil
}
