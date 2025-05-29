package console

import (
	"context"
	"fmt"
	"log"
	"os"
)

// Manager manages console commands
type Manager struct {
	commands map[string]Command
}

// NewManager creates a new command manager
func NewManager() *Manager {
	return &Manager{
		commands: make(map[string]Command),
	}
}

// Register registers a command
func (m *Manager) Register(cmd Command) {
	config := &CommandConfig{}
	cmd.Configure(config)
	log.Printf("Registering command: %s", config.Name)
	m.commands[config.Name] = cmd
}

// FindCommand finds a command by name
func (m *Manager) FindCommand(name string) Command {
	log.Printf("Looking for command: %s", name)
	log.Printf("Available commands: %v", m.commands)
	return m.commands[name]
}

// RunFromArgs runs a command from command line arguments
func (m *Manager) RunFromArgs() error {
	args := os.Args[1:]
	if len(args) == 0 {
		return m.showAvailableCommands()
	}

	cmdName := args[0]
	cmd := m.FindCommand(cmdName)
	if cmd == nil {
		return fmt.Errorf("command not found: %s", cmdName)
	}

	// Create context with arguments
	ctx := context.WithValue(context.Background(), "args", args)
	return cmd.Handle(ctx)
}

// showAvailableCommands shows all available commands
func (m *Manager) showAvailableCommands() error {
	fmt.Println("Available commands:")
	for name, cmd := range m.commands {
		fmt.Printf("  %s\t%s\n", name, cmd.GetDescription())
	}
	return nil
}
