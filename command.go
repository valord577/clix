package clix

import (
	"bytes"
	"errors"
	"flag"
	"io"
	"os"
	"sort"
	"strings"
)

// @author valor.

// Command is a command for your application.
type Command struct {
	// Name: the command's name
	Name string

	Summary     string
	Description string

	commands map[string]*Command
	parent   *Command

	// Usage4Flags: the usage of flags
	// Example: "[-H | -L | -P] [-EXdsx] [-f path] path ... [expression]"
	Usage4Flags string
	goFlagSet   *flag.FlagSet
	goFlagNum   int
	// ShowDefValue: if true, show flag's defaut value
	ShowDefValue bool

	// Run executes a command
	Run func(c *Command, args []string) error

	// Output: nil means stderr; use out() accessor
	Output io.Writer
}

// AddCmd adds one or more commands to this command.
func (c *Command) AddCmd(cmds ...*Command) error {
	length := len(cmds)
	if length == 0 {
		return nil
	}
	if c.commands == nil {
		c.commands = make(map[string]*Command)
	}

	for i := 0; i < length; i++ {
		if cmds[i] == c {
			return errors.New("can't use itself as a subcommand")
		}
		if cmds[i] == nil {
			return errors.New("can't use nil as a subcommand")
		}

		cmds[i].parent = c
		c.commands[cmds[i].Name] = cmds[i]
	}
	return nil
}

// Flags returns goFlagSet.
func (c *Command) flags() *flag.FlagSet {
	if c.goFlagSet == nil {
		c.goFlagSet = &flag.FlagSet{
			Usage: func() {
				// Use printHelp to print help information
			},
		}

		c.goFlagSet.SetOutput(c.out())
	}

	return c.goFlagSet
}

// Execute uses the args (os.Args[1:])
// and executes the command.
func (c *Command) Execute() error {
	// The command runs on root only
	if c.hasParent() {
		return c.root().Execute()
	}

	args := os.Args[1:]
	cmd, flags, err := findCmd(c, args)
	if err != nil {
		return err
	}

	err = cmd.execute(flags)
	if err == flag.ErrHelp {
		// print help message
		return cmd.printHelp()
	}
	return err
}

func (c *Command) execute(flags []string) error {
	err := c.flags().Parse(flags)
	if err != nil {
		// err may be flag.ErrHelp
		return err
	}

	args := c.flags().Args()
	if args == nil {
		args = []string{}
	}

	if c.Run == nil {
		return c.execute([]string{"-h"})
	}
	return c.Run(c, args)
}

func (c *Command) printHelp() error {
	b := &bytes.Buffer{}

	// summary
	if c.Summary != "" {
		b.WriteString(c.Summary)
		b.WriteByte('\n')
		b.WriteByte('\n')
	}

	// description
	if c.Description != "" {
		b.WriteString(c.Description)
		b.WriteByte('\n')
		b.WriteByte('\n')
	}

	name := c.fullPath()

	// usage: flags
	if c.goFlagNum > 0 {
		b.WriteString("Usage:\n")

		b.WriteString("  ")
		b.WriteString(name)
		b.WriteByte(' ')
		if c.Usage4Flags == "" {
			b.WriteString("[flags...]")
		} else {
			b.WriteString(c.Usage4Flags)
		}
		b.WriteByte('\n')
		b.WriteByte('\n')

		b.WriteString("The available flags are:\n")

		c.flags().VisitAll(func(f *flag.Flag) {

			blank := 20
			length := len(f.Name)

			b.WriteString("  -")
			if length >= blank-2 {
				b.WriteString(f.Name)
				b.WriteByte('\n')
				b.WriteString(strings.Repeat(" ", blank))
				b.WriteString(f.Usage)
			} else {
				b.WriteString(f.Name)
				b.WriteString(strings.Repeat(" ", blank-length-3))
				b.WriteString(f.Usage)
			}

			if c.ShowDefValue {
				b.WriteString(" (default ")
				b.WriteString(f.DefValue)
				b.WriteString(")")
			}
			b.WriteByte('\n')
		})

		b.WriteByte('\n')
	}

	// usage: commands
	length := len(c.commands)
	if length > 0 {
		if c.goFlagNum > 0 {
			b.WriteString("----------------------------------------\n")
			b.WriteByte('\n')
		}

		b.WriteString("Usage:\n")

		b.WriteString("  ")
		b.WriteString(name)
		b.WriteString(" <command> [arguments]")
		b.WriteByte('\n')
		b.WriteByte('\n')

		b.WriteString("The available commands are:\n")

		maxCmdLength := 0
		cmds := make([]string, 0, length)
		for k := range c.commands {
			cmds = append(cmds, k)

			cmdLength := len(k)
			if cmdLength > maxCmdLength {
				maxCmdLength = cmdLength
			}
		}
		sort.Strings(cmds)

		for _, k := range cmds {
			cmd, ok := c.commands[k]
			if cmd == nil || !ok {
				continue
			}

			b.WriteString("  ")
			b.WriteString(cmd.Name)
			b.WriteString(strings.Repeat(" ", 4+maxCmdLength-len(cmd.Name)))
			b.WriteString(cmd.Summary)

			b.WriteByte('\n')
		}

		b.WriteByte('\n')
		b.WriteString("Use \"")
		b.WriteString(name)
		b.WriteString(" <command> -h\" for more information about a command.\n")
		b.WriteByte('\n')
	}

	_, err := c.out().Write(b.Bytes())
	return err
}

func (c *Command) root() *Command {
	if c.hasParent() {
		c.parent.root()
	}
	return c
}

func (c *Command) hasParent() bool {
	return c.parent != nil
}

func (c *Command) chlid(name string) *Command {
	return c.commands[name]
}

func (c *Command) fullPath() string {
	if c.hasParent() {
		return c.parent.fullPath() + " " + c.Name
	}
	return c.Name
}

func (c *Command) out() io.Writer {
	if c.Output == nil {
		return os.Stderr
	}
	return c.Output
}

func isFlagArg(arg string) bool {
	return len(arg) >= 2 && arg[0] == '-'
}

func findCmd(c *Command, args []string) (*Command, []string, error) {
	if c == nil {
		return nil, nil, errors.New("command is nil")
	}

	length := len(args)
	if length == 0 {
		return c, args, nil
	}

	s := args[0]
	if isFlagArg(s) {
		return c, args, nil
	}

	cmd := c.chlid(s)
	if cmd == nil {
		return nil, nil, errors.New("command not found")
	}
	return findCmd(cmd, args[1:])
}
