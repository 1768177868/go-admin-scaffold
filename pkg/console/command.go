package console

import (
	"context"
	"fmt"
)

// Command represents a console command interface
type Command interface {
	GetName() string
	GetDescription() string
	Handle(ctx context.Context) error
	AddArgument(name, description string)
	GetArgument(name string) string
	SetValue(name, value string)
	Configure(config *CommandConfig)
}

// BaseCommand provides a basic implementation of Command interface
type BaseCommand struct {
	Name        string
	Description string
	Usage       string
	Arguments   map[string]string
	Options     []Option
	values      map[string]string
}

func NewCommand(name, description string) *BaseCommand {
	return &BaseCommand{
		Name:        name,
		Description: description,
		Arguments:   make(map[string]string),
		Options:     make([]Option, 0),
		values:      make(map[string]string),
	}
}

func (c *BaseCommand) GetName() string {
	return c.Name
}

func (c *BaseCommand) GetDescription() string {
	return c.Description
}

func (c *BaseCommand) AddArgument(name, description string) {
	c.Arguments[name] = description
}

func (c *BaseCommand) GetArgument(name string) string {
	return c.values[name]
}

func (c *BaseCommand) SetValue(name, value string) {
	c.values[name] = value
}

func (c *BaseCommand) Handle(ctx context.Context) error {
	return fmt.Errorf("handle method not implemented")
}

// CommandConfig holds command configuration
type CommandConfig struct {
	Name        string
	Description string
	Usage       string
	Arguments   []Argument
	Options     []Option
}

// Argument represents a command argument
type Argument struct {
	Name        string
	Description string
	Required    bool
	Value       string
}

// Option represents a command option
type Option struct {
	Name        string
	Shortcut    string
	Description string
	Required    bool
	Value       string
}

// GetUsage returns the command usage
func (c *BaseCommand) GetUsage() string {
	return c.Usage
}

// Configure sets up the command
func (c *BaseCommand) Configure(config *CommandConfig) {
	c.Name = config.Name
	c.Description = config.Description
	c.Usage = config.Usage
	for _, arg := range config.Arguments {
		c.AddArgument(arg.Name, arg.Description)
	}
	c.Options = config.Options
}

// GetOption returns the value of an option
func (c *BaseCommand) GetOption(name string) string {
	for _, opt := range c.Options {
		if opt.Name == name {
			return opt.Value
		}
	}
	return ""
}

// HasOption checks if an option exists and has a value
func (c *BaseCommand) HasOption(name string) bool {
	return c.GetOption(name) != ""
}

// Line writes a line to the output
func (c *BaseCommand) Line(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

// Info writes an info message
func (c *BaseCommand) Info(format string, args ...interface{}) {
	c.Line("INFO: "+format, args...)
}

// Error writes an error message
func (c *BaseCommand) Error(format string, args ...interface{}) {
	c.Line("ERROR: "+format, args...)
}

// Success writes a success message
func (c *BaseCommand) Success(format string, args ...interface{}) {
	c.Line("SUCCESS: "+format, args...)
}
