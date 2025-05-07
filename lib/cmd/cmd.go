package cmd

import (
	"errors"
	"os/exec"
	"sync"
)

var (
	ErrNilReciever error = errors.New("Reciever is nil")
	ErrNoSuchTask  error = errors.New("No such command")
)

type Command struct {
	backup  []string
	command *exec.Cmd
	count   uint64
}

func newCommand(count uint64, name string, args ...string) Command {
	return Command{
		backup:  append([]string{name}, args...),
		command: exec.Command(name, args...),
		count:   count,
	}
}

type CommandGroup struct {
	commandList []Command
}

func Init() *CommandGroup {
	return &CommandGroup{
		commandList: make([]Command, 0, 8),
	}
}

func (group *CommandGroup) Append(count uint64, name string, args ...string) *CommandGroup {
	if group == nil {
		return Init().Append(count, name, args...)
	}

	command := newCommand(count, name, args...)
	group.commandList = append(group.commandList, command)
	return group
}

func (group *CommandGroup) Fix() {
	if group == nil {
		return
	}

	taskTmp := make([]Command, 0, cap(group.commandList))
	for _, instance := range group.commandList {
		if instance.command != nil {
			taskTmp = append(taskTmp, instance)
		}
	}

	group.commandList = taskTmp
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

func (group *CommandGroup) Run(fn func(command *exec.Cmd)) error {
	if group == nil || fn == nil {
		return ErrNilReciever
	}

	group.Fix()
	waiter := sync.WaitGroup{}
	taskLen := group.Total()
	if taskLen <= 0 {
		return ErrNoSuchTask
	}

	waiter.Add(taskLen)
	for _, instance := range group.commandList {
		count := instance.count
		singleCount := func() {
			defer waiter.Done()
			fn(instance.command)
		}

		multiCount := func() {
			rawName := instance.backup
			rawLen := len(rawName)
			localWait := sync.WaitGroup{}

			defer waiter.Done()
			if rawLen == 0 {
				return
			}

			localWait.Add(int(count))
			for range count {
				switch rawLen {
				case 1:
					go func() {
						defer localWait.Done()
						fn(exec.Command(rawName[0]))
					}()

				default:
					go func() {
						defer localWait.Done()
						fn(exec.Command(rawName[0], rawName...))
					}()
				}
			}

			localWait.Wait()
		}

		if instance.command == nil {
			continue
		}

		switch count {
		case 0, 1:
			go singleCount()

		default:
			go multiCount()
		}
	}

	waiter.Wait()
	return nil
}
